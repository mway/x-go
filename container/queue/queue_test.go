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

package queue_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/container/queue"
)

func TestNew(t *testing.T) {
	x := queue.New[int](32)
	require.Equal(t, 0, x.Len())
	require.Equal(t, 32, x.Cap())

	x.Push(123)
	require.Equal(t, 1, x.Len())
	require.Equal(t, 123, x.Pop())
	require.Equal(t, 0, x.Len())
	require.Equal(t, 32, x.Cap())
}

func TestQueue_PushTopPopLenCap(t *testing.T) {
	var x queue.Queue[int]
	require.Equal(t, 0, x.Len())
	x.Push(123)
	require.Equal(t, 1, x.Len())
	require.NotPanics(t, func() {
		require.Equal(t, 123, x.Front())
		require.Equal(t, 1, x.Len())
		require.Equal(t, 1, x.Cap())
		require.Equal(t, 123, x.Pop())
		require.Equal(t, 0, x.Len())
		require.Equal(t, 1, x.Cap())
	})
}

func TestQueue_MaybeTopMaybePop(t *testing.T) {
	var x queue.Queue[int]

	val, ok := x.MaybeFront()
	require.False(t, ok)
	require.Zero(t, val)

	val, ok = x.MaybePop()
	require.False(t, ok)
	require.Zero(t, val)

	x.Push(123)

	val, ok = x.MaybeFront()
	require.True(t, ok)
	require.Equal(t, 123, val)
	require.Equal(t, x.Front(), val)

	val, ok = x.MaybePop()
	require.True(t, ok)
	require.Equal(t, 123, val)
	require.Equal(t, 0, x.Len())
}

func TestQueue_PeekEachPopEach(t *testing.T) {
	var (
		give = []int{1, 2, 3, 4, 5}
		want = []int{1, 2, 3, 4, 5}
		have []int
		x    queue.Queue[int]
	)

	x.PeekEach(func(int) bool {
		require.FailNow(
			t,
			"an empty Queue[T].PeekEach should not invoke a callback",
		)
		return true
	})
	x.PopEach(func(int) bool {
		require.FailNow(
			t,
			"an empty Queue[T].PopEach should not invoke a callback",
		)
		return true
	})

	for _, n := range give {
		x.Push(n)
	}

	x.PeekEach(func(n int) bool {
		have = append(have, n)
		return true
	})
	require.Equal(t, want, have)

	have = have[:0]
	x.PopEach(func(n int) bool {
		have = append(have, n)
		return true
	})
	require.Equal(t, want, have)
}

func TestQueue_PeekEachPopEach_Abort(t *testing.T) {
	var (
		give = []int{1, 2, 3, 4, 5}
		want = []int{1}
		have []int
		x    queue.Queue[int]
	)

	for _, n := range give {
		x.Push(n)
	}

	require.Equal(t, 5, x.Len())
	x.PeekEach(func(n int) bool {
		have = append(have, n)
		return false
	})
	require.Equal(t, want, have)

	have = have[:0]
	require.Equal(t, 5, x.Len())
	x.PopEach(func(n int) bool {
		have = append(have, n)
		return false
	})
	require.Equal(t, want, have)
	require.Equal(t, 4, x.Len())
}

func BenchmarkQueue_PushPop(b *testing.B) {
	depths := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}
	for _, depth := range depths {
		var q queue.Queue[string]
		for range depth {
			q.Push(b.Name())
		}

		b.Run(fmt.Sprintf("depth %d", depth), func(b *testing.B) {
			for range b.N {
				q.Push(q.Pop())
			}
		})
	}
}
