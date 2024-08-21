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

package linebuffer

import (
	"io"
)

type bufferOptions struct {
	Source io.Reader
}

func (o bufferOptions) With(opts ...Option) bufferOptions {
	for _, opt := range opts {
		opt.apply(&o)
	}
	return o
}

// An Option configures a [Buffer].
type Option interface {
	apply(*bufferOptions)
}

// WithSource returns a new [Option] that configures a [Buffer] to read lines
// from the given [io.Reader].
func WithSource(r io.Reader) Option {
	return bufferOptionFunc(func(dst *bufferOptions) {
		dst.Source = r
	})
}

type bufferOptionFunc func(*bufferOptions)

func (f bufferOptionFunc) apply(dst *bufferOptions) {
	f(dst)
}
