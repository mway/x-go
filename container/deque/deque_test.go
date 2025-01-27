package deque_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/container/deque"
)

func TestDeque(t *testing.T) {
	require.Nil(t, deque.NewLinkedWithValues[int]())

	var (
		tmps = deque.NewWithValues(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		ds   = deque.New[int](10)
		tmpl = deque.NewLinkedWithValues(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		dl   = deque.NewLinked[int]()
	)

	x, ok := ds.MaybeFront()
	require.False(t, ok)
	require.Zero(t, x)
	require.Zero(t, ds.Front())
	x, ok = dl.MaybeFront()
	require.False(t, ok)
	require.Zero(t, x)
	require.Zero(t, dl.Front())
	x, ok = ds.MaybePopFront()
	require.False(t, ok)
	require.Zero(t, x)
	x, ok = dl.MaybePopFront()
	require.False(t, ok)
	require.Zero(t, x)
	x, ok = ds.MaybeBack()
	require.False(t, ok)
	require.Zero(t, x)
	require.Zero(t, ds.Back())
	x, ok = dl.MaybeBack()
	require.False(t, ok)
	require.Zero(t, x)
	require.Zero(t, dl.Back())
	x, ok = ds.MaybePopBack()
	require.False(t, ok)
	require.Zero(t, x)
	x, ok = dl.MaybePopBack()
	require.False(t, ok)
	require.Zero(t, x)

	var (
		calls = make(map[int]int)
		count = func(x int) bool {
			calls[x]++
			return true
		}
	)
	ds.PeekEachBack(count)
	ds.PeekEachFront(count)
	ds.PopEachBack(count)
	ds.PopEachFront(count)
	dl.PeekEachBack(count)
	dl.PeekEachFront(count)
	dl.PopEachBack(count)
	dl.PopEachFront(count)
	require.Len(t, calls, 0)

	for i := range 10 {
		ds.PushBack(i + 1)
		require.Equal(t, 1, ds.Front())
		require.Equal(t, i+1, ds.Back())
		dl.PushBack(i + 1)
		require.Equal(t, 1, dl.Front())
		require.Equal(t, i+1, dl.Back())
	}
	require.Equal(t, tmps, ds)

	var (
		tmpls []int
		dls   []int
	)
	tmpl.PeekEachFront(func(i int) bool {
		tmpls = append(tmpls, i)
		return true
	})
	dl.PeekEachFront(func(i int) bool {
		dls = append(dls, i)
		return true
	})
	require.Equal(t, tmpls, dls)

	ds.PushFront(0)
	require.Equal(t, 11, ds.Len())
	dl.PushFront(0)
	require.Equal(t, 11, dl.Len())

	clear(calls)
	var (
		limit            = 14
		countWithinLimit = func(x int) bool {
			calls[x]++
			limit--
			return limit >= 0
		}
	)
	for limit > 0 {
		ds.PeekEachBack(countWithinLimit)
		ds.PeekEachFront(countWithinLimit)
		dl.PeekEachBack(countWithinLimit)
		dl.PeekEachFront(countWithinLimit)
	}
	ds.PeekEachBack(countWithinLimit)
	ds.PeekEachFront(countWithinLimit)
	dl.PeekEachBack(countWithinLimit)
	dl.PeekEachFront(countWithinLimit)
	require.Equal(
		t,
		map[int]int{
			0:  5,
			1:  2,
			2:  2,
			3:  2,
			4:  1,
			5:  1,
			6:  1,
			7:  1,
			8:  1,
			9:  1,
			10: 4,
		},
		calls,
	)

	clear(calls)
	ds.PeekEachBack(count)
	ds.PeekEachFront(count)
	ds.PopEachBack(count)
	dl.PeekEachBack(count)
	dl.PeekEachFront(count)
	dl.PopEachBack(count)
	for i := range 11 {
		ds.PushBack(i)
		dl.PushBack(i)
	}
	ds.PopEachFront(count)
	dl.PopEachFront(count)
	require.Equal(
		t,
		map[int]int{
			0:  8,
			1:  8,
			2:  8,
			3:  8,
			4:  8,
			5:  8,
			6:  8,
			7:  8,
			8:  8,
			9:  8,
			10: 8,
		},
		calls,
	)

	for i := range 11 {
		ds.PushBack(i)
		dl.PushBack(i)
	}
	for i := range 11 {
		require.Equal(t, i, ds.PopFront())
		require.Equal(t, i, dl.PopFront())
	}
	for i := range 11 {
		ds.PushFront(i)
		dl.PushFront(i)
	}
	for i := range 11 {
		require.Equal(t, i, ds.PopBack())
		require.Equal(t, i, dl.PopBack())
	}
}

func BenchmarkDeque(b *testing.B) {
	sizes := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}
	for _, size := range sizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			b.Run("slice-backed", func(b *testing.B) {
				b.Run("front to back", func(b *testing.B) {
					d := deque.New[int](size)
					for i := range size {
						d.PushBack(i)
					}

					b.ResetTimer()
					for range b.N {
						d.PushBack(d.PopFront())
					}
				})

				b.Run("back to front", func(b *testing.B) {
					d := deque.New[int](size)
					for i := range size {
						d.PushBack(i)
					}

					b.ResetTimer()
					for range b.N {
						d.PushFront(d.PopBack())
					}
				})
			})

			b.Run("doubly-linked-list", func(b *testing.B) {
				b.Run("front to back", func(b *testing.B) {
					d := deque.NewLinked[int]()
					for i := range size {
						d.PushBack(i)
					}

					b.ResetTimer()
					for range b.N {
						d.PushBack(d.PopFront())
					}
				})

				b.Run("back to front", func(b *testing.B) {
					d := deque.NewLinked[int]()
					for i := range size {
						d.PushBack(i)
					}

					b.ResetTimer()
					for range b.N {
						d.PushFront(d.PopBack())
					}
				})
			})
		})
	}
}
