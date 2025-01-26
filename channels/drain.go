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

package channels

import (
	"context"
)

// Drain drains ch until it does not yield further values, returning the number
// of values received.
func Drain[T any](ch <-chan T) int {
	var n int
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return n
			}
			n++
		default:
			return n
		}
	}
}

// DrainContext drains ch until either ctx is canceled or ch does not yield
// further values, returning the number of values received.
func DrainContext[T any](ctx context.Context, ch <-chan T) int {
	var n int
	for {
		select {
		case <-ctx.Done():
			return n
		case _, ok := <-ch:
			if !ok {
				return n
			}
			n++
		}
	}
}
