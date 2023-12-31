// Copyright (c) 2023 Matt Way
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

// Package stub provides stubbing utilities, primarily for tests.
package stub

// With executes fn with the value at dst replaced by stub for the span of fn.
func With[T any](dst *T, stub T, fn func()) {
	defer swap(dst, stub)()
	fn()
}

// WithError executes fn and returns its result, with the value at dst replaced
// by stub for the span of fn.
func WithError[T any](dst *T, stub T, fn func() error) error {
	defer swap(dst, stub)()
	return fn()
}

func swap[T any](dst *T, stub T) func() {
	orig := *dst
	*dst = stub
	return func() {
		*dst = orig
	}
}
