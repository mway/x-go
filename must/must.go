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

// Package must provides helpers that narrow operands or execute functions
// under threat of panic.
package must

import (
	"errors"
	"fmt"
	"io"
)

// ErrConditionFailed is a base error used in must package function panics, in
// order to allow any rescuers to understand that the panic being recovered was
// due to this package.
var ErrConditionFailed = errors.New("must: condition failed")

// Get returns value if err is nil, or panics with err otherwise.
func Get[T any](value T, err error) T {
	if err != nil {
		panic(fmt.Errorf(
			"%w: must.Get[%T]: received a non-nil error: %w",
			ErrConditionFailed,
			value,
			err,
		))
	}
	return value
}

// True returns value if ok is true, or panics otherwise.
func True[T any](value T, ok bool) T {
	if !ok {
		panic(fmt.Errorf(
			"%w: must.True[%T]: received a false value",
			ErrConditionFailed,
			value,
		))
	}
	return value
}

// False returns value if ok is false, or panics otherwise.
func False[T any](value T, ok bool) T {
	if ok {
		panic(fmt.Errorf(
			"%w: must.False[%T]: received a true value",
			ErrConditionFailed,
			value,
		))
	}
	return value
}

// A PredicateFunc returns whether a value should be considered valid.
type PredicateFunc[T any] interface {
	~func(T) bool | ~func() bool
}

// Predicate returns value if pred returns true, or panics otherwise.
func Predicate[T any, P PredicateFunc[T]](value T, pred P) T {
	type (
		FnT = func(T) bool
		Fn  = func() bool
	)
	if fn, ok := any(pred).(FnT); ok {
		return True(value, fn(value))
	}
	return True(value, any(pred).(Fn)()) //nolint:errcheck
}

// Any delegates to [Must], [MustBool], or [MustPredicate], depending on P,
// returning the result.
func Any[T any, P comparable](value T, pred P) T {
	if err, ok := any(pred).(error); ok {
		return Get(value, err)
	}

	switch x := any(pred).(type) {
	case bool:
		return True(value, x)
	case func(T) bool:
		return Predicate(value, x)
	case func() bool:
		return Predicate(value, x)
	default:
		var zero P
		return True(value, pred == zero)
	}
}

// A MustFunc is any function that returns a T and either an error or a bool.
//
//nolint:revive
type MustFunc[T any] interface {
	~func() (T, error) | ~func() (T, bool)
}

// Func calls Must with the values returned by fn.
func Func[T any, F MustFunc[T]](fn F) T {
	type (
		fnError = func() (T, error)
		fnBool  = func() (T, bool)
	)
	if fn, ok := any(fn).(fnError); ok {
		return Get(fn())
	}
	return True(any(fn).(fnBool)()) //nolint:errcheck
}

// Do panics if fn returns an error.
func Do(fn func() error) {
	if err := fn(); err != nil {
		panic(fmt.Errorf(
			"%w: must.Do: received a non-nil error: %w",
			ErrConditionFailed,
			err,
		))
	}
}

// Close panics if closer.Close returns an error.
func Close(closer io.Closer) {
	if err := closer.Close(); err != nil {
		panic(fmt.Errorf(
			"%w: must.Close: received a non-nil error when closing: %w",
			ErrConditionFailed,
			err,
		))
	}
}

// NotError returns value if err is nil, or panics with err otherwise.
func NotError(fn func() error) {
	if err := fn(); err != nil {
		panic(fmt.Errorf(
			"%w: must.NotError: fn produced a non-nil error: %w",
			ErrConditionFailed,
			err,
		))
	}
}

// As attempts to convert the given value from type [In] to type [Out] and
// return the result. if the value cannot be converted, As will panic.
func As[Out any, In any](in In) Out {
	out, ok := any(in).(Out)
	if !ok {
		panic(fmt.Errorf(
			"%w: must.As: cannot convert from %T to %T",
			ErrConditionFailed,
			in,
			out,
		))
	}
	return out
}
