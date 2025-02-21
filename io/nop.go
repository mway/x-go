// Copyright (c) 2025 Matt Way
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

// Package io provides I/O-related helpers and utilities.
package io

import (
	"io"
)

var (
	// Nop is an [io.ReadWriteCloser] that does nothing.
	Nop io.ReadWriteCloser = nop{}
	// NopReader is an [io.Reader] that reads no data and always returns EOF.
	NopReader io.Reader = nop{}
	// NopReadCloser is an [io.ReadCloser] that reads no data, always returns
	// EOF, and does not error on close.
	NopReadCloser io.ReadCloser = nop{}
	// NopWriter is an [io.Writer] that writes no data and always returns nil.
	NopWriter io.Writer = nop{}
	// NopWriteCloser is an [io.WriteCloser] that writes no data, always
	// returns nil, and does not error on close.
	NopWriteCloser io.WriteCloser = nop{}
	// NopStringWriter is an [io.StringWriter] that writes no data and always
	// returns nil.
	NopStringWriter io.StringWriter = nop{}
)

type nop struct{}

func (nop) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (nop) Write(p []byte) (int, error) {
	return io.Discard.Write(p)
}

func (nop) WriteString(s string) (int, error) {
	return io.Discard.(io.StringWriter).WriteString(s) //nolint:errcheck
}

func (nop) Close() error {
	return nil
}
