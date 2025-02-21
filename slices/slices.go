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

// Package slices provides slice-related utilities.
package slices

import (
	"iter"
	"slices"
)

// HasPrefix evaluates if x contains prefix as its first len(prefix) elements.
func HasPrefix[T comparable](x []T, prefix []T) bool {
	plen := len(prefix)
	if plen > len(x) {
		return false
	}
	return slices.Equal(x[:plen], prefix)
}

// Filter returns a copy of x with any elements for which pred returns true.
func Filter[T any, P ~func(T) bool](x []T, pred P) []T {
	if len(x) == 0 {
		return nil
	}

	var tmp []T
	for i := range x {
		if !pred(x[i]) {
			continue
		}

		if tmp == nil {
			tmp = make([]T, 0, len(x)-i)
		}

		tmp = append(tmp, x[i])
	}
	return tmp
}

// Transform returns a copy of x with all elements' values passed through the
// given mapping function.
func Transform[From any, To any, P ~func(From) To](x []From, mapper P) []To {
	if len(x) == 0 || mapper == nil {
		return nil
	}

	dst := make([]To, len(x))
	for i := range x {
		dst[i] = mapper(x[i])
	}

	return dst
}

// Iter returns an [iter.Seq[V]] that ranges over s.
func Iter[T ~[]V, V any](s T) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, value := range s {
			if !yield(value) {
				return
			}
		}
	}
}

// Iter2 returns an [iter.Seq[int, V]] that ranges over s.
func Iter2[T ~[]V, V any](s T) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, value := range s {
			if !yield(i, value) {
				return
			}
		}
	}
}
