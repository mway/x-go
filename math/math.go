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

// Package math provides math-related types and utilities.
package math

import (
	"golang.org/x/exp/constraints"
)

// Numeric describes any basic number type.
type Numeric interface {
	constraints.Integer | constraints.Float
}

// Signed describes the portion of [Numeric] that can hold negative values.
type Signed interface {
	constraints.Signed | constraints.Float
}

// Abs returns the absolute value of the given signed number.
func Abs[T Signed](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Mean returns a truncated average of all given numbers.
func Mean[T Numeric](x ...T) T {
	return T(MeanFloat64(x...))
}

// MeanFloat64 returns a precise average of all given numbers.
func MeanFloat64[T Numeric](x ...T) float64 {
	var (
		size  = len(x)
		total T
	)
	for len(x) >= 8 {
		total += x[0] + x[1] + x[2] + x[3] + x[4] + x[5] + x[6] + x[7]
		x = x[8:]
	}
	for i := range len(x) {
		total += x[i]
	}
	return float64(total) / float64(size)
}
