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

// Package limits provides helpers for dynamically determining numeric type limits.
package limits

import (
	"math"

	xmath "go.mway.dev/x/math"
)

// Max returns the maximum value that can be represented by T.
//
//nolint:gocyclo
func Max[T xmath.Numeric]() T {
	var x T
	switch any(x).(type) {
	case uint:
		var x uint = math.MaxUint
		return T(x)
	case uint8:
		var x uint8 = math.MaxUint8
		return T(x)
	case uint16:
		var x uint16 = math.MaxUint16
		return T(x)
	case uint32:
		var x uint32 = math.MaxUint32
		return T(x)
	case uint64:
		var x uint64 = math.MaxUint64
		return T(x)
	case int:
		x := math.MaxInt
		return T(x)
	case int8:
		var x int8 = math.MaxInt8
		return T(x)
	case int16:
		var x int16 = math.MaxInt16
		return T(x)
	case int32:
		var x int32 = math.MaxInt32
		return T(x)
	case int64:
		var x int64 = math.MaxInt64
		return T(x)
	case float32:
		var x float32 = math.MaxFloat32
		return T(x)
	default: // float64
		x := math.MaxFloat64
		return T(x)
	}
}

// Min returns the minimum value that can be represented by T.
func Min[T xmath.Numeric]() T {
	var x T
	switch any(x).(type) {
	case uint, uint8, uint16, uint32, uint64:
		return 0
	case int:
		x := math.MinInt
		return T(x)
	case int8:
		var x int8 = math.MinInt8
		return T(x)
	case int16:
		var x int16 = math.MinInt16
		return T(x)
	case int32:
		var x int32 = math.MinInt32
		return T(x)
	case int64:
		var x int64 = math.MinInt64
		return T(x)
	case float32:
		var x float32 = -math.MaxFloat32
		return T(x)
	default: // float64
		x := -math.MaxFloat64
		return T(x)
	}
}
