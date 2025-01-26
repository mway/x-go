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

package limits_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/math/limits"
)

func TestMax(t *testing.T) {
	require.Equal(t, uint(math.MaxUint), limits.Max[uint]())
	require.Equal(t, uint8(math.MaxUint8), limits.Max[uint8]())
	require.Equal(t, uint16(math.MaxUint16), limits.Max[uint16]())
	require.Equal(t, uint32(math.MaxUint32), limits.Max[uint32]())
	require.Equal(t, uint64(math.MaxUint64), limits.Max[uint64]())
	require.Equal(t, int(math.MaxInt), limits.Max[int]())
	require.Equal(t, int8(math.MaxInt8), limits.Max[int8]())
	require.Equal(t, int16(math.MaxInt16), limits.Max[int16]())
	require.Equal(t, int32(math.MaxInt32), limits.Max[int32]())
	require.Equal(t, int64(math.MaxInt64), limits.Max[int64]())
	require.Equal(t, float32(math.MaxFloat32), limits.Max[float32]())
	require.Equal(t, math.MaxFloat64, limits.Max[float64]())
}

func TestMin(t *testing.T) {
	require.Equal(t, uint(0), limits.Min[uint]())
	require.Equal(t, uint8(0), limits.Min[uint8]())
	require.Equal(t, uint16(0), limits.Min[uint16]())
	require.Equal(t, uint32(0), limits.Min[uint32]())
	require.Equal(t, uint64(0), limits.Min[uint64]())
	require.Equal(t, int(math.MinInt), limits.Min[int]())
	require.Equal(t, int8(math.MinInt8), limits.Min[int8]())
	require.Equal(t, int16(math.MinInt16), limits.Min[int16]())
	require.Equal(t, int32(math.MinInt32), limits.Min[int32]())
	require.Equal(t, int64(math.MinInt64), limits.Min[int64]())
	require.Equal(t, -float32(math.MaxFloat32), limits.Min[float32]())
	require.Equal(t, -math.MaxFloat64, limits.Min[float64]())
}

func BenchmarkMax(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = limits.Max[int]()
		}
	})
}

func BenchmarkMin(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = limits.Min[int]()
		}
	})
}
