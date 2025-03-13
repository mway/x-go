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

// Package future provides future-related types and utilities.
package future

import (
	"context"
	"sync"
)

// A Future is a type that may hold a value now or in the future.
type Future[T any] struct {
	value    T
	rmu      sync.RWMutex
	wmu      sync.Mutex
	wg       sync.WaitGroup
	isset    bool
	canceled bool
}

// New creates a new [Future].
func New[T any]() *Future[T] {
	v := &Future[T]{}
	v.rmu.Lock()
	return v
}

// Cancel cancels the future, if it has not already been canceled or had a
// value set.
func (f *Future[T]) Cancel() {
	f.wmu.Lock()
	defer f.wmu.Unlock()
	if f.isset || f.canceled {
		return
	}
	f.canceled = true
	f.rmu.Unlock()
}

// IsCanceled returns whether the future has been canceled.
func (f *Future[T]) IsCanceled() bool {
	f.wmu.Lock()
	defer f.wmu.Unlock()
	return f.canceled
}

// Set sets the future's value to the given value, if it has not been canceled
// or previously had a value set.
func (f *Future[T]) Set(value T) {
	f.wmu.Lock()
	defer f.wmu.Unlock()
	if f.isset || f.canceled {
		return
	}
	f.isset = true
	f.value = value
	f.rmu.Unlock()
}

// IsSet returns whether the future has a value set.
func (f *Future[T]) IsSet() bool {
	f.wmu.Lock()
	defer f.wmu.Unlock()
	return f.isset
}

// Get returns the future's value, if it currently holds one. The boolean
// indicates whether the returned value is valid (i.e., the returned value was
// explicitly set).
func (f *Future[T]) Get() (T, bool) {
	f.wmu.Lock()
	defer f.wmu.Unlock()
	if !f.isset {
		var zero T
		return zero, false
	}
	return f.value, true
}

// Wait waits for the future to hold a value or to be canceled, and returns its
// value afterward. The boolean indicates whether the returned value is valid
// (i.e., the returned value was explicitly set).
func (f *Future[T]) Wait() (T, bool) {
	f.rmu.RLock()
	defer f.rmu.RUnlock()
	return f.value, f.isset
}

// WaitContext waits for the future to hold a value, to be canceled, or for ctx
// to be canceled, and returns its value afterward. The boolean indicates
// whether the returned value is valid (i.e., the returned value was explicitly
// set).
func (f *Future[T]) WaitContext(ctx context.Context) (T, bool) {
	if x, ok := f.Get(); ok {
		return x, true
	}

	var (
		update = make(chan T, 1)
		value  T
		valid  bool
	)

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		defer close(update)
		if val, ok := f.Wait(); ok {
			update <- val
		}
	}()

	select {
	case <-ctx.Done():
		// n.b. If the context is done but there's a value, prefer the value.
		select {
		case value, valid = <-update:
		default:
		}
	case value, valid = <-update:
	}
	return value, valid
}

// Promise returns a [Promise] that is bound to the future.
func (f *Future[T]) Promise() Promise[T] {
	return Promise[T]{
		future: f,
	}
}

// A Promise is the receiving portion of a [Future].
type Promise[T any] struct {
	future *Future[T]
}

// IsSet returns whether the promise has a value set.
func (p *Promise[T]) IsSet() bool {
	return p.future.IsSet()
}

// IsCanceled returns whether the promise has been canceled.
func (p *Promise[T]) IsCanceled() bool {
	return p.future.IsCanceled()
}

// Get returns the promise's value, if it currently holds one. The boolean
// indicates whether the returned value is valid (i.e., the returned value was
// explicitly set).
func (p *Promise[T]) Get() (T, bool) {
	return p.future.Get()
}

// Wait waits for the promise to hold a value or to be canceled, and returns
// its value afterward. The boolean indicates whether the returned value is
// valid (i.e., the returned value was explicitly set).
func (p *Promise[T]) Wait() (T, bool) {
	return p.future.Wait()
}

// WaitContext waits for the promise to hold a value, to be canceled, or for
// ctx to be canceled, and returns its value afterward. The boolean indicates
// whether the returned value is valid (i.e., the returned value was explicitly
// set).
func (p *Promise[T]) WaitContext(ctx context.Context) (T, bool) {
	return p.future.WaitContext(ctx)
}
