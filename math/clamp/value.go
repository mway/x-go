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

package clamp

import (
	xmath "go.mway.dev/x/math"
)

// A Value is a [Range[T]]-bounded value.
type Value[T xmath.Numeric] struct {
	cur    T
	bounds Range[T]
}

// NewValue creates a new [Value[T]] holding value as clamped to [lower,upper].
// The resulting [Value[T]] will continue to constrain any value to this range.
func NewValue[T xmath.Numeric](value T, lower T, upper T) Value[T] {
	r := NewRange(lower, upper)
	return Value[T]{
		cur:    r.Clamp(value),
		bounds: r,
	}
}

// NewValueWithRange creates a new [Value[T]] holding value as clamped to the
// given range. The resulting [Value[T]] will continue to constrain any value
// to this range.
func NewValueWithRange[T xmath.Numeric](value T, bounds Range[T]) Value[T] {
	return Value[T]{
		cur:    bounds.Clamp(value),
		bounds: bounds,
	}
}

// InRange indicates whether value is in v's range.
func (v *Value[T]) InRange(value T) bool {
	return v.bounds.InRange(value)
}

// Min returns the minimum bound (inclusive) for v's range.
func (v *Value[T]) Min() T {
	return v.bounds.Min
}

// SetMin sets the minimum bound (inclusive) for v's range, clamping the
// current value if necessary.
func (v *Value[T]) SetMin(value T) {
	v.bounds.Min = value
	v.Set(v.cur)
}

// Max returns the maximum bound (inclusive) for v's range.
func (v *Value[T]) Max() T {
	return v.bounds.Max
}

// SetMax sets the maximum bound (inclusive) for v's range, clamping the
// current value if necessary.
func (v *Value[T]) SetMax(value T) {
	v.bounds.Max = value
	v.Set(v.cur)
}

// Value returns the current underlying value of v.
func (v *Value[T]) Value() T {
	return v.cur
}

// Set sets v to value, clamping within v's range.
func (v *Value[T]) Set(value T) T {
	v.cur = v.bounds.Clamp(value)
	return v.cur
}

// Add adds delta to v, clamping within v's range.
func (v *Value[T]) Add(delta T) T {
	v.cur = v.bounds.Add(v.cur, delta)
	return v.cur
}

// Sub subtracts delta from v, clamping within v's range.
func (v *Value[T]) Sub(delta T) T {
	v.cur = v.bounds.Sub(v.cur, delta)
	return v.cur
}
