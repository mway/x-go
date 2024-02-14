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

// Package flow provides workflow-related types and utilities.
package flow

import (
	"context"
	"errors"
	"slices"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrSkipRemainder is a sentinel value used to skip remaining actions.
	ErrSkipRemainder = errors.New("skip")

	_ Action = ActionFunc(nil)
	_ Action = Func(nil)
)

// ActionState is the cumulative state for an action.
type ActionState struct{}

// NewActionState creates a new [ActionState].
func NewActionState() ActionState {
	return ActionState{}
}

// An Action is something that runs.
type Action interface {
	// Run the action. The call must respect the lifetime of the given context.
	Run(context.Context, ActionState) error
}

// An ActionFunc is a function that is also an [Action].
type ActionFunc func(context.Context, ActionState) error

// Run the function with the given context and state.
func (f ActionFunc) Run(ctx context.Context, state ActionState) error {
	return f(ctx, state)
}

// A Func is a function that is also an [Action].
type Func func(context.Context) error

// Run the function with the given context.
func (f Func) Run(ctx context.Context, _ ActionState) error {
	return f(ctx)
}

// Linear creates a new linear [Action] comprised of the given sub-actions.
func Linear(actions ...Action) Action {
	return newGroup(actions...)
}

// Async creates a new async [Action] comprised of the given sub-actions.
func Async(actions ...Action) Action {
	return ThrottledAsync(-1, actions...)
}

// ThrottledAsync creates a new async [Action] comprised of the given
// sub-actions. The action will not use more than the given number of workers.
// Worker values < 0 indicate no limit.
func ThrottledAsync(workers int, actions ...Action) Action {
	return newAsyncGroup(workers, actions...)
}

type group []Action

func newGroup(actions ...Action) group {
	return slices.Clone(actions)
}

func (g group) Run(ctx context.Context, state ActionState) error {
	for _, action := range g {
		err := action.Run(ctx, state)
		switch {
		case errors.Is(err, ErrSkipRemainder):
			return nil
		case err == nil:
			// continue
		default:
			return err
		}
	}
	return nil
}

type asyncGroup struct {
	actions []Action
	workers int
}

func newAsyncGroup(workers int, actions ...Action) asyncGroup {
	return asyncGroup{
		actions: slices.Clone(actions),
		workers: workers,
	}
}

func (g asyncGroup) Run(ctx context.Context, state ActionState) error {
	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(g.workers)

	for _, action := range g.actions {
		action := action
		eg.Go(func() error {
			return action.Run(egCtx, state)
		})
	}

	err := eg.Wait()
	if !errors.Is(err, ErrSkipRemainder) {
		return err
	}

	return nil
}
