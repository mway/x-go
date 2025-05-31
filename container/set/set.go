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

// Package set provides helpers for working with sets.
package set

import (
	"maps"
	"slices"
)

// A Set is a collection of unique values of type T.
type Set[T comparable] struct {
	data map[T]struct{}
}

// A Callback handles a value of type T during iteration and returns whether
// iteration should continue.
type Callback[T comparable] func(T) bool

// New creates a new [Set[T]] containing the given values.
func New[T comparable](values ...T) Set[T] {
	data := make(map[T]struct{}, len(values))
	for _, value := range values {
		data[value] = struct{}{}
	}
	return Set[T]{
		data: data,
	}
}

// Add adds value to the set if it is not present, returning whether the value
// was added.
func (s *Set[T]) Add(value T) bool {
	if s.data == nil {
		s.data = map[T]struct{}{
			value: {},
		}
		return true
	} else if _, ok := s.data[value]; ok {
		return false
	}

	s.data[value] = struct{}{}
	return true
}

// AddN adds each of the given values to the set if they are not present,
// returning the number of values added.
func (s *Set[T]) AddN(values ...T) (added int) {
	for _, value := range values {
		if ok := s.Add(value); ok {
			added++
		}
	}
	return
}

// AddSet adds each of the values in other to the set if they are not present,
// returning the number of values added.
func (s *Set[T]) AddSet(other Set[T]) (added int) {
	other.ForEach(func(value T) bool {
		if s.Add(value) {
			added++
		}
		return true
	})
	return
}

// AddOrderedSet adds each of the values in other to the set if they are not
// present, returning the number of values added.
func (s *Set[T]) AddOrderedSet(other OrderedSet[T]) (added int) {
	other.ForEach(func(value T) bool {
		if s.Add(value) {
			added++
		}
		return true
	})
	return
}

// Clear resets the set, removing all data.
func (s Set[T]) Clear() {
	clear(s.data)
}

// Contains indicates if the set contains the given value.
func (s Set[T]) Contains(value T) bool {
	_, ok := s.data[value]
	return ok
}

// ContainsAny indicates if the set contains any of the given values.
func (s Set[T]) ContainsAny(values ...T) bool {
	return s.count(values) > 0
}

// ContainsAll indicates if the set contains all of the given values.
func (s Set[T]) ContainsAll(values ...T) bool {
	return len(values) > 0 && s.count(values) == len(values)
}

// ForEach invokes the given callback for each element in the set. Iteration
// order is not deterministic.
func (s Set[T]) ForEach(fn Callback[T]) {
	for k := range s.data {
		if !fn(k) {
			break
		}
	}
}

// Intersect returns a [Set[T]] containing the intersection between s and
// other.
func (s Set[T]) Intersect(other Set[T]) (result Set[T]) {
	if len(s.data) == 0 || len(other.data) == 0 {
		return
	}

	for k := range s.data {
		if other.Contains(k) {
			if result.data == nil {
				result = New(k)
			} else {
				result.Add(k)
			}
		}
	}

	return result
}

// OrderedIntersect returns a [Set[T]] containing the intersection between s
// and other.
func (s Set[T]) OrderedIntersect(other OrderedSet[T]) (result Set[T]) {
	if len(s.data) == 0 || len(other.data) == 0 {
		return
	}

	for k := range s.data {
		if other.Contains(k) {
			if result.data == nil {
				result = New(k)
			} else {
				result.Add(k)
			}
		}
	}

	return result
}

// Merge returns a [Set[T]] containing all unique values between s and other.
func (s Set[T]) Merge(other Set[T]) (result Set[T]) {
	switch {
	case len(s.data) == 0:
		return other
	case len(other.data) == 0:
		return s
	default:
	}

	result.data = maps.Clone(s.data)
	for k := range other.data {
		result.Add(k)
	}

	return
}

// Len returns the number of values held in the set.
func (s Set[T]) Len() int {
	return len(s.data)
}

// ToSlice returns the set as a slice. Note that the order of elements is not
// guaranteed.
func (s Set[T]) ToSlice() []T {
	if len(s.data) == 0 {
		return nil
	}
	return slices.Collect(maps.Keys(s.data))
}

// Union returns a [Set[T]] containing the union between s and other.
func (s Set[T]) Union(other Set[T]) (result Set[T]) {
	switch {
	case len(s.data) == 0:
		return other
	case len(other.data) == 0:
		return s
	default:
	}

	for k := range s.data {
		if !other.Contains(k) {
			if result.data == nil {
				result = New(k)
			} else {
				result.Add(k)
			}
		}
	}

	for k := range other.data {
		if !s.Contains(k) {
			if result.data == nil {
				result = New(k)
			} else {
				result.Add(k)
			}
		}
	}

	return result
}

// OrderedUnion returns a [Set[T]] containing the union between s and other.
func (s Set[T]) OrderedUnion(other OrderedSet[T]) (result Set[T]) {
	switch {
	case len(s.data) == 0:
		return other.ToSet()
	case len(other.data) == 0:
		return s
	default:
	}

	for k := range s.data {
		if !other.Contains(k) {
			if result.data == nil {
				result = New(k)
			} else {
				result.Add(k)
			}
		}
	}

	for k := range other.data {
		if !s.Contains(k) {
			if result.data == nil {
				result = New(k)
			} else {
				result.Add(k)
			}
		}
	}

	return
}

func (s Set[T]) count(values []T) int {
	var found int
	for _, value := range values {
		if _, ok := s.data[value]; ok {
			found++
		}
	}
	return found
}
