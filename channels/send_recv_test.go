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

package channels

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mway.dev/chrono/clock"

	"go.mway.dev/x/stub"
)

func TestRecv(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		var (
			values = []int{123, 234, 345}
			src    = newChannel(values...)
		)

		for _, want := range values {
			have, ok := Recv(context.Background(), src)
			require.True(t, ok)
			require.Equal(t, want, have)
		}
	})

	t.Run("closed channel", func(t *testing.T) {
		src := make(chan int)
		close(src)

		have, ok := Recv(context.Background(), src)
		require.False(t, ok)
		require.Zero(t, have)
	})

	t.Run("blocked channel", func(t *testing.T) {
		var (
			ready = make(chan struct{})
			done  = make(chan struct{})
			src   = make(chan int)
			want  = 12345
		)
		go func() {
			defer close(done)
			close(ready)
			have, ok := Recv(context.Background(), src)
			require.True(t, ok)
			require.Equal(t, have, want)
		}()

		timer := time.NewTimer(time.Second)
		defer timer.Stop()

		<-ready

		select {
		case <-timer.C:
			require.FailNow(t, "timed out waiting for send")
		case src <- want:
		}

		select {
		case <-timer.C:
			require.FailNow(t, "timed out waiting for return")
		case <-done:
		}
	})

	t.Run("blocked channel context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		var (
			ready = make(chan struct{})
			done  = make(chan struct{})
			src   chan int
		)
		go func() {
			defer close(done)
			close(ready)
			have, ok := Recv(ctx, src)
			require.False(t, ok)
			require.Zero(t, have)
		}()

		timer := time.NewTimer(time.Second)
		defer timer.Stop()

		<-ready
		time.AfterFunc(100*time.Millisecond, cancel)

		select {
		case <-timer.C:
			require.FailNow(t, "timed out waiting for return")
		case <-done:
		}
	})
}

func TestRecvWithTimeout(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		clk := clock.NewFakeClock()
		stub.With(&_newTimer, clk.NewTimer, func() {
			var (
				values = []int{123, 234, 345}
				src    = newChannel(values...)
			)

			for _, want := range values {
				have, ok := RecvWithTimeout(
					context.Background(),
					src,
					time.Second,
				)
				require.True(t, ok)
				require.Equal(t, want, have)
			}
		})
	})

	t.Run("closed channel", func(t *testing.T) {
		src := make(chan int)
		close(src)

		have, ok := RecvWithTimeout(context.Background(), src, time.Second)
		require.False(t, ok)
		require.Zero(t, have)
	})

	t.Run("blocked channel", func(t *testing.T) {
		var (
			ready = make(chan struct{})
			done  = make(chan struct{})
			src   = make(chan int)
			want  = 12345
		)
		go func() {
			defer close(done)
			close(ready)
			have, ok := RecvWithTimeout(context.Background(), src, time.Second)
			require.True(t, ok)
			require.Equal(t, have, want)
		}()

		timer := time.NewTimer(time.Second)
		defer timer.Stop()

		<-ready

		select {
		case <-timer.C:
			require.FailNow(t, "timed out waiting for send")
		case src <- want:
		}

		select {
		case <-timer.C:
			require.FailNow(t, "timed out waiting for return")
		case <-done:
		}
	})

	t.Run("blocked channel context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		var (
			ready = make(chan struct{})
			done  = make(chan struct{})
			src   chan int
		)
		go func() {
			defer close(done)
			close(ready)
			have, ok := RecvWithTimeout(ctx, src, time.Second)
			require.False(t, ok)
			require.Zero(t, have)
		}()

		timer := time.NewTimer(time.Second)
		defer timer.Stop()

		<-ready
		time.AfterFunc(100*time.Millisecond, cancel)

		select {
		case <-timer.C:
			require.FailNow(t, "timed out waiting for return")
		case <-done:
		}
	})

	t.Run("blocked channel timeout", func(t *testing.T) {
		clk := clock.NewFakeClock()
		stub.With(&_newTimer, clk.NewTimer, func() {
			var (
				ready = make(chan struct{})
				done  = make(chan struct{})
				src   chan int
			)
			go func() {
				defer close(done)
				close(ready)
				have, ok := RecvWithTimeout(
					context.Background(),
					src,
					time.Second,
				)
				require.False(t, ok)
				require.Zero(t, have)
			}()

			timer := time.NewTimer(time.Second)
			defer timer.Stop()

			<-ready
			time.AfterFunc(100*time.Millisecond, func() {
				clk.Add(time.Second)
			})

			select {
			case <-timer.C:
				require.FailNow(t, "timed out waiting for return")
			case <-done:
			}
		})
	})
}

func TestSend(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		dst := make(chan int, 1)
		require.True(t, Send(context.Background(), dst, 12345))
		require.Len(t, dst, 1)
		require.Equal(t, 12345, <-dst)
	})

	t.Run("nil channel", func(t *testing.T) {
		var dst chan int
		require.False(t, Send(context.Background(), dst, 12345))
	})

	t.Run("blocked channel context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		dst := make(chan int)
		defer close(dst)

		var (
			ready = make(chan struct{})
			done  = make(chan struct{})
		)
		go func() {
			defer close(done)
			close(ready)
			require.False(t, Send(ctx, dst, 12345))
		}()

		timer := time.NewTimer(time.Second)
		defer timer.Stop()

		<-ready
		time.AfterFunc(100*time.Millisecond, cancel)

		select {
		case <-timer.C:
			require.FailNow(t, "timed out waiting for return")
		case <-done:
			require.Len(t, dst, 0)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		dst := make(chan int, 1)
		require.False(t, Send(ctx, dst, 12345))
	})
}

func TestSendWithTimeout(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		clk := clock.NewFakeClock()
		stub.With(&_newTimer, clk.NewTimer, func() {
			dst := make(chan int, 1)
			require.True(t, SendWithTimeout(
				context.Background(),
				dst,
				12345,
				time.Second,
			))
			require.Len(t, dst, 1)
			require.Equal(t, 12345, <-dst)
		})
	})

	t.Run("timeout", func(t *testing.T) {
		clk := clock.NewFakeClock()
		stub.With(&_newTimer, clk.NewTimer, func() {
			time.AfterFunc(100*time.Millisecond, func() {
				clk.Add(time.Second)
			})

			dst := make(chan int)
			require.False(t, SendWithTimeout(
				context.Background(),
				dst,
				12345,
				time.Second,
			))
			require.Len(t, dst, 0)
		})
	})

	t.Run("nil channel", func(t *testing.T) {
		var dst chan int
		require.False(
			t,
			SendWithTimeout(context.Background(), dst, 12345, time.Second),
		)
	})

	t.Run("blocked channel context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		dst := make(chan int)
		defer close(dst)

		var (
			ready = make(chan struct{})
			done  = make(chan struct{})
		)
		go func() {
			defer close(done)
			close(ready)
			require.False(t, SendWithTimeout(ctx, dst, 12345, time.Second))
		}()

		timer := time.NewTimer(time.Second)
		defer timer.Stop()

		<-ready
		time.AfterFunc(100*time.Millisecond, cancel)

		select {
		case <-timer.C:
			require.FailNow(t, "timed out waiting for return")
		case <-done:
			require.Len(t, dst, 0)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		dst := make(chan int, 1)
		require.False(t, Send(ctx, dst, 12345))
	})
}

func newChannel[T any](values ...T) chan T {
	ch := make(chan T, len(values))
	for _, value := range values {
		ch <- value
	}
	return ch
}
