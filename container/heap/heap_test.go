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

package heap_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mway.dev/x/container/heap"
)

func TestHeap(t *testing.T) {
	var (
		minh = heap.NewMinHeap(2, 4, 5, 1, 3)
		maxh = heap.NewMaxHeap(2, 4, 5, 1, 3)
	)

	require.Equal(t, 5, minh.Len())
	require.Equal(t, 1, minh.Min())
	require.Equal(t, 5, maxh.Len())
	require.Equal(t, 5, maxh.Max())

	minh.Reset()
	maxh.Reset()
	require.Equal(t, 0, minh.Min())
	require.Equal(t, 0, maxh.Max())
	require.Equal(t, 0, minh.Len())

	const times = 100
	for i := times; i > 0; i-- {
		minh.Push(i)
		require.Equal(t, i, minh.Min())
		maxh.Push(i)
		require.Equal(t, times, maxh.Max())
	}
	require.Equal(t, times, minh.Len())
	require.Equal(t, times, maxh.Len())

	for i := 1; i <= times; i++ {
		minh.Push(i)
		require.Equal(t, 1, minh.Min())
		maxh.Push(i)
		require.Equal(t, times, maxh.Max())
	}
	require.Equal(t, 2*times, minh.Len())
	require.Equal(t, 2*times, maxh.Len())

	var (
		lastMin = math.MinInt
		lastMax = math.MaxInt
	)

	for minh.Len() > 2 {
		var (
			curMin = minh.Pop()
			curMax = maxh.Pop()
		)
		require.Greater(t, curMin, lastMin)
		require.Equal(t, curMin, minh.Remove(0))
		lastMin = curMin
		require.Less(t, curMax, lastMax)
		require.Equal(t, curMax, maxh.Remove(0))
		lastMax = curMax
	}

	require.Equal(t, 100, minh.Min())
	require.Equal(t, 100, minh.Pop())
	require.Equal(t, 100, minh.Remove(0))
	require.Equal(t, 0, minh.Min())
	require.Equal(t, 0, minh.Len())

	require.Equal(t, 1, maxh.Max())
	require.Equal(t, 1, maxh.Pop())
	require.Equal(t, 1, maxh.Remove(0))
	require.Equal(t, 0, maxh.Max())
	require.Equal(t, 0, maxh.Len())
}

func BenchmarkMinHeap_PushPop(b *testing.B) {
	var h heap.MinHeap[int]

	for i := range b.N {
		h.Push(i)
		h.Pop()
	}
}
