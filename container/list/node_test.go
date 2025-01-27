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
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/container/list"
)

func TestLinkWithTail(t *testing.T) {
	head, tail := list.LinkWithTail(1, 2, 3, 4, 5)
	require.Equal(t, 1, head.Value())

	var (
		cur     = head
		checked bool
	)
	for cur != nil {
		if cur.Next == nil {
			require.Equal(t, 5, tail.Value())
			checked = true
		}
		cur = cur.Next
	}
	require.True(t, checked)
}

func TestLinkToSlice(t *testing.T) {
	cases := map[string]struct {
		give *list.Node[int]
		want []int
	}{
		"forward": {
			give: list.Link(1, 2, 3, 4, 5),
			want: []int{1, 2, 3, 4, 5},
		},
		"backward": {
			give: list.Link(5, 4, 3, 2, 1),
			want: []int{5, 4, 3, 2, 1},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.give.ToSlice())
		})
	}
}

func TestNode_SetVsUnset(t *testing.T) {
	var node *list.Node[int]
	require.False(t, node.IsSet())

	node = list.NewNode(123)
	require.True(t, node.IsSet())
	node.Unset()
	require.False(t, node.IsSet())

	x, ok := node.Get()
	require.False(t, ok)
	require.Zero(t, x)
	require.Zero(t, node.Value())
	require.Nil(t, node.ToSlice())
	require.NotPanics(t, func() {
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
	node.Next = &list.Node[int]{}
	require.Equal(t, []int{123}, node.ToSlice())
	node.Next = list.NewNode(456)

	clear(calls)
	node.ForEachValue(func(have int) bool {
		calls[have]++
		return true
	})
	require.Equal(t, map[int]int{123: 1, 456: 1}, calls)
}

func TestNode_InsertBefore(t *testing.T) {
	cases := map[string]struct {
		head     *list.Node[int]
		node     *list.Node[int]
		wantHead []int
		wantNode []int
		skip     int
	}{
		"singles": {
			head:     list.Link(1),
			node:     list.Link(2),
			skip:     0,
			wantHead: []int{2, 1},
			wantNode: []int{2, 1},
		},
		"doubles prepend": {
			head:     list.Link(1, 2),
			node:     list.Link(3, 4),
			skip:     0,
			wantHead: []int{3, 4, 1, 2},
			wantNode: []int{3, 4, 1, 2},
		},
		"doubles mid": {
			head:     list.Link(1, 2),
			node:     list.Link(3, 4),
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

func TestNode_InsertValueBefore(t *testing.T) {
	var (
		head = list.Link(1, 2, 3, 4, 5)
		cur  = head
	)

	for cur != nil {
		cur.InsertValueBefore(-1)
		cur = cur.Next.Next
	}

	require.Equal(
		t,
		[]int{-1, 1, -1, 2, -1, 3, -1, 4, -1, 5},
		head.ToSlice(),
	)
}

func TestNode_InsertAfter(t *testing.T) {
	var (
		head = list.Link(1, 2, 3, 4, 5)
		cur  = head
	)

	for cur != nil {
		cur.InsertAfter(list.Link(-1, -1))
		cur = cur.Next.Next.Next
	}

	require.Equal(
		t,
		[]int{1, -1, -1, 2, -1, -1, 3, -1, -1, 4, -1, -1, 5, -1, -1},
		head.ToSlice(),
	)
}

func TestNode_InsertValueAfter(t *testing.T) {
	var (
		head = list.Link(1, 2, 3, 4, 5)
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

func TestNode_Append(t *testing.T) {
	head := list.Link(1, 2)
	for i := 0; i < 8; i += 2 {
		head.Append(list.Link(i+3, i+4))
	}
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, head.ToSlice())
}

func TestNode_AppendValue(t *testing.T) {
	var (
		head = list.Link(1)
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

func TestNode_Sort(t *testing.T) {
	head := list.Link(1, 5, 2, 4, 3)

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

func TestNode_Iter(t *testing.T) {
	head := list.Link(1, 2, 3, 4, 5)
	require.Equal(t, []int{1, 2, 3, 4, 5}, slices.Collect(head.Iter()))
	head = list.Link(5, 4, 3, 2, 1)
	require.Equal(t, []int{5, 4, 3, 2, 1}, slices.Collect(head.Iter()))

	next, stop := iter.Pull(list.Link(123, 456, 789).Iter())
	value, ok := next()
	require.True(t, ok)
	require.Equal(t, 123, value)
	stop()
	value, ok = next()
	require.False(t, ok)
	require.Zero(t, value)
}

func TestNode_ForEach(t *testing.T) {
	var (
		head = list.Link(1, 2, 3, 4, 5)
		want = []int{1, 2, 3, 4, 5}
		have = make([]int, 0, len(want))
	)

	head.ForEach(func(x *list.Node[int]) bool {
		have = append(have, x.Value())
		return true
	})
	require.Equal(t, want, have)

	have = have[:0]
	head.ForEach(func(x *list.Node[int]) bool {
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

func TestNode_PopMaybePop(t *testing.T) {
	var (
		head1 = list.Link(0, 1, 2, 3, 4)
		head2 = list.Link(0, 1, 2, 3, 4)
	)

	for want := range 5 {
		have1, ok1 := head1.MaybePop()
		require.True(t, ok1)
		require.Equal(t, want, have1)
		require.Equal(t, want, head2.Pop())
	}

	x, ok := head1.MaybePop()
	require.False(t, ok)
	require.Zero(t, x)
}

func BenchmarkNode_AppendValue(b *testing.B) {
	depths := []int{1, 2, 4, 8, 16, 32, 64}
	for _, depth := range depths {
		node := list.Link(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
		b.Run(strconv.Itoa(depth), func(b *testing.B) {
			for range b.N {
				x := node.Pop()
				node.AppendValue(x + 10)
			}
		})
	}
}
