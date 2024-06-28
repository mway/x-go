// Copyright (c) 2024 Matt Way
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

package set_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/container/set"
)

func TestOrderedSet_Add(t *testing.T) {
	var x set.OrderedSet[int]
	require.True(t, x.Add(1))

	x = set.NewOrdered[int]()
	require.True(t, x.Add(1))
	require.False(t, x.Add(1))
	require.True(t, x.Add(2))
}

func TestOrderedSet_AddN(t *testing.T) {
	cases := map[string]struct {
		base      set.OrderedSet[int]
		give      []int
		want      []int
		wantAdded int
	}{
		"zero base": {
			base:      set.OrderedSet[int]{},
			give:      []int{1, 2, 3},
			want:      []int{1, 2, 3},
			wantAdded: 3,
		},
		"empty base": {
			base:      set.NewOrdered[int](),
			give:      []int{1, 2, 3},
			want:      []int{1, 2, 3},
			wantAdded: 3,
		},
		"no overlap": {
			base:      set.NewOrdered(1, 2, 3),
			give:      []int{4, 5, 6},
			want:      []int{1, 2, 3, 4, 5, 6},
			wantAdded: 3,
		},
		"some overlap": {
			base:      set.NewOrdered(1, 2, 3),
			give:      []int{2, 3, 4},
			want:      []int{1, 2, 3, 4},
			wantAdded: 1,
		},
		"total overlap": {
			base:      set.NewOrdered(1, 2, 3),
			give:      []int{1, 2, 3},
			want:      []int{1, 2, 3},
			wantAdded: 0,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			haveAdded := tt.base.AddN(tt.give...)
			require.Equal(t, tt.wantAdded, haveAdded)
			require.Equal(t, tt.want, tt.base.ToSlice())
		})
	}
}

func TestOrderedSet_AddSet(t *testing.T) {
	cases := map[string]struct {
		base      set.OrderedSet[int]
		give      set.Set[int]
		want      []int
		wantAdded int
	}{
		"zero base": {
			base:      set.OrderedSet[int]{},
			give:      set.New(1, 2, 3),
			want:      []int{1, 2, 3},
			wantAdded: 3,
		},
		"empty base": {
			base:      set.NewOrdered[int](),
			give:      set.New(1, 2, 3),
			want:      []int{1, 2, 3},
			wantAdded: 3,
		},
		"no overlap": {
			base:      set.NewOrdered(1, 2, 3),
			give:      set.New(4, 5, 6),
			want:      []int{1, 2, 3, 4, 5, 6},
			wantAdded: 3,
		},
		"some overlap": {
			base:      set.NewOrdered(1, 2, 3),
			give:      set.New(2, 3, 4),
			want:      []int{1, 2, 3, 4},
			wantAdded: 1,
		},
		"total overlap": {
			base:      set.NewOrdered(1, 2, 3),
			give:      set.New(1, 2, 3),
			want:      []int{1, 2, 3},
			wantAdded: 0,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			haveAdded := tt.base.AddSet(tt.give)
			require.Equal(t, tt.wantAdded, haveAdded)
			require.ElementsMatch(t, tt.want, tt.base.ToSlice())
		})
	}
}

func TestOrderedSet_AddOrderedSet(t *testing.T) {
	cases := map[string]struct {
		base      set.OrderedSet[int]
		give      set.OrderedSet[int]
		want      []int
		wantAdded int
	}{
		"zero base": {
			base:      set.OrderedSet[int]{},
			give:      set.NewOrdered(1, 2, 3),
			want:      []int{1, 2, 3},
			wantAdded: 3,
		},
		"empty base": {
			base:      set.NewOrdered[int](),
			give:      set.NewOrdered(1, 2, 3),
			want:      []int{1, 2, 3},
			wantAdded: 3,
		},
		"no overlap": {
			base:      set.NewOrdered(1, 2, 3),
			give:      set.NewOrdered(4, 5, 6),
			want:      []int{1, 2, 3, 4, 5, 6},
			wantAdded: 3,
		},
		"some overlap": {
			base:      set.NewOrdered(1, 2, 3),
			give:      set.NewOrdered(2, 3, 4),
			want:      []int{1, 2, 3, 4},
			wantAdded: 1,
		},
		"total overlap": {
			base:      set.NewOrdered(1, 2, 3),
			give:      set.NewOrdered(1, 2, 3),
			want:      []int{1, 2, 3},
			wantAdded: 0,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			haveAdded := tt.base.AddOrderedSet(tt.give)
			require.Equal(t, tt.wantAdded, haveAdded)
			require.Equal(t, tt.want, tt.base.ToSlice())
		})
	}
}

func TestOrderedSet_ForEach(t *testing.T) {
	var (
		x     set.OrderedSet[int]
		calls int
	)

	x.ForEach(func(int) bool {
		calls++
		return true
	})
	require.Equal(t, 0, calls)

	seen := make([]int, 0, 3)
	x = set.NewOrdered(1, 2, 3)
	x.ForEach(func(value int) bool {
		seen = append(seen, value)
		calls++
		return true
	})
	require.Equal(t, 3, calls)
	require.Equal(t, []int{1, 2, 3}, seen)

	seen = seen[:0]
	calls = 0
	x.ForEach(func(value int) bool {
		seen = append(seen, value)
		calls++
		return false
	})
	require.Equal(t, 1, calls)
	require.Equal(t, []int{1}, seen)
}

func TestOrderedSet_Contains(t *testing.T) {
	x := set.NewOrdered(1, 2, 3, 4, 5)
	require.False(t, x.Contains(0))
	for i := 1; i <= 5; i++ {
		require.True(t, x.Contains(i))
	}
}

func TestOrderedSet_ContainsAny(t *testing.T) {
	x := set.NewOrdered(1, 2, 3, 4, 5)
	require.False(t, x.ContainsAny())
	require.False(t, x.ContainsAny(0))
	require.False(t, x.ContainsAny(0, -1, -2))
	require.True(t, x.ContainsAny(1))
	require.True(t, x.ContainsAny(1, 0, -1))
}

func TestOrderedSet_ContainsAll(t *testing.T) {
	x := set.NewOrdered(1, 2, 3, 4, 5)
	require.False(t, x.ContainsAll())
	require.False(t, x.ContainsAll(0))
	require.False(t, x.ContainsAll(0, -1, -2))
	require.False(t, x.ContainsAll(0, 1, 2))
	require.True(t, x.ContainsAll(1))
	require.True(t, x.ContainsAll(3, 2, 1))
}

func TestOrderedSet_Clear(t *testing.T) {
	var x set.OrderedSet[int]
	x.Clear() // for sanity

	x = set.NewOrdered(1, 2, 3)
	require.Equal(t, 3, x.Len())

	x.Clear()
	require.Equal(t, 0, x.Len())
}

func TestOrderedSet_Len(t *testing.T) {
	var x set.OrderedSet[int]
	require.Equal(t, 0, x.Len())

	x = set.NewOrdered[int]()
	require.Equal(t, 0, x.Len())

	require.True(t, x.Add(1))
	require.Equal(t, 1, x.Len())
	require.Equal(t, 1, set.New(1).Len())
}

func TestOrderedSet_ToSlice(t *testing.T) {
	cases := map[string]struct {
		give set.OrderedSet[int]
		want []int
	}{
		"empty": {
			give: set.OrderedSet[int]{},
			want: nil,
		},
		"single": {
			give: set.NewOrdered(1),
			want: []int{1},
		},
		"multiple": {
			give: set.NewOrdered(1, 2, 3),
			want: []int{1, 2, 3},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.give.ToSlice())
		})
	}
}

func TestOrderedSet_Intersect(t *testing.T) {
	cases := map[string]struct {
		base set.OrderedSet[int]
		give set.OrderedSet[int]
		want []int
	}{
		"empty base": {
			base: set.OrderedSet[int]{},
			give: set.NewOrdered(1, 2, 3),
			want: nil,
		},
		"empty upper": {
			base: set.NewOrdered(1, 2, 3),
			give: set.OrderedSet[int]{},
			want: nil,
		},
		"no overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.NewOrdered(4, 5, 6),
			want: nil,
		},
		"some overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.NewOrdered(2, 3, 4),
			want: []int{2, 3},
		},
		"total overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.NewOrdered(1, 2, 3),
			want: []int{1, 2, 3},
		},
		"base subset": {
			base: set.NewOrdered(1, 2, 3),
			give: set.NewOrdered(1, 2, 3, 4),
			want: []int{1, 2, 3},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.base.Intersect(tt.give).ToSlice())
		})
	}
}

func TestOrderedSet_UnorderedIntersect(t *testing.T) {
	cases := map[string]struct {
		base set.OrderedSet[int]
		give set.Set[int]
		want []int
	}{
		"empty base": {
			base: set.OrderedSet[int]{},
			give: set.New(1, 2, 3),
			want: nil,
		},
		"empty upper": {
			base: set.NewOrdered(1, 2, 3),
			give: set.Set[int]{},
			want: nil,
		},
		"no overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.New(4, 5, 6),
			want: nil,
		},
		"some overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.New(2, 3, 4),
			want: []int{2, 3},
		},
		"total overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.New(1, 2, 3),
			want: []int{1, 2, 3},
		},
		"base subset": {
			base: set.NewOrdered(1, 2, 3),
			give: set.New(1, 2, 3, 4),
			want: []int{1, 2, 3},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.ElementsMatch(
				t,
				tt.want,
				tt.base.UnorderedIntersect(tt.give).ToSlice(),
			)
		})
	}
}

func TestOrderedSet_Union(t *testing.T) {
	cases := map[string]struct {
		base set.OrderedSet[int]
		give set.OrderedSet[int]
		want []int
	}{
		"empty base": {
			base: set.OrderedSet[int]{},
			give: set.NewOrdered(1, 2, 3),
			want: []int{1, 2, 3},
		},
		"empty upper": {
			base: set.NewOrdered(1, 2, 3),
			give: set.OrderedSet[int]{},
			want: []int{1, 2, 3},
		},
		"no overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.NewOrdered(4, 5, 6),
			want: []int{1, 2, 3, 4, 5, 6},
		},
		"some overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.NewOrdered(2, 3, 4),
			want: []int{1, 4},
		},
		"total overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.NewOrdered(1, 2, 3),
			want: nil,
		},
		"base subset": {
			base: set.NewOrdered(1, 2, 3),
			give: set.NewOrdered(1, 2, 3, 4),
			want: []int{4},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.base.Union(tt.give).ToSlice())
		})
	}
}

func TestOrderedSet_UnorderedUnion(t *testing.T) {
	cases := map[string]struct {
		base set.OrderedSet[int]
		give set.Set[int]
		want []int
	}{
		"empty base": {
			base: set.OrderedSet[int]{},
			give: set.New(1, 2, 3),
			want: []int{1, 2, 3},
		},
		"empty upper": {
			base: set.NewOrdered(1, 2, 3),
			give: set.Set[int]{},
			want: []int{1, 2, 3},
		},
		"no overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.New(4, 5, 6),
			want: []int{1, 2, 3, 4, 5, 6},
		},
		"some overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.New(2, 3, 4),
			want: []int{1, 4},
		},
		"total overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.New(1, 2, 3),
			want: nil,
		},
		"base subset": {
			base: set.NewOrdered(1, 2, 3),
			give: set.New(1, 2, 3, 4),
			want: []int{4},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.ElementsMatch(
				t,
				tt.want,
				tt.base.UnorderedUnion(tt.give).ToSlice(),
			)
		})
	}
}

func TestOrderedSet_Merge(t *testing.T) {
	cases := map[string]struct {
		base set.OrderedSet[int]
		give set.OrderedSet[int]
		want []int
	}{
		"empty base": {
			base: set.OrderedSet[int]{},
			give: set.NewOrdered(1, 2, 3),
			want: []int{1, 2, 3},
		},
		"empty upper": {
			base: set.NewOrdered(1, 2, 3),
			give: set.OrderedSet[int]{},
			want: []int{1, 2, 3},
		},
		"no overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.NewOrdered(4, 5, 6),
			want: []int{1, 2, 3, 4, 5, 6},
		},
		"some overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.NewOrdered(2, 3, 4),
			want: []int{1, 2, 3, 4},
		},
		"total overlap": {
			base: set.NewOrdered(1, 2, 3),
			give: set.NewOrdered(1, 2, 3),
			want: []int{1, 2, 3},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.ElementsMatch(
				t,
				tt.want,
				tt.base.Merge(tt.give).ToSlice(),
			)
		})
	}
}

func TestUnorderedSet_ToSet(t *testing.T) {
	cases := map[string]struct {
		give set.OrderedSet[int]
		want []int
	}{
		"empty": {
			give: set.OrderedSet[int]{},
			want: nil,
		},
		"single": {
			give: set.NewOrdered(1),
			want: []int{1},
		},
		"multiple": {
			give: set.NewOrdered(1, 2, 3),
			want: []int{1, 2, 3},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				have = tt.give.ToSet()
				seen []int
			)

			have.ForEach(func(value int) bool {
				seen = append(seen, value)
				return true
			})

			require.ElementsMatch(t, tt.want, seen)
		})
	}
}
