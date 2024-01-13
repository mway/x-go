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

// Package http provides http-related utilities and helpers.
package http

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"go.mway.dev/errors"
	xos "go.mway.dev/x/os"
	"golang.org/x/sys/unix"
)

var (
	// ErrDestNotWritable is returned when a given destination is a directory
	// that is not writable.
	ErrDestNotWritable = errors.New("destination is not writable")

	_httpGet    = http.Get
	_osStat     = os.Stat
	_osMkdirAll = os.MkdirAll
	_unixAccess = unix.Access
)

// GetFile retrieves a file referenced by the given url and writes it to the
// given destination, returning any errors in the process. The resulting file
// may have been written despite returning a non-nil error. If dst does not
// exist, GetFile will make its parent directory if needed before writing; if
// dst exists and is a directory, GetFile will use the basename of the URL to
// inform the resulting filename.
func GetFile(url string, dst string) (err error) {
	var resp *http.Response
	if resp, err = _httpGet(url); err != nil {
		return err
	}

	if resp.Body != nil {
		defer func() {
			err = errors.Join(err, errors.Wrap(
				resp.Body.Close(),
				"failed to close response body",
			))
		}()
	}

	var dstInfo fs.FileInfo
	if dstInfo, err = _osStat(dst); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err = _osMkdirAll(filepath.Dir(dst), 0o755)
		}
	}

	switch {
	case err != nil:
		return err
	case dstInfo != nil && dstInfo.IsDir():
		if base := filepath.Base(url); len(base) > 0 {
			dst = filepath.Join(dst, base)
		}
	}

	if parent := filepath.Dir(dst); !existsAndWritable(parent) {
		return errors.Wrap(ErrDestNotWritable, parent)
	}

	_, err = xos.WriteReaderToFile(dst, resp.Body)
	return err
}

func existsAndWritable(path string) bool {
	return _unixAccess(path, unix.W_OK) == nil
}
