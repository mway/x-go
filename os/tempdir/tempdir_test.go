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

package tempdir_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/errors"

	"go.mway.dev/x/os/tempdir"
)

func TestWith_PathFunc(t *testing.T) {
	haveErr := tempdir.With(func(path string) {
		requireInValidTempDir(t, path)
	})
	require.NoError(t, haveErr)
}

func TestWith_PathErrorFunc(t *testing.T) {
	var (
		wantErr = errors.New(t.Name())
		haveErr = tempdir.With(func(path string) error {
			requireInValidTempDir(t, path)
			return wantErr
		})
	)
	require.ErrorIs(t, haveErr, wantErr)
}

func TestNewWithBase_BadTempDir(t *testing.T) {
	_, err := tempdir.NewWithBase("-/does/not/exist/-")
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestDir_Nominal(t *testing.T) {
	d, err := tempdir.New()
	require.NoError(t, err)

	var called bool
	require.NoError(t, d.With(func(_ string) error {
		called = true
		return nil
	}))
	require.True(t, called)
}

func TestDir_NoCallsAfterClose(t *testing.T) {
	d, err := tempdir.New()
	require.NoError(t, err)
	require.NoError(t, d.Close())

	var called bool
	require.NoError(t, d.With(func(_ string) error {
		called = true
		return nil
	}))
	require.False(t, called)

	require.NoError(t, d.Close()) // double close
}

func TestDir_InvalidFunc(t *testing.T) {
	d, err := tempdir.New()
	require.NoError(t, err)

	require.ErrorIs(t, d.With(nil), tempdir.ErrInvalidFuncType)
	require.ErrorIs(t, d.With("foo"), tempdir.ErrInvalidFuncType)
	require.ErrorIs(
		t,
		d.With(func() { /* bad */ }),
		tempdir.ErrInvalidFuncType,
	)
}

func requireInValidTempDir(t *testing.T, path string) {
	t.Helper()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	absPath, err := filepath.EvalSymlinks(path)
	require.NoError(t, err)
	require.Equal(t, cwd, absPath)

	absTemp, err := filepath.EvalSymlinks(os.TempDir())
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(absPath, absTemp))
}
