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

package context

import (
	"context"
	"time"
)

type (
	// A Context carries a deadline, a cancellation signal, and other values
	// across API boundaries.
	//
	// Context's methods may be called by multiple goroutines simultaneously.
	Context = context.Context

	// A CancelFunc tells an operation to abandon its work. A CancelFunc does
	// not wait for the work to stop. A CancelFunc may be called by multiple
	// goroutines simultaneously. After the first call, subsequent calls to a
	// CancelFunc do nothing.
	CancelFunc = context.CancelFunc
)

var (
	// Canceled is the error returned by Context.Err when the context is
	// canceled.
	Canceled = context.Canceled

	// DeadlineExceeded is the error returned by Context.Err when the context's
	// deadline passes.
	DeadlineExceeded = context.DeadlineExceeded
)

// Background returns a non-nil, empty Context. It is never canceled, has no
// values, and has no deadline. It is typically used by the main function,
// initialization, and tests, and as the top-level Context for incoming
// requests.
func Background() context.Context {
	return context.Background()
}

// TODO returns a non-nil, empty Context. Code should use context.TODO when
// it's unclear which Context to use or it is not yet available (because the
// surrounding function has not yet been extended to accept a Context
// parameter).
func TODO() context.Context {
	return context.TODO()
}

// WithCancel returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed when the returned cancel function is called
// or when the parent context's Done channel is closed, whichever happens
// first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func WithCancel(ctx Context) (Context, CancelFunc) {
	return context.WithCancel(ctx)
}

// WithDeadline returns a copy of the parent context with the deadline adjusted
// to be no later than d. If the parent's deadline is already earlier than d,
// WithDeadline(parent, d) is semantically equivalent to parent. The returned
// context's Done channel is closed when the deadline expires, when the
// returned cancel function is called, or when the parent context's Done
// channel is closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func WithDeadline(ctx Context, t time.Time) (Context, CancelFunc) {
	return context.WithDeadline(ctx, t)
}

// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete:
//
// 	func slowOperationWithTimeout(ctx context.Context) (Result, error) {
// 		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
// 		defer cancel()  // releases resources if slowOperation completes before
// 		                // timeout elapses
// 		return slowOperation(ctx)
// 	}
func WithTimeout(ctx context.Context, d time.Duration) (Context, CancelFunc) {
	return context.WithTimeout(ctx, d)
}

// WithValue returns a copy of parent in which the value associated with key is
// val.
//
// Use context Values only for request-scoped data that transits processes and
// APIs, not for passing optional parameters to functions.
//
// The provided key must be comparable and should not be of type string or any
// other built-in type to avoid collisions between packages using context.
// Users of WithValue should define their own types for keys. To avoid
// allocating when assigning to an interface{}, context keys often have
// concrete type struct{}. Alternatively, exported context key variables'
// static type should be a pointer or interface.
func WithValue(ctx Context, key any, value any) Context {
	return context.WithValue(ctx, key, value)
}
