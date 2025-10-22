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

package list_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mway.dev/x/container/list"
)

func TestLinkDoublyWithTail(t *testing.T) {
	head, tail := list.LinkDoublyWithTail(1, 2, 3, 4, 5)
	require.Equal(t, []int{1, 2, 3, 4, 5}, head.ToSlice())
	require.Equal(t, []int{5, 4, 3, 2, 1}, tail.ToSliceRev())
}

func TestLinkDoublyToSlice(t *testing.T) {
	cases := map[string]struct {
		give *list.DoubleNode[int]
		want []int
	}{
		"forward": {
			give: list.LinkDoubly(1, 2, 3, 4, 5),
			want: []int{1, 2, 3, 4, 5},
		},
		"backward": {
			give: list.LinkDoubly(5, 4, 3, 2, 1),
			want: []int{5, 4, 3, 2, 1},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.give.ToSlice())
			var zero list.DoubleNode[int]
			require.Equal(t, []int(nil), zero.ToSlice())
			require.Equal(t, []int(nil), zero.ToSliceRev())

			slices.Reverse(tt.want)

			last := tt.give
			for last.Next != nil {
				last = last.Next
			}
			require.Equal(t, tt.want, last.ToSliceRev())
		})
	}
}

func TestDoubleNode_WithNextPrev(t *testing.T) {
	node := list.NewDoubleNode(456).
		WithNext(list.NewDoubleNode(789).WithNext(&list.DoubleNode[int]{})).
		WithPrev(list.NewDoubleNode(123).WithPrev(&list.DoubleNode[int]{}))
	require.Equal(t, []int{456, 789}, node.ToSlice())
	require.Equal(t, []int{456, 123}, node.ToSliceRev())
	require.Equal(t, []int{123, 456, 789}, node.Prev.ToSlice())
	require.Equal(t, []int{123}, node.Prev.ToSliceRev())
	require.Equal(t, []int{789}, node.Next.ToSlice())
	require.Equal(t, []int{789, 456, 123}, node.Next.ToSliceRev())
}

func TestDoubleNode_SetVsUnset(t *testing.T) {
	var node *list.DoubleNode[int]
	require.False(t, node.IsSet())

	node = list.NewDoubleNode(123)
	require.True(t, node.IsSet())
	node.Unset()
	require.False(t, node.IsSet())

	x, ok := node.Get()
	require.False(t, ok)
	require.Zero(t, x)
	require.Zero(t, node.Value())
	require.Nil(t, node.ToSlice())
	require.NotPanics(t, func() {
		node.DetachNext()
		node.DetachPrev()
		node.Sort(nil)
	})

	calls := make(map[int]int)
	node.ForEachValue(func(have int) bool {
		calls[have]++
		return true
	})
	require.Len(t, calls, 0)

	collected := slices.Collect(node.Iter())
	require.Nil(t, collected)

	node.Set(123)
	require.True(t, node.IsSet())
	x, ok = node.Get()
	require.True(t, ok)
	require.Equal(t, 123, x)
	require.Equal(t, 123, node.Value())
	require.Equal(t, []int{123}, node.ToSlice())
	node.Next = &list.DoubleNode[int]{}
	require.Equal(t, []int{123}, node.ToSlice())
	node.InsertAfter(list.NewDoubleNode(456))

	clear(calls)
	node.ForEachValue(func(have int) bool {
		calls[have]++
		return true
	})
	require.Equal(t, map[int]int{123: 1, 456: 1}, calls)
}

func TestDoubleNode_InsertBefore(t *testing.T) {
	cases := map[string]struct {
		head     *list.DoubleNode[int]
		node     *list.DoubleNode[int]
		wantHead []int
		wantNode []int
		skip     int
	}{
		"singles": {
			head:     list.LinkDoubly(1),
			node:     list.LinkDoubly(2),
			skip:     0,
			wantHead: []int{1},
			wantNode: []int{2, 1},
		},
		"doubles prepend": {
			head:     list.LinkDoubly(1, 2),
			node:     list.LinkDoubly(3, 4),
			skip:     0,
			wantHead: []int{1, 2},
			wantNode: []int{3, 4, 1, 2},
		},
		"doubles mid": {
			head:     list.LinkDoubly(1, 2),
			node:     list.LinkDoubly(3, 4),
			skip:     1,
			wantHead: []int{1, 3, 4, 2},
			wantNode: []int{3, 4, 2},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			cur := tt.head
			for range tt.skip {
				cur = cur.Next
			}
			cur.InsertBefore(tt.node)
			require.Equal(t, tt.wantHead, tt.head.ToSlice())
			require.Equal(t, tt.wantNode, tt.node.ToSlice())
		})
	}
}

func TestDoubleNode_InsertValueBefore(t *testing.T) {
	var (
		head = list.LinkDoubly(1, 2, 3, 4, 5)
		cur  = head
	)

	for cur != nil {
		cur.InsertValueBefore(-1)
		cur = cur.Next
	}

	for head.Prev != nil {
		head = head.Prev
	}

	require.Equal(
		t,
		[]int{-1, 1, -1, 2, -1, 3, -1, 4, -1, 5},
		head.ToSlice(),
	)
}

func TestDoubleNode_InsertAfter(t *testing.T) {
	var (
		head = list.LinkDoubly(1, 2, 3, 4, 5)
		cur  = head
	)

	for cur != nil {
		cur.InsertAfter(list.LinkDoubly(-1, -1))
		cur = cur.Next.Next.Next
	}

	require.Equal(
		t,
		[]int{1, -1, -1, 2, -1, -1, 3, -1, -1, 4, -1, -1, 5, -1, -1},
		head.ToSlice(),
	)
}

func TestDoubleNode_InsertValueAfter(t *testing.T) {
	var (
		head = list.LinkDoubly(1, 2, 3, 4, 5)
		cur  = head
	)

	for cur != nil {
		cur.InsertValueAfter(-1)
		cur = cur.Next.Next
	}

	require.Equal(
		t,
		[]int{1, -1, 2, -1, 3, -1, 4, -1, 5, -1},
		head.ToSlice(),
	)
}

func TestDoubleNode_Append(t *testing.T) {
	head := list.LinkDoubly(1, 2)
	for i := 0; i < 8; i += 2 {
		head.Append(list.LinkDoubly(i+3, i+4))
	}
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, head.ToSlice())
}

func TestDoubleNode_AppendValue(t *testing.T) {
	var (
		head = list.LinkDoubly(1)
		cur  = head
	)
	for i := range 4 {
		head.AppendValue(i + 2)
		cur = cur.Next
	}
	for i := range 5 {
		cur.AppendValue(i + 6)
		cur = cur.Next
	}
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, head.ToSlice())
}

func TestDoubleNode_Sort(t *testing.T) {
	head := list.LinkDoubly(1, 5, 2, 4, 3)

	head.Sort(func(a int, b int) int {
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	})
	require.Equal(t, []int{1, 2, 3, 4, 5}, head.ToSlice())

	head.Sort(func(a int, b int) int {
		switch {
		case a < b:
			return 1
		case a > b:
			return -1
		default:
			return 0
		}
	})
	require.Equal(t, []int{5, 4, 3, 2, 1}, head.ToSlice())
}

func TestDoubleNode_Iter(t *testing.T) {
	head := list.LinkDoubly(1, 2, 3, 4, 5)
	require.Equal(t, []int{1, 2, 3, 4, 5}, slices.Collect(head.Iter()))
	head = list.LinkDoubly(5, 4, 3, 2, 1)
	require.Equal(t, []int{5, 4, 3, 2, 1}, slices.Collect(head.Iter()))

	next, stop := iter.Pull(list.LinkDoubly(123, 456, 789).Iter())
	value, ok := next()
	require.True(t, ok)
	require.Equal(t, 123, value)
	stop()
	value, ok = next()
	require.False(t, ok)
	require.Zero(t, value)
}

func TestDoubleNode_ForEach(t *testing.T) {
	var (
		head = list.LinkDoubly(1, 2, 3, 4, 5)
		want = []int{1, 2, 3, 4, 5}
		have = make([]int, 0, len(want))
	)

	head.ForEach(func(x *list.DoubleNode[int]) bool {
		have = append(have, x.Value())
		return true
	})
	require.Equal(t, want, have)

	have = have[:0]
	head.ForEach(func(x *list.DoubleNode[int]) bool {
		have = append(have, x.Value())
		return false
	})
	require.Equal(t, want[:1], have)

	have = have[:0]
	head.ForEachValue(func(x int) bool {
		have = append(have, x)
		return true
	})
	require.Equal(t, want, have)

	have = have[:0]
	head.ForEachValue(func(x int) bool {
		have = append(have, x)
		return false
	})
	require.Equal(t, want[:1], have)
}

func TestDoubleNode_DetachPrev(t *testing.T) {
	var (
		head = list.LinkDoubly(1, 2, 3, 4, 5, 6)
		cur  = head
	)
	for range 3 {
		cur = cur.Next
	}
	cur.DetachPrev()
	require.Nil(t, cur.Prev)
	require.Equal(t, []int{1, 2, 3}, head.ToSlice())
	require.Equal(t, []int{4, 5, 6}, cur.ToSlice())
}

func TestDoubleNode_DetachNext(t *testing.T) {
	var (
		head = list.LinkDoubly(1, 2, 3, 4, 5, 6)
		tail = head
	)
	for tail.Next != nil {
		tail = tail.Next
	}

	cur := tail
	for range 3 {
		cur = cur.Prev
	}
	cur.DetachNext()
	require.Nil(t, cur.Next)
	require.Equal(t, []int{6, 5, 4}, tail.ToSliceRev())
	require.Equal(t, []int{3, 2, 1}, cur.ToSliceRev())
}
