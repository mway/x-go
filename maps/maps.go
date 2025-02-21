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

// Package maps provides slice-related utilities.
package maps

import (
	"iter"
)

// Predicate constrains the types of predicate functions supported by [Filter].
type Predicate[K comparable, V any] interface {
	~func(K, V) bool | ~func(K, struct{}) bool | ~func(struct{}, V) bool
}

// ByKey returns fn as a [Predicate] in the form of func(K, struct{}).
func ByKey[K comparable](fn func(K) bool) func(K, struct{}) bool {
	return func(k K, _ struct{}) bool {
		return fn(k)
	}
}

// ByValue returns fn as a [Predicate] in the form of func(struct{}, V).
func ByValue[V any](fn func(V) bool) func(struct{}, V) bool {
	return func(_ struct{}, v V) bool {
		return fn(v)
	}
}

// Filter filters src, returning a copy containing all elements for which the
// given [Predicate] evaluates to true.
func Filter[K comparable, V any, M ~map[K]V, P Predicate[K, V]](
	src M,
	pred P,
) M {
	var (
		dst = make(M, len(src))
		fn  = wrap[K, V](pred)
	)

	for k, v := range src {
		if fn(k, v) {
			dst[k] = v
		}
	}

	if len(dst) == 0 {
		dst = nil
	}

	return dst
}

// Transform transforms src from type In to Out using the provided mapper.
func Transform[
	Out ~map[OutK]OutV, OutK comparable, OutV any,
	In ~map[InK]InV, InK comparable, InV any,
	Mapper func(InK, InV) (OutK, OutV),
](src In, mapper Mapper) Out {
	dst := make(Out, len(src))
	for ink, inv := range src {
		outk, outv := mapper(ink, inv)
		dst[outk] = outv
	}
	return dst
}

// Iter returns an [iter.Seq[V]] that ranges over s.
func Iter[T ~map[K]V, K comparable, V any](s T) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

// Iter2 returns an [iter.Seq2[K, V]] that ranges over s.
func Iter2[T ~map[K]V, K comparable, V any](m T) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func wrap[K comparable, V any, P Predicate[K, V]](pred P) func(K, V) bool {
	switch pred := any(pred).(type) {
	case func(K, V) bool:
		return pred
	case func(K, struct{}) bool:
		return func(k K, _ V) bool {
			return pred(k, struct{}{})
		}
	default:
		// n.b. This conversion cannot fail as the type of P is restricted to
		//      the two switch cases above and the one type below.
		x := pred.(func(struct{}, V) bool) //nolint:errcheck
		return func(_ K, v V) bool {
			return x(struct{}{}, v)
		}
	}
}
