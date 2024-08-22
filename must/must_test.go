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

package must_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/must"
)

func TestMust(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.Must[%T] given a non-nil error: %s",
			0,
			t.Name(),
		)
		require.PanicsWithError(t, want, func() {
			must.Must(0, errors.New(t.Name()))
		})
	})

	t.Run("no panic", func(t *testing.T) {
		want := 123
		require.NotPanics(t, func() {
			require.Equal(t, want, must.Must(want, nil))
		})
	})
}

func TestMustBool(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.MustBool[%T] given a false boolean value",
			0,
		)
		require.PanicsWithError(t, want, func() {
			must.MustBool(0, false)
		})
	})

	t.Run("no panic", func(t *testing.T) {
		require.NotPanics(t, func() {
			require.Equal(t, 123, must.MustBool(123, true))
		})
	})
}

func TestMustFunc(t *testing.T) {
	t.Run("panic error", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.MustFunc[%T] received a non-nil error: %s",
			0,
			t.Name(),
		)
		require.PanicsWithError(t, want, func() {
			must.MustFunc[int](func() (int, error) {
				return 0, errors.New(t.Name())
			})
		})
	})

	t.Run("panic bool", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.MustFunc[%T] received a false boolean value",
			0,
		)
		require.PanicsWithError(t, want, func() {
			must.MustFunc[int](func() (int, bool) {
				return 0, false
			})
		})
	})

	t.Run("no panic error", func(t *testing.T) {
		want := 123
		require.NotPanics(t, func() {
			require.Equal(t, want, must.MustFunc[int](func() (int, error) {
				return want, nil
			}))
		})
	})

	t.Run("no panic bool", func(t *testing.T) {
		want := 123
		require.NotPanics(t, func() {
			require.Equal(t, want, must.MustFunc[int](func() (int, bool) {
				return want, true
			}))
		})
	})
}

func TestDo(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.Do received a non-nil error: %s",
			t.Name(),
		)
		require.PanicsWithError(t, want, func() {
			must.Do(func() error {
				return errors.New(t.Name())
			})
		})
	})

	t.Run("no panic", func(t *testing.T) {
		require.NotPanics(t, func() {
			must.Do(func() error {
				return nil
			})
		})
	})
}

func TestClose(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.Close received a non-nil error when closing: %s",
			t.Name(),
		)
		require.PanicsWithError(t, want, func() {
			must.Close(newTestCloser(errors.New(t.Name())))
		})
	})

	t.Run("no panic", func(t *testing.T) {
		require.NotPanics(t, func() {
			must.Close(newTestCloser(nil))
		})
	})
}

type testCloser struct {
	err    error
	closed bool
}

func newTestCloser(err error) *testCloser {
	return &testCloser{
		err: err,
	}
}

func (c *testCloser) Closed() bool { return c.closed }
func (c *testCloser) Close() error { return c.err }
