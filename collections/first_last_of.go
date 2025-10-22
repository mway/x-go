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

import (
	"iter"
)

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

// FirstOfOr returns the first T in values that is not a zero value of T, or
// otherwise returns the given fallback.
func FirstOfOr[T comparable](fallback T, values ...T) T {
	var zero T
	if x := FirstOf(values...); x != zero {
		return x
	}
	return fallback
}

// FirstOfOrElse returns the first T in values that is not a zero value of T,
// or otherwise invokes fallback (if non-nil) to produce a return value.
func FirstOfOrElse[T comparable](fallback func() T, values ...T) T {
	var zero T
	if x := FirstOf(values...); x != zero {
		return x
	}
	if fallback != nil {
		return fallback()
	}
	return zero
}

// FirstOfSeq returns the first T in seq that is not a zero value of T.
func FirstOfSeq[T comparable](seq iter.Seq[T]) (value T) {
	var zero T
	seq(func(x T) bool {
		value = x
		return value == zero
	})
	return value
}

// FirstOfSeqOr returns the first V in seq that is not a zero value of V, or
// otherwise returns the given fallback.
func FirstOfSeqOr[T comparable](fallback T, seq iter.Seq[T]) T {
	var zero T
	if x := FirstOfSeq(seq); x != zero {
		return x
	}
	return fallback
}

// FirstOfSeqOrElse returns the first V in seq that is not a zero value of V,
// or otherwise invokes fallback (if non-nil) to produce a return value.
func FirstOfSeqOrElse[T comparable](fallback func() T, seq iter.Seq[T]) T {
	var zero T
	if x := FirstOfSeq(seq); x != zero {
		return x
	}
	if fallback != nil {
		return fallback()
	}
	return zero
}

// FirstOfSeq2 returns the first (K,V) in seq where V is not the zero value.
func FirstOfSeq2[K comparable, V comparable](
	seq iter.Seq2[K, V],
) (key K, value V) {
	var zero V
	seq(func(k K, v V) bool {
		if v != zero {
			key = k
			value = v
			return false
		}
		return true
	})
	return key, value
}

// FirstOfSeq2Or returns the first (K,V) in seq where V is not the zero value,
// or otherwise returns the given fallbacks.
func FirstOfSeq2Or[K comparable, V comparable](
	fallbackKey K,
	fallbackValue V,
	seq iter.Seq2[K, V],
) (K, V) {
	var zero V
	if k, v := FirstOfSeq2(seq); v != zero {
		return k, v
	}
	return fallbackKey, fallbackValue
}

// FirstOfSeq2OrElse returns the first (K,V) in seq where V is not the zero
// value, or otherwise invokes fallback (if non-nil) to produce a return value.
func FirstOfSeq2OrElse[K comparable, V comparable](
	fallback func() (K, V),
	seq iter.Seq2[K, V],
) (K, V) {
	var zero V
	if k, v := FirstOfSeq2(seq); v != zero {
		return k, v
	}
	if fallback != nil {
		return fallback()
	}
	var key K
	return key, zero
}

// FirstOfFuncs returns the first T produced by a non-nil function in fns that
// is not a zero value of T. Note that because the given functions will be
// evaluated in order, higher-indexed functions may not be called.
func FirstOfFuncs[T comparable](fns ...func() T) T {
	var zero T
	for i := range len(fns) {
		if fns[i] == nil {
			continue
		}
		if x := fns[i](); x != zero {
			return x
		}
	}
	return zero
}

// FirstOfFuncsOr returns the first T produced by a non-nil function in fns
// that is not a zero value of T, or otherwise returns the given fallback.
// Note that because the given functions will be evaluated in order,
// higher-indexed functions may not be called.
func FirstOfFuncsOr[T comparable](fallback T, fns ...func() T) T {
	var zero T
	if x := FirstOfFuncs(fns...); x != zero {
		return x
	}
	return fallback
}

// FirstOfFuncsOrElse returns the first T produced by a non-nil function in fns
// that is not a zero value of T, or otherwise invokes fallback (if non-nil) to
// produce a return value. Note that because the given functions will be
// evaluated in order, higher-indexed functions may not be called.
func FirstOfFuncsOrElse[T comparable](fallback func() T, fns ...func() T) T {
	var zero T
	if x := FirstOfFuncs(fns...); x != zero {
		return x
	}
	if fallback != nil {
		return fallback()
	}
	return zero
}

// LastOf returns the last T in values that is not a zero value of T.
func LastOf[T comparable](values ...T) T {
	var zero T
	for i := len(values) - 1; i >= 0; i-- {
		if x := values[i]; x != zero {
			return x
		}
	}
	return zero
}

// LastOfOr returns the last T in values that is not a zero value of T, or
// otherwise returns the given fallback.
func LastOfOr[T comparable](fallback T, values ...T) T {
	var zero T
	if x := LastOf(values...); x != zero {
		return x
	}
	return fallback
}

// LastOfOrElse returns the last T in values that is not a zero value of T, or
// otherwise invokes fallback (if non-nil) to produce a return value.
func LastOfOrElse[T comparable](fn func() T, values ...T) T {
	var zero T
	for i := len(values) - 1; i >= 0; i-- {
		if value := values[i]; value != zero {
			return value
		}
	}
	if fn != nil {
		return fn()
	}
	return zero
}

// LastOfSeq returns the last T in seq that is not a zero value of T.
func LastOfSeq[T comparable](seq iter.Seq[T]) (value T) {
	var zero T
	seq(func(x T) bool {
		if x != zero {
			value = x
		}
		return true
	})
	return value
}

// LastOfSeqOr returns the last V in seq that is not a zero value of V, or
// otherwise returns the given fallback.
func LastOfSeqOr[T comparable](fallback T, seq iter.Seq[T]) T {
	var zero T
	if x := LastOfSeq(seq); x != zero {
		return x
	}
	return fallback
}

// LastOfSeqOrElse returns the last V in seq that is not a zero value of V,
// or otherwise invokes fallback (if non-nil) to produce a return value.
func LastOfSeqOrElse[T comparable](fallback func() T, seq iter.Seq[T]) T {
	var zero T
	if x := LastOfSeq(seq); x != zero {
		return x
	}
	if fallback != nil {
		return fallback()
	}
	return zero
}

// LastOfSeq2 returns the last (K,V) in seq where V is not the zero value.
func LastOfSeq2[K comparable, V comparable](
	seq iter.Seq2[K, V],
) (key K, value V) {
	var zero V
	seq(func(k K, v V) bool {
		if v != zero {
			key = k
			value = v
		}
		return true
	})
	return key, value
}

// LastOfSeq2Or returns the last (K,V) in seq where V is not the zero value, or
// otherwise returns the given fallbacks.
func LastOfSeq2Or[K comparable, V comparable](
	fallbackKey K,
	fallbackValue V,
	seq iter.Seq2[K, V],
) (K, V) {
	var zero V
	if k, v := LastOfSeq2(seq); v != zero {
		return k, v
	}
	return fallbackKey, fallbackValue
}

// LastOfSeq2OrElse returns the last (K,V) in seq where V is not the zero
// value, or otherwise invokes fallback (if non-nil) to produce a return value.
func LastOfSeq2OrElse[K comparable, V comparable](
	fallback func() (K, V),
	seq iter.Seq2[K, V],
) (K, V) {
	var zero V
	if k, v := LastOfSeq2(seq); v != zero {
		return k, v
	}
	if fallback != nil {
		return fallback()
	}
	var key K
	return key, zero
}

// LastOfFuncs returns the last T produced by a non-nil function in fns that
// is not a zero value of T. Note that because the given functions will be
// evaluated in reverse order, lower-indexed functions may not be called.
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

// LastOfFuncsOr returns the last T produced by a non-nil function in fns
// that is not a zero value of T, or otherwise returns the given fallback.
// Note that because the given functions will be evaluated in order,
// lower-indexed functions may not be called.
func LastOfFuncsOr[T comparable](fallback T, fns ...func() T) T {
	var zero T
	if x := LastOfFuncs(fns...); x != zero {
		return x
	}
	return fallback
}

// LastOfFuncsOrElse returns the last T produced by a non-nil function in fns
// that is not a zero value of T, or otherwise invokes fallback (if non-nil) to
// produce a return value. Note that because the given functions will be
// evaluated in order, lower-indexed functions may not be called.
func LastOfFuncsOrElse[T comparable](fallback func() T, fns ...func() T) T {
	var zero T
	if x := LastOfFuncs(fns...); x != zero {
		return x
	}
	if fallback != nil {
		return fallback()
	}
	return zero
}
