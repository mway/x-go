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

// Error returns value if err is nil, or panics with err otherwise.
func Error[T any](value T, err error) T {
	if err != nil {
		panic(fmt.Errorf(
			"%w: must.Error[%T]() given a non-nil error: %w",
			ErrConditionFailed,
			value,
			err,
		))
	}
	return value
}

// Bool returns value if ok is true, or panics otherwise.
func Bool[T any](value T, ok bool) T {
	if !ok {
		panic(fmt.Errorf(
			"%w: must.Bool[%T]() given a false boolean value",
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
		return Bool(value, fn(value))
	}
	return Bool(value, any(pred).(Fn)())
}

// Any delegates to [Must], [MustBool], or [MustPredicate], depending on P,
// returning the result.
func Any[T any, P comparable](value T, pred P) T {
	if err, ok := any(pred).(error); ok {
		return Error(value, err)
	}

	switch x := any(pred).(type) {
	case bool:
		return Bool(value, x)
	case func(T) bool:
		return Predicate(value, x)
	case func() bool:
		return Predicate(value, x)
	default:
		var zero P
		return Bool(value, pred == zero)
	}
}

// A MustFunc is any function that returns a T and either an error or a bool.
//
//nolint:revive
type MustFunc[T any] interface {
	~func() (T, error) | ~func() (T, bool)
}

// Func calls Must with the values returned by fn.
//
//nolint:revive
func Func[T any, F MustFunc[T]](fn F) T {
	type (
		fnError = func() (T, error)
		fnBool  = func() (T, bool)
	)
	if fn, ok := any(fn).(fnError); ok {
		return Error(fn())
	}
	return Bool(any(fn).(fnBool)())
}

// Do panics if fn returns an error.
func Do(fn func() error) {
	if err := fn(); err != nil {
		panic(fmt.Errorf(
			"%w: must.Do() received a non-nil error: %w",
			ErrConditionFailed,
			err,
		))
	}
}

// Close panics if closer.Close returns an error.
func Close(closer io.Closer) {
	if err := closer.Close(); err != nil {
		panic(fmt.Errorf(
			"%w: must.Close() received a non-nil error when closing: %w",
			ErrConditionFailed,
			err,
		))
	}
}
