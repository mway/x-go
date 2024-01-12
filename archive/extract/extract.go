// Copyright (c) 2024 Matt Way
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE THE SOFTWARE.

// Package extract provides archive extraction helpers.
package extract

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"
	"go.mway.dev/color"
	"go.mway.dev/errors"
	xos "go.mway.dev/x/os"
)

const (
	// ExtractToTempDir is a sentinel destination value that will cause
	// [Extract] to generate a temporary directory for extraction, and remove
	// the temporary directory once extraction (and callbacks) have completed.
	// ExtractToTempDir is only useful if a [Callback] is passed to Extract.
	ExtractToTempDir = ""

	_sep = string(os.PathSeparator)
)

// Extract extracts the given archive to dst using any provided [Option]s.
//
//nolint:gocyclo
func Extract(
	ctx context.Context,
	dst string,
	archive string,
	opts ...Option,
) (err error) {
	var options Options
	for _, opt := range opts {
		opt.apply(&options)
	}

	if options.Output == nil {
		options.Output = io.Discard
	}

	if dst == ExtractToTempDir {
		if dst, err = os.MkdirTemp("", "extract"); err != nil {
			return errors.Wrap(err, "failed to create temporary directory")
		}

		defer func() {
			err = errors.Join(err, errors.Wrapf(
				os.RemoveAll(dst),
				"failed to remove temporary destination %q",
				dst,
			))
		}()
	}

	dirs := make(map[string]struct{})
	handler := archiver.FileHandler(func(
		_ context.Context,
		f archiver.File,
	) (err error) {
		// Ignore dirs; they are created lazily below.
		if f.IsDir() {
			return nil
		}

		fpath, stripped, stripErr := stripPrefix(
			f.NameInArchive,
			options.StripPrefix,
		)
		switch {
		case stripErr != nil:
			return errors.Wrap(stripErr, "failed to strip prefix")
		case !stripped:
			return nil
		default:
			// passthrough
		}

		for _, exclude := range options.ExcludePaths {
			matched, matchErr := filepath.Match(exclude, fpath)
			if matchErr != nil {
				return errors.Wrapf(matchErr, "bad match pattern %q", exclude)
			}
			if matched {
				return nil
			}
		}

		dstpath := filepath.Join(dst, fpath)
		if len(options.IncludePaths) > 0 {
			var matched bool
			for include, explicitDst := range options.IncludePaths {
				if matched, err = filepath.Match(include, fpath); err != nil {
					return errors.Wrapf(err, "bad match pattern %q", include)
				}

				switch {
				case !matched:
					continue
				case len(explicitDst) == 0:
					explicitDst = fpath
				case !filepath.IsAbs(explicitDst):
					explicitDst = filepath.Join(dst, explicitDst)
				}

				dstpath = explicitDst
				break
			}

			if !matched {
				return nil
			}
		}

		parent := filepath.Dir(fpath)
		if _, parentCreated := dirs[parent]; !parentCreated {
			if err = xos.MkdirAllInherit(parent); err != nil {
				return errors.Wrapf(
					err,
					"failed to create destination parent(s) %q",
					parent,
				)
			}
			dirs[fpath] = struct{}{}
		}

		_, err = xos.WriteReaderToFileWithFlags(
			dstpath,
			f.Open,
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			f.Mode(),
		)
		if err != nil {
			return errors.Wrapf(
				err,
				"failed to extract %q to %q",
				fpath,
				dstpath,
			)
		}

		_, err = color.FgHiGreen.Fprint(options.Output, "Extracting:")
		if err != nil {
			return errors.Wrap(err, "failed to write to output")
		}

		_, err = fmt.Fprintln(options.Output, "", fpath, "->", dstpath)
		return errors.Wrap(err, "failed to write to output")
	})

	var src io.ReadCloser
	if src, err = os.Open(archive); err != nil {
		return errors.Wrapf(err, "failed to open source file %q", archive)
	}
	defer func() {
		err = errors.Join(err, errors.Wrapf(
			src.Close(),
			"failed to close source file %q",
			archive,
		))
	}()

	return xos.WithCwd(dst, func() (err error) {
		var (
			format archiver.Format
			reader io.Reader
		)
		if format, reader, err = archiver.Identify(archive, src); err != nil {
			return errors.Wrapf(err, "failed to detect format of %q", archive)
		}

		ex, ok := format.(archiver.Extractor)
		if !ok {
			return errors.Newf(
				"bug: identified format (%T) is not extractable",
				format,
			)
		}

		if _, ok := ex.(archiver.Zip); ok {
			raw, readErr := io.ReadAll(reader)
			if readErr != nil && !errors.Is(readErr, io.EOF) {
				return errors.Wrap(readErr, "failed to read data buffer")
			}
			reader = bytes.NewReader(raw)
		}

		var extractAll []string
		if err = ex.Extract(ctx, reader, extractAll, handler); err != nil {
			return errors.Wrapf(err, "failed to extract %q", archive)
		}

		if options.Callback != nil {
			return errors.Wrap(
				options.Callback(ctx, dst),
				"error during callback",
			)
		}

		return nil
	})
}

func stripPrefix(path string, prefix string) (string, bool, error) {
	prefix = strings.Trim(prefix, _sep)
	if len(prefix) == 0 {
		return path, true, nil
	}

	path = filepath.Clean(path)
	idx := strings.IndexByte(path, _sep[0])
	for idx > 0 {
		matched, err := filepath.Match(prefix, path[:idx])
		if err != nil {
			return "", false, errors.Wrapf(err, "bad prefix %q", prefix)
		}

		if matched {
			return path[idx+1:], true, nil
		}

		offset := strings.IndexByte(path[idx+1:], _sep[0])
		if offset < 0 {
			break
		}

		idx += offset + 1
	}

	return path, false, nil
}
