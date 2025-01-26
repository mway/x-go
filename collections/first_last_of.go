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

// Package collections provides helpers for collections of objects, like
// slices, maps, and sets.
package collections

// FirstOf returns the first T in values that is not a zero value of T.
func FirstOf[T comparable](values ...T) T {
	var zero T
	for _, value := range values {
		if value != zero {
			return value
		}
	}
	return zero
}

// FirstOfFuncs returns the first T produced by a function in fns that is not a
// zero value of T.
func FirstOfFuncs[T comparable](fns ...func() T) T {
	var zero T
	for i := 0; i < len(fns); i++ {
		if fns[i] == nil {
			continue
		}
		if x := fns[i](); x != zero {
			return x
		}
	}
	return zero
}

// FirstOfOr returns the first T in values that is not a zero value of T, or
// invokes fn to produce a T otherwise.
func FirstOfOr[T comparable](fn func() T, values ...T) T {
	var zero T
	if x := FirstOf(values...); x != zero {
		return x
	}
	if fn != nil {
		return fn()
	}
	return zero
}

// LastOf returns the last T in values that is not a zero value of T.
func LastOf[T comparable](values ...T) T {
	var zero T
	for i := len(values) - 1; i >= 0; i-- {
		if value := values[i]; value != zero {
			return value
		}
	}
	return zero
}

// LastOfFuncs returns the last T produced by a function in fns that is not a
// zero value of T. The given functions will be evaluated in reverse order;
// functions provided earlier than a successful function will not be called.
func LastOfFuncs[T comparable](fns ...func() T) T {
	var zero T
	for i := len(fns) - 1; i >= 0; i-- {
		if fns[i] == nil {
			continue
		}
		if value := fns[i](); value != zero {
			return value
		}
	}
	return zero
}

// LastOfOr returns the last T in values that is not a zero value of T, or
// invokes fn to produce a T otherwise.
func LastOfOr[T comparable](fn func() T, values ...T) T {
	var zero T
	if x := LastOf(values...); x != zero {
		return x
	}
	if fn != nil {
		return fn()
	}
	return zero
}
