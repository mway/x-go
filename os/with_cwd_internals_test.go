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

package os

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithCwdFixedCwd(t *testing.T) {
	prevwd := _getwd
	defer func() {
		_getwd = prevwd
	}()
	_getwd = func() (string, error) { return "sometestdir", nil }

	var called bool
	err := WithCwd("sometestdir", func() {
		called = true
	})
	require.NoError(t, err)
	require.True(t, called)
}

func TestWithCwdGetwdError(t *testing.T) {
	var (
		expectErr = errors.New("os.Getwd error")
		prevwd    = _getwd
	)
	defer func() {
		_getwd = prevwd
	}()
	_getwd = func() (string, error) { return "", expectErr }

	err := WithCwd(".", func() {
		require.FailNow(t, "WithCwd func argument should not be called")
	})
	require.ErrorIs(t, err, expectErr)
}

func TestWithCwdChdirError(t *testing.T) {
	var (
		expectErr = errors.New("os.Chdir error")
		prevChdir = _chdir
	)
	defer func() {
		_chdir = prevChdir
	}()
	_chdir = func(string) error { return expectErr }

	require.NotPanics(t, func() {
		err := WithCwd(".", func() {
			require.FailNow(t, "WithCwd func argument should not be called")
		})
		require.ErrorIs(t, err, expectErr)
	})
}
