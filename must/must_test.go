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

func TestError(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.Error[%T]() given a non-nil error: %s",
			0,
			t.Name(),
		)
		require.PanicsWithError(t, want, func() {
			must.Error(0, errors.New(t.Name()))
		})
	})

	t.Run("no panic", func(t *testing.T) {
		want := 123
		require.NotPanics(t, func() {
			require.Equal(t, want, must.Error(want, nil))
		})
	})
}

func TestBool(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.Bool[%T]() given a false boolean value",
			0,
		)
		require.PanicsWithError(t, want, func() {
			must.Bool(0, false)
		})
	})

	t.Run("no panic", func(t *testing.T) {
		require.NotPanics(t, func() {
			require.Equal(t, 123, must.Bool(123, true))
		})
	})
}

func TestPredicate(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.Bool[%T]() given a false boolean value",
			0,
		)
		require.PanicsWithError(t, want, func() {
			must.Predicate(0, func(int) bool {
				return false
			})
		})
	})

	t.Run("no panic", func(t *testing.T) {
		require.NotPanics(t, func() {
			must.Predicate(0, func(int) bool {
				return true
			})
		})
	})
}

func TestMustFunc(t *testing.T) {
	t.Run("panic error", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.Error[%T]() given a non-nil error: %s",
			0,
			t.Name(),
		)
		require.PanicsWithError(t, want, func() {
			must.Func[int](func() (int, error) {
				return 0, errors.New(t.Name())
			})
		})
	})

	t.Run("panic bool", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.Bool[%T]() given a false boolean value",
			0,
		)
		require.PanicsWithError(t, want, func() {
			must.Func[int](func() (int, bool) {
				return 0, false
			})
		})
	})

	t.Run("no panic error", func(t *testing.T) {
		want := 123
		require.NotPanics(t, func() {
			require.Equal(t, want, must.Func[int](func() (int, error) {
				return want, nil
			}))
		})
	})

	t.Run("no panic bool", func(t *testing.T) {
		want := 123
		require.NotPanics(t, func() {
			require.Equal(t, want, must.Func[int](func() (int, bool) {
				return want, true
			}))
		})
	})
}

func TestAny(t *testing.T) {
	cases := map[string]struct {
		givePanic   any
		giveNoPanic any
	}{
		"error": {
			givePanic:   errors.New(t.Name()),
			giveNoPanic: error(nil),
		},
		"bool": {
			givePanic:   false,
			giveNoPanic: true,
		},
		"func T bool": {
			givePanic:   func(int) bool { return false },
			giveNoPanic: func(int) bool { return true },
		},
		"func bool": {
			givePanic:   func() bool { return false },
			giveNoPanic: func() bool { return true },
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Panics(t, func() {
				must.Any(0, tt.givePanic)
			})
			require.NotPanics(t, func() {
				must.Any(0, tt.giveNoPanic)
			})
		})
	}
}

func TestDo(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		want := fmt.Sprintf(
			"must: condition failed: must.Do() received a non-nil error: %s",
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
			"must: condition failed: must.Close() received a non-nil error when closing: %s",
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
