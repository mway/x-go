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

package maps_test

import (
	gomaps "maps"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/maps"
)

func TestFilter(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.Nil(t, maps.Filter((map[int]int)(nil), func(int, int) bool {
			return true
		}))
	})

	t.Run("completely filtered out", func(t *testing.T) {
		var (
			src      = map[int]int{1: 1, 2: 2, 3: 3}
			haveNone = maps.Filter(src, maps.ByKey(func(k int) bool {
				return k < 0
			}))
		)
		require.Nil(t, haveNone)
	})
}

func TestFilter_DiscreteTypes(t *testing.T) {
	src := map[int]string{
		1: "one",
		2: "two",
		3: "three",
		4: "four",
	}

	haveNone := maps.Filter(src, maps.ByKey(func(k int) bool {
		return k < 0
	}))
	require.Nil(t, haveNone)

	haveEvens := maps.Filter(src, maps.ByKey(func(k int) bool {
		return k%2 == 0
	}))
	require.Equal(t, map[int]string{
		2: "two",
		4: "four",
	}, haveEvens)

	haveShortWords := maps.Filter(src, maps.ByValue(func(v string) bool {
		return len(v) <= 3
	}))
	require.Equal(t, map[int]string{
		1: "one",
		2: "two",
	}, haveShortWords)

	haveKeyMatchesWordLength := maps.Filter(src, func(k int, v string) bool {
		return k == len(v)
	})
	require.Equal(t, map[int]string{
		4: "four",
	}, haveKeyMatchesWordLength)
}

func TestFilter_EqualTypes(t *testing.T) {
	src := map[int]int{
		1: 2,
		2: 3,
		3: 4,
		4: 5,
	}

	haveEvenKeys := maps.Filter(src, maps.ByKey(func(k int) bool {
		return k%2 == 0
	}))
	require.Equal(t, map[int]int{
		2: 3,
		4: 5,
	}, haveEvenKeys)

	haveEvenValues := maps.Filter(src, maps.ByValue(func(v int) bool {
		return v%2 == 0
	}))
	require.Equal(t, map[int]int{
		1: 2,
		3: 4,
	}, haveEvenValues)

	haveProductFactor4 := maps.Filter(src, func(k int, v int) bool {
		return (k*v)%4 == 0
	})
	require.Equal(t, map[int]int{
		3: 4,
		4: 5,
	}, haveProductFactor4)
}

func TestTransform(t *testing.T) {
	t.Run("same types", func(t *testing.T) {
		var (
			give   = map[int]int{1: 1, 2: 2, 3: 3}
			want   = map[int]int{2: 2, 4: 4, 6: 6}
			mapper = func(k int, v int) (int, int) { return k * 2, v * 2 }
		)
		require.Equal(t, want, maps.Transform[map[int]int](give, mapper))
	})

	t.Run("different key type", func(t *testing.T) {
		var (
			give   = map[int]int{1: 1, 2: 2, 3: 3}
			want   = map[string]int{"2": 2, "4": 4, "6": 6}
			mapper = func(k int, v int) (string, int) {
				return strconv.Itoa(k * 2), v * 2
			}
		)
		require.Equal(t, want, maps.Transform[map[string]int](give, mapper))
	})

	t.Run("different value type", func(t *testing.T) {
		var (
			give   = map[int]int{1: 1, 2: 2, 3: 3}
			want   = map[int]string{2: "2", 4: "4", 6: "6"}
			mapper = func(k int, v int) (int, string) {
				return k * 2, strconv.Itoa(v * 2)
			}
		)
		require.Equal(t, want, maps.Transform[map[int]string](give, mapper))
	})

	t.Run("different types", func(t *testing.T) {
		var (
			give   = map[int]int{1: 1, 2: 2, 3: 3}
			want   = map[string]string{"2": "2", "4": "4", "6": "6"}
			mapper = func(k int, v int) (string, string) {
				return strconv.Itoa(k * 2), strconv.Itoa(v * 2)
			}
		)
		require.Equal(t, want, maps.Transform[map[string]string](give, mapper))
	})
}

func TestIterIter2(t *testing.T) {
	t.Run("sanity", func(t *testing.T) {
		var (
			want = map[int]int{
				1: 1,
				2: 2,
				3: 3,
			}
			iter  = maps.Iter(want)
			iter2 = maps.Iter2(want)
		)

		require.ElementsMatch(
			t,
			slices.Collect(gomaps.Values(want)),
			slices.Collect(iter),
		)
		require.Equal(t, want, gomaps.Collect(iter2))
	})

	t.Run("early return", func(t *testing.T) {
		var (
			want = map[int]int{
				1: 1,
				2: 2,
				3: 3,
			}
			iter       = maps.Iter(want)
			iter2      = maps.Iter2(want)
			haveValues []int
			haveMap    = make(map[int]int)
		)

		iter(func(v int) bool {
			haveValues = append(haveValues, v)
			return false
		})
		require.Len(t, haveValues, 1)
		require.Contains(t, slices.Collect(gomaps.Values(want)), haveValues[0])

		iter2(func(k int, v int) bool {
			haveMap[k] = v
			return false
		})
		require.Len(t, haveMap, 1)
		for k, v := range haveMap {
			require.Contains(t, want, k)
			require.Equal(t, want[k], v)
		}
	})
}
