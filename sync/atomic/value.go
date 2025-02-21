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

// Package atomic provides atomicity-related types and helpers.
package atomic

import (
	"sync/atomic"
)

// Value is a strongly-typed atomic value. It is otherwise identical to the
// standard library's atomic.Value.
type Value[T any] struct {
	value atomic.Value
}

// NewValue creates a new Value that can store values of type T, initializing
// the Value with the given initial value.
func NewValue[T any](initial T) *Value[T] {
	v := &Value[T]{}
	v.value.Store(initial)
	return v
}

// CompareAndSwap performs an atomic compare and swap using oldval and newval.
// The return value indicates whether a swap took place.
func (v *Value[T]) CompareAndSwap(oldval T, newval T) bool {
	return v.value.CompareAndSwap(oldval, newval)
}

// Load loads the currently held T value.
func (v *Value[T]) Load() T {
	return v.value.Load().(T) //nolint:errcheck
}

// Store stores the given T value.
func (v *Value[T]) Store(val T) {
	v.value.Store(val)
}

// Swap swaps the currently held T value with the given value, returning the
// previous T.
func (v *Value[T]) Swap(val T) T {
	return v.value.Swap(val).(T) //nolint:errcheck
}
