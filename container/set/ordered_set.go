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
	"sort"

	"golang.org/x/exp/maps"
)

// An OrderedSet is a collection of unique values of type T where the order of
// additions is remembered.
type OrderedSet[T comparable] struct {
	data map[T]int // value -> index
}

// NewOrdered creates a new [OrderedSet[T]] containing the given values.
func NewOrdered[T comparable](values ...T) OrderedSet[T] {
	data := make(map[T]int, len(values))
	for i, value := range values {
		data[value] = i
	}
	return OrderedSet[T]{
		data: data,
	}
}

// Add adds value to the set if it is not present, returning whether the value
// was added.
func (s *OrderedSet[T]) Add(value T) bool {
	if s.data == nil {
		s.data = map[T]int{
			value: 0,
		}
		return true
	} else if _, ok := s.data[value]; ok {
		return false
	}

	s.data[value] = len(s.data)
	return true
}

// AddN adds each of the given values to the set if they are not present,
// returning the number of values added.
func (s *OrderedSet[T]) AddN(values ...T) (added int) {
	for _, value := range values {
		if ok := s.Add(value); ok {
			added++
		}
	}
	return
}

// AddSet adds each of the values in other to the set if they are not present,
// returning the number of values added.
func (s *OrderedSet[T]) AddSet(other Set[T]) (added int) {
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
func (s *OrderedSet[T]) AddOrderedSet(other OrderedSet[T]) (added int) {
	other.ForEach(func(value T) bool {
		if s.Add(value) {
			added++
		}
		return true
	})
	return
}

// Clear resets the set, removing all data.
func (s OrderedSet[T]) Clear() {
	clear(s.data)
}

// Contains indicates if the set contains the given value.
func (s OrderedSet[T]) Contains(value T) bool {
	_, ok := s.data[value]
	return ok
}

// ContainsAny indicates if the set contains any of the given values.
func (s OrderedSet[T]) ContainsAny(values ...T) bool {
	return s.count(values) > 0
}

// ContainsAll indicates if the set contains all of the given values.
func (s OrderedSet[T]) ContainsAll(values ...T) bool {
	return len(values) > 0 && s.count(values) == len(values)
}

// ForEach invokes the given callback for each element in the set. Iteration
// order is deterministic.
func (s OrderedSet[T]) ForEach(fn Callback[T]) {
	if len(s.data) == 0 {
		return
	}

	type pair struct {
		key   T
		value int
	}

	pairs := make([]pair, 0, len(s.data))
	for key, value := range s.data {
		pairs = append(pairs, pair{
			key:   key,
			value: value,
		})
	}
	sort.Slice(pairs, func(i int, j int) bool {
		return pairs[i].value < pairs[j].value
	})

	for _, pair := range pairs {
		if !fn(pair.key) {
			break
		}
	}
}

// Intersect returns an [OrderedSet[T]] containing the intersection between s
// and other.
func (s OrderedSet[T]) Intersect(
	other OrderedSet[T],
) (result OrderedSet[T]) {
	if len(s.data) == 0 || len(other.data) == 0 {
		return
	}

	s.ForEach(func(k T) bool {
		if other.Contains(k) {
			if result.data == nil {
				result = NewOrdered(k)
			} else {
				result.Add(k)
			}
		}
		return true
	})
	return
}

// UnorderedIntersect returns a [Set[T]] containing the intersection between s
// and other.
func (s OrderedSet[T]) UnorderedIntersect(other Set[T]) (result Set[T]) {
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

// Merge returns an [OrderedSet[T]] containing all unique values between s and
// other.
func (s OrderedSet[T]) Merge(other OrderedSet[T]) (result OrderedSet[T]) {
	switch {
	case len(s.data) == 0:
		return other
	case len(other.data) == 0:
		return s
	default:
	}

	result.data = maps.Clone(s.data)
	other.ForEach(func(value T) bool {
		result.Add(value)
		return true
	})

	return
}

// Len returns the number of values held in the set.
func (s OrderedSet[T]) Len() int {
	return len(s.data)
}

// ToSet converts the set to an unordered [Set[T]].
func (s OrderedSet[T]) ToSet() Set[T] {
	if len(s.data) == 0 {
		return Set[T]{}
	}

	return New(maps.Keys(s.data)...)
}

// ToSlice returns the set as a slice. Note that the order of elements is not
// guaranteed.
func (s OrderedSet[T]) ToSlice() []T {
	if len(s.data) == 0 {
		return nil
	}

	type pair struct {
		key   T
		value int
	}

	pairs := make([]pair, 0, len(s.data))
	for key, value := range s.data {
		pairs = append(pairs, pair{
			key:   key,
			value: value,
		})
	}
	sort.Slice(pairs, func(i int, j int) bool {
		return pairs[i].value < pairs[j].value
	})

	keys := make([]T, len(pairs))
	for i, pair := range pairs {
		keys[i] = pair.key
	}
	return keys
}

// UnorderedUnion returns a [Set[T]] containing the union between s and other.
func (s OrderedSet[T]) UnorderedUnion(other Set[T]) (result Set[T]) {
	switch {
	case len(s.data) == 0:
		return other
	case len(other.data) == 0:
		return s.ToSet()
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

// Union returns an [OrderedSet[T]] containing the union between s and other.
func (s OrderedSet[T]) Union(other OrderedSet[T]) (result OrderedSet[T]) {
	switch {
	case len(s.data) == 0:
		return other
	case len(other.data) == 0:
		return s
	default:
	}

	s.ForEach(func(k T) bool {
		if !other.Contains(k) {
			if result.data == nil {
				result = NewOrdered(k)
			} else {
				result.Add(k)
			}
		}
		return true
	})

	other.ForEach(func(k T) bool {
		if !s.Contains(k) {
			if result.data == nil {
				result = NewOrdered(k)
			} else {
				result.Add(k)
			}
		}
		return true
	})

	return
}

func (s OrderedSet[T]) count(values []T) (found int) {
	for _, value := range values {
		if _, ok := s.data[value]; ok {
			found++
		}
	}
	return
}
