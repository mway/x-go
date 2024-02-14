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

package flow_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/flow"
)

func TestActionFunc_Run(t *testing.T) {
	var (
		wantCtx, cancel = context.WithCancel(context.Background())
		wantState       = flow.ActionState{}
		calls           int
		action          = flow.ActionFunc(func(
			ctx context.Context,
			state flow.ActionState,
		) error {
			calls++
			require.Equal(t, wantCtx, ctx)
			require.Equal(t, wantState, state)
			return nil
		})
	)

	defer cancel()
	require.NoError(t, action.Run(wantCtx, wantState))
	require.Equal(t, 1, calls)

	wantErr := errors.New("want error")
	action = flow.ActionFunc(func(context.Context, flow.ActionState) error {
		return wantErr
	})
	require.ErrorIs(t, action.Run(wantCtx, wantState), wantErr)
}

func TestLinear(t *testing.T) {
	var (
		calls     = make([]int, 0, 3)
		wantCalls = []int{0, 1, 2}
		action    = flow.Linear(
			flow.Func(func(context.Context) error {
				calls = append(calls, wantCalls[0])
				return nil
			}),
			flow.Func(func(context.Context) error {
				calls = append(calls, wantCalls[1])
				return nil
			}),
			flow.Func(func(context.Context) error {
				calls = append(calls, wantCalls[2])
				return nil
			}),
		)
		ctx   = context.Background()
		state = flow.NewActionState()
	)

	for i := 0; i < 128; i++ {
		require.NoError(t, action.Run(ctx, state))
		require.Equal(t, wantCalls, calls)
		calls = calls[:0]
	}
}

func TestLinear_Error(t *testing.T) {
	var (
		wantError = errors.New(t.Name())
		action    = flow.Linear(
			flow.Func(func(context.Context) error {
				return nil
			}),
			flow.Func(func(context.Context) error {
				return wantError
			}),
			flow.Func(func(context.Context) error {
				return nil
			}),
		)
		ctx   = context.Background()
		state = flow.NewActionState()
	)

	for i := 0; i < 128; i++ {
		require.ErrorIs(t, action.Run(ctx, state), wantError)
	}
}

func TestLinear_SkipRemainder(t *testing.T) {
	var (
		calls     = make([]int, 0, 3)
		wantCalls = []int{0}
		action    = flow.Linear(
			flow.Func(func(context.Context) error {
				calls = append(calls, wantCalls[0])
				return flow.ErrSkipRemainder
			}),
			flow.Func(func(context.Context) error {
				calls = append(calls, wantCalls[1])
				return nil
			}),
			flow.Func(func(context.Context) error {
				calls = append(calls, wantCalls[2])
				return nil
			}),
		)
		ctx   = context.Background()
		state = flow.NewActionState()
	)

	for i := 0; i < 128; i++ {
		require.NoError(t, action.Run(ctx, state))
		require.Equal(t, wantCalls, calls)
		calls = calls[:0]
	}
}

func TestAsync(t *testing.T) {
	var (
		mu     sync.Mutex
		calls  = make([]string, 0, 3)
		seen   = make(map[string]int)
		action = flow.Async(
			flow.Func(func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				calls = append(calls, "0")
				return nil
			}),
			flow.Func(func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				calls = append(calls, "1")
				return nil
			}),
			flow.Func(func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				calls = append(calls, "2")
				return nil
			}),
		)
		ctx   = context.Background()
		state = flow.NewActionState()
	)

	for i := 0; i < 1024; i++ {
		require.NoError(t, action.Run(ctx, state))
		str := strings.Join(calls, "")
		seen[str]++
		calls = calls[:0]
	}

	require.Len(t, seen, 6)
	for order, times := range seen {
		require.GreaterOrEqual(t, times, 1, order)
	}
}

func TestAsync_Error(t *testing.T) {
	var (
		wantError = errors.New(t.Name())
		action    = flow.Async(
			flow.Func(func(context.Context) error {
				return nil
			}),
			flow.Func(func(context.Context) error {
				return wantError
			}),
			flow.Func(func(context.Context) error {
				return nil
			}),
		)
		ctx   = context.Background()
		state = flow.NewActionState()
	)

	for i := 0; i < 1024; i++ {
		require.ErrorIs(t, action.Run(ctx, state), wantError)
	}
}

func TestThrottledAsync(t *testing.T) {
	var (
		mu     sync.Mutex
		calls  = make([]string, 0, 3)
		seen   = make(map[string]int)
		action = flow.ThrottledAsync(
			2,
			flow.Func(func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				calls = append(calls, "0")
				return nil
			}),
			flow.Func(func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				calls = append(calls, "1")
				return nil
			}),
			flow.Func(func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				calls = append(calls, "2")
				return nil
			}),
		)
		ctx   = context.Background()
		state = flow.NewActionState()
	)

	for i := 0; i < 1024; i++ {
		require.NoError(t, action.Run(ctx, state))
		str := strings.Join(calls, "")
		seen[str]++
		calls = calls[:0]
	}

	require.Len(t, seen, 4)
	for order, times := range seen {
		require.GreaterOrEqual(t, times, 1, order)
	}
}

func TestThrottledAsync_Linear(t *testing.T) {
	var (
		mu        sync.Mutex
		wantCalls = []int{0, 1, 2}
		calls     = make([]int, 0, 3)
		action    = flow.ThrottledAsync(
			1, // max 1 worker at a time
			flow.Func(func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				calls = append(calls, 0)
				return nil
			}),
			flow.Func(func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				calls = append(calls, 1)
				return nil
			}),
			flow.Func(func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				calls = append(calls, 2)
				return nil
			}),
		)
		ctx   = context.Background()
		state = flow.NewActionState()
	)

	for i := 0; i < 1024; i++ {
		require.NoError(t, action.Run(ctx, state))
		require.Equal(t, wantCalls, calls)
		calls = calls[:0]
	}
}

func TestThrottledAsync_Linear_SkipRemainder(t *testing.T) {
	var (
		mu        sync.Mutex
		wantCalls = []int{0}
		calls     = make([]int, 0, 3)
		action    = flow.ThrottledAsync(
			1, // max 1 worker at a time
			flow.Func(func(context.Context) error {
				mu.Lock()
				defer mu.Unlock()
				calls = append(calls, 0)
				return flow.ErrSkipRemainder
			}),
			flow.Func(func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return nil
				default:
					mu.Lock()
					defer mu.Unlock()
					calls = append(calls, 1)
					return nil
				}
			}),
			flow.Func(func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return nil
				default:
					mu.Lock()
					defer mu.Unlock()
					calls = append(calls, 2)
					return nil
				}
			}),
		)
		ctx   = context.Background()
		state = flow.NewActionState()
	)

	for i := 0; i < 1024; i++ {
		require.NoError(t, action.Run(ctx, state))
		require.Equal(t, wantCalls, calls)
		calls = calls[:0]
	}
}
