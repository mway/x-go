// Copyright (c) 2022 Matt Way
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

package context_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mway.dev/chrono/clock"
	"go.mway.dev/x/context"
)

func TestNewDebouncedFactory(t *testing.T) {
	cases := map[string]struct {
		timeout     time.Duration
		opts        []context.DebouncedOption
		expectError error
	}{
		"timeout <0s defaults": {
			timeout:     -1,
			opts:        nil,
			expectError: context.ErrInvalidTimeout,
		},
		"timeout 0s debounce <0s": {
			timeout: 0,
			opts: []context.DebouncedOption{
				context.WithDebounce(-1),
			},
			expectError: context.ErrInvalidDebounce,
		},
		"timeout 0s debounce <0s with options": {
			timeout: 0,
			opts: []context.DebouncedOption{
				context.DebouncedOptions{
					Debounce: -1,
				},
			},
			expectError: context.ErrInvalidDebounce,
		},
		"timeout 0s defaults": {
			timeout:     0,
			opts:        nil,
			expectError: context.ErrInvalidDebounce,
		},
		"timeout 1s defaults": {
			timeout:     time.Second,
			opts:        nil,
			expectError: nil,
		},
		"timeout 1s debounce 2s": {
			timeout: time.Second,
			opts: []context.DebouncedOption{
				context.WithDebounce(2 * time.Second),
			},
			expectError: context.ErrInvalidDebounce,
		},
		"nil clock": {
			timeout: 0,
			opts: []context.DebouncedOption{
				context.WithClock(nil),
			},
			expectError: context.ErrNilClock,
		},
		"nil parent context": {
			timeout: 0,
			opts: []context.DebouncedOption{
				context.WithContext(nil), //nolint:govet,staticcheck
			},
			expectError: context.ErrNilContext,
		},
		"nil context func": {
			timeout: 0,
			opts: []context.DebouncedOption{
				context.WithContextFunc(nil),
			},
			expectError: context.ErrNilContextFunc,
		},
		"full options": {
			timeout: time.Second,
			opts: []context.DebouncedOption{
				context.DebouncedOptions{
					Debounce:    time.Second,
					Clock:       clock.NewFakeClock(),
					Context:     context.TODO(),
					ContextFunc: context.WithTimeout,
				},
			},
			expectError: nil,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			factory, err := context.NewDebouncedFactory(tt.timeout, tt.opts...)
			require.ErrorIs(t, err, tt.expectError)

			if tt.expectError != nil {
				require.Nil(t, factory)
			} else {
				require.NotNil(t, factory)
			}
		})
	}
}

func TestDebouncedFactory_Get_Deadline(t *testing.T) {
	var (
		factory, clk = newFactory(
			t,
			time.Second,
			context.WithDebounce(time.Second),
		)
		ctxA           = factory.Get()
		deadlineA, okA = ctxA.Deadline()
	)

	require.True(t, okA)

	for i := 0; i < 100; i++ {
		ctx := factory.Get()
		require.Equal(t, ctxA, ctx)

		dl, ok := ctx.Deadline()
		require.True(t, ok)
		require.True(t, dl.Equal(deadlineA))
	}

	clk.Add(10 * time.Second)

	var (
		ctxB           = factory.Get()
		deadlineB, okB = ctxB.Deadline()
	)

	require.NotEqual(t, ctxA, ctxB)
	require.True(t, okB)
	require.False(t, deadlineB.Equal(deadlineA))
}

func TestDebouncedFactory_ContextReuse(t *testing.T) {
	var (
		factory, _ = newFactory(
			t,
			time.Second,
			context.WithDebounce(time.Second),
		)
		ctx = factory.Get()
	)

	for i := 0; i < 1000; i++ {
		require.Equal(t, ctx, factory.Get())
	}
}

func TestDebouncedFactory_ParentContext(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		factory, _ = newFactory(t, time.Hour, context.WithContext(parent))
		ctx        = factory.Get()
	)

	select {
	case <-ctx.Done():
		require.FailNow(t, "context unexpectedly done")
	default:
	}

	factory.Cancel()

	select {
	case <-ctx.Done():
	default:
		require.FailNow(t, "context unexpectedly not done")
	}
}

func TestDebouncedFactory_Cancel(t *testing.T) {
	var (
		factory, _ = newFactory(t, time.Hour)
		ctx        = factory.Get()
	)

	select {
	case <-ctx.Done():
		require.FailNow(t, "context unexpectedly done")
	default:
	}

	factory.Cancel()

	select {
	case <-ctx.Done():
	default:
		require.FailNow(t, "context unexpectedly not done")
	}
}

func newFactory(
	tb testing.TB,
	timeout time.Duration,
	opts ...context.DebouncedOption,
) (*context.DebouncedFactory, *clock.FakeClock) {
	clk := clock.NewFakeClock()
	opts = append(
		append([]context.DebouncedOption(nil), opts...),
		context.WithClock(clk),
	)

	f, err := context.NewDebouncedFactory(timeout, opts...)
	require.NoError(tb, err)
	return f, clk
}
