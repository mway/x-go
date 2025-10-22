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

// Package channels provides helpers for working with sets.
package channels

import (
	"context"
	"time"

	"go.mway.dev/chrono/clock"
)

var (
	_clock    = clock.NewMonotonicClock()
	_newTimer = func(d time.Duration) *clock.Timer {
		return _clock.NewTimer(d)
	}
)

// Send will send value to ch, blocking until either the send succeeds or the
// given context is canceled. Returns whether the send was successful.
func Send[T any](ctx context.Context, ch chan<- T, value T) bool {
	if ch == nil || ctx.Err() != nil {
		return false
	}

	select {
	case <-ctx.Done():
		return false
	case ch <- value:
		return true
	}
}

// SendWithTimeout will send value to ch, blocking until the send succeeds, the
// given context is canceled, or the timeout is reached. Returns whether the
// send was successful.
func SendWithTimeout[T any](
	ctx context.Context,
	ch chan<- T,
	value T,
	timeout time.Duration,
) bool {
	if ch == nil || ctx.Err() != nil {
		return false
	}

	timer := _newTimer(timeout)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return false
	case ch <- value:
		return true
	}
}

// Recv will receive from ch, blocking until either the receive succeeds or the
// given context is canceled. Returns any value received and whether the
// receive was successful.
func Recv[T any](ctx context.Context, ch <-chan T) (value T, ok bool) {
	select {
	case <-ctx.Done():
		return value, ok
	case value, ok = <-ch:
		return value, ok
	}
}

// RecvWithTimeout will receive from ch, blocking until the receive succeeds,
// the given context is canceled, or the timeout is reached. Returns any value
// received and whether the receive was successful.
func RecvWithTimeout[T any](
	ctx context.Context,
	ch <-chan T,
	timeout time.Duration,
) (value T, ok bool) {
	timer := _newTimer(timeout)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return value, ok
	case <-timer.C:
		return value, ok
	case value, ok = <-ch:
		return value, ok
	}
}
