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

package os

import (
	"errors"
	"io"
	"io/fs"
	"os"
)

const _createFileFlags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY

// WithFileReader opens the file at path for reading and passes the file's
// reader to fn, cleaning up once fn returns. WithFileReader returns any error
// returned by fn as well as any other errors encountered.
func WithFileReader(path string, fn func(r io.Reader) error) (err error) {
	var src io.ReadCloser
	if src, err = os.Open(path); err != nil { //nolint:gosec
		return err
	}
	defer func() {
		err = errors.Join(err, src.Close())
	}()

	return fn(src)
}

// WithFileWriter opens the file at path for writing and passes the file's
// writer to fn, cleaning up once fn returns. WithFileWriter returns any error
// returned by fn as well as any other errors encountered.
func WithFileWriter(path string, fn func(w io.Writer) error) (err error) {
	return WithFileModeWriter(path, 0o744, fn)
}

// WithFileModeWriter opens (creating with the given mode or truncating, as
// necessary) the file at path for writing and passes the file's writer to fn,
// cleaning up once fn returns. WithFileModeWriter returns any error returned
// by fn as well as any other errors encountered.
func WithFileModeWriter(
	path string,
	mode fs.FileMode,
	fn func(w io.Writer) error,
) (err error) {
	var dst io.WriteCloser
	if dst, err = os.OpenFile(path, _createFileFlags, mode); err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, dst.Close())
	}()

	return fn(dst)
}
