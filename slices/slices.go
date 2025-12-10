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
	"slices"
)

type (
	// A PredicateFunc returns true or false based off of a provided value.
	PredicateFunc[T any] = func(T) bool
	// A TransformFunc (infallibly) transforms an object of type In into an
	// object of type Out, returning the result.
	TransformFunc[In any, Out any] = func(In) Out
	// A TransformErrorFunc transforms an object of type In into an object of
	// type Out, returning the result or an error.
	TransformErrorFunc[In any, Out any] = func(In) (Out, error)
)

// HasPrefix evaluates if x contains prefix as its first len(prefix) elements.
func HasPrefix[S ~[]T, T comparable](x S, prefix S) bool {
	plen := len(prefix)
	if plen > len(x) {
		return false
	}
	return slices.Equal(x[:plen], prefix)
}

// Filter returns a copy of x with any elements for which pred returns true.
func Filter[S ~[]T, T any, Fn ~PredicateFunc[T]](x S, pred Fn) []T {
	if len(x) == 0 {
		return nil
	}

	var tmp []T //nolint:prealloc
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
func Transform[S ~[]In, In any, Out any, Fn ~TransformFunc[In, Out]](
	x S,
	transform Fn,
) []Out {
	if len(x) == 0 || transform == nil {
		return nil
	}

	dst := make([]Out, len(x))
	for i := range x {
		dst[i] = transform(x[i])
	}

	return dst
}

// TransformError returns a copy of x with all elements' values passed through
// the given mapping function. If an error is encountered, iteration halts, any
// transformed items are discarded, and the error is immediately returned to
// the caller.
func TransformError[S ~[]In, In any, Out any, Fn ~TransformErrorFunc[In, Out]](
	x S,
	transform Fn,
) ([]Out, error) {
	if len(x) == 0 || transform == nil {
		return nil, nil
	}

	var (
		dst = make([]Out, len(x))
		err error
	)
	for i := range x {
		if dst[i], err = transform(x[i]); err != nil {
			return nil, err
		}
	}

	return dst, nil
}

// Count counts the number of elements in s for which pred returns true.
func Count[S ~[]T, T any, P ~PredicateFunc[T]](s S, pred P) int {
	var n int
	for i := range s {
		if pred(s[i]) {
			n++
		}
	}
	return n
}

// CountPtr counts the number of elements in s for which pred returns true.
func CountPtr[S ~[]T, T any, P ~PredicateFunc[*T]](s S, pred P) int {
	var n int
	for i := range s {
		if pred(&s[i]) {
			n++
		}
	}
	return n
}

// An AssociativeFunc outputs both a key and value from a single input.
type AssociativeFunc[In any, OutK comparable, OutV any] = func(In) (OutK, OutV)

// Map converts s into map[OutK]OutV according to fn.
func Map[
	OutK comparable,
	OutV any,
	S ~[]In,
	F AssociativeFunc[In, OutK, OutV],
	In any,
](s S, fn F) map[OutK]OutV {
	if len(s) == 0 || fn == nil {
		return nil
	}

	dst := make(map[OutK]OutV, len(s))
	for i := range s {
		k, v := fn(s[i])
		dst[k] = v
	}
	return dst
}
