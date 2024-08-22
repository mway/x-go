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

// Package must provides helpers that narrow operands or execute functions
// under threat of panic.
package must

import (
	"fmt"
	"io"

	"go.mway.dev/errors"
)

// ErrConditionFailed is a base error used in must package function panics, in
// order to allow any rescuers to understand that the panic being recovered was
// due to this package.
var ErrConditionFailed = errors.New("must: condition failed")

// Must returns value if err is nil, or panics with err otherwise.
func Must[T any](value T, err error) T {
	if err != nil {
		panic(fmt.Errorf(
			"%w: must.Must[%T] given a non-nil error: %w",
			ErrConditionFailed,
			value,
			err,
		))
	}
	return value
}

// MustBool returns value if ok is true, or panics otherwise.
//
//nolint:revive
func MustBool[T any](value T, ok bool) T {
	if !ok {
		panic(fmt.Errorf(
			"%w: must.MustBool[%T] given a false boolean value",
			ErrConditionFailed,
			value,
		))
	}
	return value
}

// A Func is any function that returns a T and either an error or a bool.
type Func[T any] interface {
	~func() (T, error) | ~func() (T, bool)
}

// MustFunc calls Must with the values returned by fn.
//
//nolint:revive
func MustFunc[T any, F Func[T]](fn F) T {
	switch fn := any(fn).(type) {
	case func() (T, error):
		value, err := fn()
		if err != nil {
			panic(fmt.Errorf(
				"%w: must.MustFunc[%T] received a non-nil error: %w",
				ErrConditionFailed,
				value,
				err,
			))
		}
		return value
	case func() (T, bool):
		value, ok := fn()
		if !ok {
			panic(fmt.Errorf(
				"%w: must.MustFunc[%T] received a false boolean value",
				ErrConditionFailed,
				value,
			))
		}
		return value
	default:
		panic("unreachable")
	}
}

// Do panics if fn returns an error.
func Do(fn func() error) {
	if err := fn(); err != nil {
		panic(fmt.Errorf(
			"%w: must.Do received a non-nil error: %w",
			ErrConditionFailed,
			err,
		))
	}
}

// Close panics if closer.Close returns an error.
func Close(closer io.Closer) {
	if err := closer.Close(); err != nil {
		panic(fmt.Errorf(
			"%w: must.Close received a non-nil error when closing: %w",
			ErrConditionFailed,
			err,
		))
	}
}
