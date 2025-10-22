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

package future_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.mway.dev/x/future"
)

func TestNew(t *testing.T) {
	testFuture(t, future.New[int])
}

func TestWaitContext(t *testing.T) {
	var (
		f = future.New[int]()
		p = f.Promise()
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.False(t, f.IsSet())
	require.False(t, f.IsCanceled())
	have, ok := f.WaitContext(ctx)
	require.False(t, ok)
	require.Zero(t, have)
	have, ok = p.WaitContext(ctx)
	require.False(t, ok)
	require.Zero(t, have)

	f.Set(123)

	require.True(t, f.IsSet())
	require.False(t, f.IsCanceled())

	have, ok = f.Wait()
	require.True(t, ok)
	require.Equal(t, 123, have)
	have, ok = p.Wait()
	require.True(t, ok)
	require.Equal(t, 123, have)
	have, ok = f.WaitContext(ctx)
	require.True(t, ok)
	require.Equal(t, 123, have)
	have, ok = p.WaitContext(ctx)
	require.True(t, ok)
	require.Equal(t, 123, have)
}

func TestWaitRoutine(t *testing.T) {
	const want = 123
	var (
		f     = future.New[int]()
		p     = f.Promise()
		ready = make(chan struct{}, 2)
		wg    sync.WaitGroup
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ready <- struct{}{}
		have, ok := f.WaitContext(ctx)
		require.True(t, f.IsSet())
		require.False(t, f.IsCanceled())
		require.True(t, ok)
		require.Equal(t, want, have)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ready <- struct{}{}
		have, ok := p.WaitContext(ctx)
		require.True(t, p.IsSet())
		require.False(t, p.IsCanceled())
		require.True(t, ok)
		require.Equal(t, want, have)
	}()

	for range cap(ready) {
		<-ready
	}

	time.AfterFunc(100*time.Millisecond, func() {
		f.Set(want)
	})

	// n.b. This wait will nominally be zero, but is otherwise implicitly bound
	//      by the lifetime of ctx above.
	wg.Wait()
}

func BenchmarkNew(b *testing.B) {
	var x *future.Future[int]
	for range b.N {
		x = future.New[int]()
	}
	_ = x
}

func testFuture(t *testing.T, fn func() *future.Future[int]) {
	t.Helper()

	t.Run("smoke", func(t *testing.T) {
		var (
			f = fn()
			p = f.Promise()
		)

		have, ok := f.Get()
		require.False(t, ok)
		require.Zero(t, have)
		have, ok = p.Get()
		require.False(t, ok)
		require.Zero(t, have)

		f.Set(123)
		f.Set(456)
		for range 3 {
			have, ok = f.Get()
			require.True(t, ok)
			require.Equal(t, 123, have)
			have, ok = p.Get()
			require.True(t, ok)
			require.Equal(t, 123, have)
			have, ok = f.Wait()
			require.True(t, ok)
			require.Equal(t, 123, have)
			have, ok = p.Wait()
			require.True(t, ok)
			require.Equal(t, 123, have)
		}
	})

	t.Run("cancel", func(t *testing.T) {
		var (
			f = fn()
			p = f.Promise()
		)

		require.False(t, f.IsSet())
		require.False(t, f.IsCanceled())
		f.Cancel()
		require.False(t, f.IsSet())
		require.True(t, f.IsCanceled())
		f.Cancel() // double cancel for sanity
		f.Set(123)
		require.False(t, f.IsSet())

		have, ok := f.Get()
		require.False(t, ok)
		require.Zero(t, have)
		have, ok = p.Get()
		require.False(t, ok)
		require.Zero(t, have)
		have, ok = f.Wait()
		require.False(t, ok)
		require.Zero(t, have)
		have, ok = p.Wait()
		require.False(t, ok)
		require.Zero(t, have)
	})
}
