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

package ptr_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/container/ptr"
)

func TestTo(t *testing.T) {
	have := ptr.To(t.Name())
	require.NotNil(t, have)
	require.Equal(t, t.Name(), *have)
}

func TestLoad(t *testing.T) {
	require.Zero(t, ptr.Load((*string)(nil)))
	require.Equal(t, t.Name(), ptr.Load(ptr.To(t.Name())))
}

func TestLoadOr(t *testing.T) {
	require.Equal(t, 123, ptr.LoadOr((*int)(nil), 123))
	require.Equal(t, 456, ptr.LoadOr(ptr.To(456), 123))
}

func TestLoadOrElse(t *testing.T) {
	require.Equal(
		t,
		123,
		ptr.LoadOrElse((*int)(nil), func() int { return 123 }),
	)
	require.Equal(
		t,
		456,
		ptr.LoadOrElse(ptr.To(456), func() int { return 123 }),
	)
}

func TestIf(t *testing.T) {
	var (
		x     = ptr.To(t.Name())
		calls int
	)
	require.True(t, ptr.If(x, func(strp *string) {
		require.NotNil(t, strp)
		require.Equal(t, t.Name(), *strp)
		calls++
	}))
	require.Equal(t, 1, calls)
	require.False(t, ptr.If(nil, func(*int) {
		require.FailNow(t, "unexpected If callback call")
	}))
}

func TestLoadIf(t *testing.T) {
	var (
		x     = ptr.To(t.Name())
		calls int
	)
	require.True(t, ptr.LoadIf(x, func(str string) {
		require.Equal(t, t.Name(), str)
		calls++
	}))
	require.Equal(t, 1, calls)
	require.False(t, ptr.LoadIf(nil, func(int) {
		require.FailNow(t, "unexpected LoadIf callback call")
	}))
}

func TestNew(t *testing.T) {
	var p ptr.Pointer[int]
	require.False(t, p.Held())

	p = ptr.New(123)
	require.True(t, p.Held())
	require.Equal(t, 123, p.Load())

	tmp := ptr.To(456)
	p = ptr.NewPtr(tmp)
	require.True(t, p.Held())
	require.Equal(t, 456, p.Load())
	require.True(t, tmp == p.Raw())

	tmp = ptr.To(789)
	p = ptr.NewPtrCopy(tmp)
	require.True(t, p.Held())
	require.Equal(t, 789, p.Load())
	require.False(t, tmp == p.Raw())
}

func TestPointer_Call(t *testing.T) {
	var p ptr.Pointer[string]
	requirePointerCallsNotHeld(t, p)

	p = ptr.New(t.Name())
	requirePointerCallsHeld(t, p, t.Name())
}

func TestPointer_Store(t *testing.T) {
	var p ptr.Pointer[int]
	require.False(t, p.Held())
	require.True(t, p.MaybeStore(123))
	require.True(t, p.Held())
	require.Equal(t, 123, p.Load())
	require.False(t, p.MaybeStore(456))
	require.Equal(t, 123, p.Load())
	p.Store(456)
	require.Equal(t, 456, p.Load())
}

func TestPointer_StorePtr(t *testing.T) {
	var p ptr.Pointer[int]
	require.False(t, p.Held())
	require.True(t, p.MaybeStorePtr(ptr.To(123)))
	require.True(t, p.Held())
	require.Equal(t, 123, p.Load())
	require.False(t, p.MaybeStorePtr(ptr.To(456)))
	require.Equal(t, 123, p.Load())
	p.StorePtr(ptr.To(456))
	require.Equal(t, 456, p.Load())
}

func TestPointer_StorePtrCopy(t *testing.T) {
	var p ptr.Pointer[int]
	require.False(t, p.Held())
	require.True(t, p.MaybeStorePtrCopy(ptr.To(123)))
	require.True(t, p.Held())
	require.Equal(t, 123, p.Load())
	require.False(t, p.MaybeStorePtrCopy(ptr.To(456)))
	require.Equal(t, 123, p.Load())
	p.StorePtrCopy(ptr.To(456))
	require.Equal(t, 456, p.Load())
}

func TestPointer_Load(t *testing.T) {
	var p ptr.Pointer[int]
	require.Equal(t, 0, p.Load())
	require.Equal(t, 123, p.LoadOr(123))
	require.Equal(t, 123, p.LoadOrElse(func() int { return 123 }))

	p = ptr.New(456)
	require.Equal(t, 456, p.Load())
	require.Equal(t, 456, p.LoadOr(123))
	require.Equal(t, 456, p.LoadOrElse(func() int { return 123 }))
}

func TestPointer_Move(t *testing.T) {
	var (
		src  = ptr.New(123)
		dst  ptr.Pointer[int]
		zero ptr.Pointer[int]
	)

	require.False(t, zero.Move(&dst))
	require.False(t, dst.Held())
	require.True(t, src.Held())
	require.True(t, src.Move(&dst))
	require.True(t, dst.Held())
	require.False(t, src.Held())
	require.Equal(t, 123, dst.Load())

	var pp *ptr.Pointer[int]
	require.False(t, pp.Move(&dst))
	require.True(t, dst.Held())
	require.Equal(t, 123, dst.Load())
}

func BenchmarkPointer_Load(b *testing.B) {
	b.Run("ptr", func(b *testing.B) {
		p := ptr.New(123)
		b.ReportAllocs()
		b.ResetTimer()

		var x int
		for i := 0; i < b.N; i++ {
			x = p.Load()
		}
		_ = x
	})

	b.Run("raw", func(b *testing.B) {
		p := ptr.To(123)
		b.ReportAllocs()
		b.ResetTimer()

		var x int
		for i := 0; i < b.N; i++ {
			if p != nil {
				x = *p
			}
		}
		_ = x
	})
}

func BenchmarkPointer_LoadOr(b *testing.B) {
	b.Run("ptr", func(b *testing.B) {
		p := ptr.New(123)
		b.ReportAllocs()
		b.ResetTimer()

		var x int
		for i := 0; i < b.N; i++ {
			x = p.LoadOr(b.N)
		}
		_ = x
	})

	b.Run("raw", func(b *testing.B) {
		p := ptr.To(123)
		b.ReportAllocs()
		b.ResetTimer()

		var x int
		for i := 0; i < b.N; i++ {
			if p != nil {
				x = *p
			} else {
				x = b.N
			}
		}
		_ = x
	})
}

func BenchmarkPointer_LoadOrElse(b *testing.B) {
	b.Run("ptr", func(b *testing.B) {
		p := ptr.New(123)
		b.ReportAllocs()
		b.ResetTimer()

		var x int
		for i := 0; i < b.N; i++ {
			x = p.LoadOrElse(func() int { return b.N })
		}
		_ = x
	})

	b.Run("raw", func(b *testing.B) {
		p := ptr.To(123)
		b.ReportAllocs()
		b.ResetTimer()

		var x int
		for i := 0; i < b.N; i++ {
			if p != nil {
				x = *p
			} else {
				x = func() int { return b.N }()
			}
		}
		_ = x
	})
}

func requirePointerCallsNotHeld[T any](t *testing.T, p ptr.Pointer[T]) {
	t.Helper()

	require.False(t, p.MaybeCall(func(T) {
		require.FailNow(t, fmt.Sprintf("unexpected call: %T.MaybeCall", p))
	}))

	require.False(t, p.MaybeCallPtr(func(*T) {
		require.FailNow(t, fmt.Sprintf("unexpected call: %T.MaybeCallPtr", p))
	}))

	p.Call(func(have T, held bool) {
		require.Zero(t, have)
		require.False(t, held)
	})

	p.CallPtr(func(have *T, held bool) {
		require.Nil(t, have)
		require.False(t, held)
	})
}

func requirePointerCallsHeld[T any](t *testing.T, p ptr.Pointer[T], want T) {
	t.Helper()

	require.True(t, p.MaybeCall(func(have T) {
		require.Equal(t, want, have)
	}))

	require.True(t, p.MaybeCallPtr(func(have *T) {
		require.NotNil(t, have)
		require.Equal(t, want, *have)
	}))

	p.Call(func(have T, held bool) {
		require.Equal(t, want, have)
		require.True(t, held)
	})

	p.CallPtr(func(have *T, held bool) {
		require.NotNil(t, have)
		require.Equal(t, want, *have)
		require.True(t, held)
	})
}
