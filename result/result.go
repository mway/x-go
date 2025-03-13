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

// Package result provides a helpful Result[T,error] type wrapper.
package result

const (
	_value = 1
	_err   = 2
)

// A Result holds either a value of type T or an error.
type Result[T any] struct {
	value T
	err   error
	flag  uint8
}

// Ok constructs a [Result] holding the given value.
func Ok[T any](value T) Result[T] {
	return Result[T]{
		value: value,
		flag:  _value,
	}
}

// Err constructs a [Result] holding the given error.
func Err[T any](err error) Result[T] {
	return Result[T]{
		err:  err,
		flag: _err,
	}
}

// HasValue returns whether the result holds a value.
func (r *Result[T]) HasValue() bool {
	return r.flag == _value
}

// Value returns the result's held value, if any. The boolean indicates whether
// a value is held.
func (r *Result[T]) Value() (T, bool) {
	if r.flag != _value {
		var zero T
		return zero, false
	}
	return r.value, true
}

// ValueOr returns the result's held value, if any; the given fallback is
// returned otherwise.
func (r *Result[T]) ValueOr(fallback T) T {
	x, ok := r.Value()
	if !ok {
		x = fallback
	}
	return x
}

// ValueOrElse returns the result's held value, if any; the given function is
// called to produce and return a value otherwise.
func (r *Result[T]) ValueOrElse(fn func() T) T {
	x, ok := r.Value()
	if !ok {
		x = fn()
	}
	return x
}

// HasErr returns whether the result holds an error.
func (r *Result[T]) HasErr() bool {
	return r.flag == _err
}

// Err returns the result's held error, if any. The boolean indicates whether
// an error is held.
func (r *Result[T]) Err() (error, bool) { //nolint:revive
	if r.flag != _err {
		return nil, false
	}
	return r.err, true
}

// ErrOr returns the result's held error, if any; the given fallback is
// returned otherwise.
func (r *Result[T]) ErrOr(fallback error) error {
	x, ok := r.Err()
	if !ok {
		x = fallback
	}
	return x
}

// ErrOrElse returns the result's held error, if any; the given function is
// called to produce and return a value otherwise.
func (r *Result[T]) ErrOrElse(fn func() error) error {
	x, ok := r.Err()
	if !ok {
		x = fn()
	}
	return x
}
