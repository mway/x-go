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
	"fmt"

	xmath "go.mway.dev/x/math"
)

// A Range describes an inclusive [min,max] range.
type Range[T xmath.Numeric] struct {
	// Min is the minimum value allowed in this range. Lesser values will be
	// clamped to this value.
	Min T
	// Max is the maximum value allowed in this range. Greater values will be
	// clamped to this value.
	Max T
}

// NewRange creates a new [Range[T]] with the given min and max values.
func NewRange[T xmath.Numeric](lo T, hi T) Range[T] {
	if lo > hi {
		panic(fmt.Sprintf(
			"clamp.NewRange: lower (%v) > upper (%v)",
			lo,
			hi,
		))
	}

	return Range[T]{
		Min: lo,
		Max: hi,
	}
}

// InRange indicates if value is within [r.Min,r.Max].
func (r Range[T]) InRange(value T) bool {
	return value >= r.Min && value <= r.Max
}

// Clamp returns value clamped to the range [r.Min,r.Max].
func (r Range[T]) Clamp(value T) T {
	return Clamp(value, r.Min, r.Max)
}

// Add adds delta to value, clamping if necessary, and returns the result.
func (r Range[T]) Add(value T, delta T) T {
	return Add(value, delta, r.Min, r.Max)
}

// Sub subtracts delta from value, clamping if necessary, and returns the
// result.
func (r Range[T]) Sub(value T, delta T) T {
	return Sub(value, delta, r.Min, r.Max)
}
