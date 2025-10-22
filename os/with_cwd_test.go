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

package os_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	xos "go.mway.dev/x/os"
)

func TestWithCwd(t *testing.T) {
	orig := resolvedTempDir(t)
	require.NoError(t, os.Chdir(orig))
	requireCwd(t, orig)

	var (
		newdir = resolvedTempDir(t)
		err    = xos.WithCwd(newdir, func() {
			requireCwd(t, newdir)
		})
	)

	require.NoError(t, err)
	requireCwd(t, orig)
}

func TestWithCwdSameDir(t *testing.T) {
	orig := resolvedTempDir(t)
	require.NoError(t, os.Chdir(orig))
	requireCwd(t, orig)

	err := xos.WithCwd(orig, func() {
		requireCwd(t, orig)
	})

	require.NoError(t, err)
	requireCwd(t, orig)
}

func TestWithCwdPanic(t *testing.T) {
	orig := resolvedTempDir(t)
	require.NoError(t, os.Chdir(orig))
	requireCwd(t, orig)

	var (
		newdir = resolvedTempDir(t)
		err    error
	)

	require.Panics(t, func() {
		err = xos.WithCwd(newdir, func() {
			requireCwd(t, newdir)
			panic("oops")
		})
	})
	require.NoError(t, err)
	requireCwd(t, orig)
}

func TestWithCwdOrigDirGone(t *testing.T) {
	orig := resolvedTempDir(t)
	require.NoError(t, os.Chdir(orig))
	requireCwd(t, orig)

	var (
		newdir = resolvedTempDir(t)
		err    = xos.WithCwd(newdir, func() {
			requireCwd(t, newdir)
			require.NoError(t, os.Remove(orig))
		})
	)

	require.Error(t, err)
	requireCwd(t, newdir)
}

func TestWithCwdOrigDirRenamed(t *testing.T) {
	orig := resolvedTempDir(t)
	require.NoError(t, os.Chdir(orig))
	requireCwd(t, orig)

	var (
		moved  = orig + "-moved"
		newdir = resolvedTempDir(t)
		err    = xos.WithCwd(newdir, func() {
			requireCwd(t, newdir)
			require.NoError(t, os.Rename(orig, moved))
		})
	)

	require.Error(t, err)
	requireCwd(t, newdir)
}

func TestWithCwdEmptyTargetDir(t *testing.T) {
	err := xos.WithCwd("", func() {
		require.FailNow(t, "WithCwd func argument should not be called")
	})
	require.Error(t, err)
}

func TestWithCwdTargetDirDoesNotExist(t *testing.T) {
	bad := []string{
		"/foo/bar/baz",
		"            ",
		"!@#$%^&*()",
		"`\x00",
	}

	for _, dir := range bad {
		err := xos.WithCwd(dir, func() {
			require.FailNow(t, "WithCwd func argument should not be called")
		})
		require.Error(t, err)
	}
}

func TestWithCwdErrorFunc(t *testing.T) {
	orig := resolvedTempDir(t)
	require.NoError(t, os.Chdir(orig))
	requireCwd(t, orig)

	var (
		wantErr = errors.New(t.Name())
		haveErr = xos.WithCwd(orig, func() error {
			requireCwd(t, orig)
			return wantErr
		})
	)
	require.ErrorIs(t, haveErr, wantErr)
}

func resolvedTempDir(t *testing.T) string {
	dir, err := filepath.Abs(t.TempDir())
	require.NoError(t, err)

	dir, err = filepath.EvalSymlinks(dir)
	require.NoError(t, err)

	return dir
}

func requireCwd(t *testing.T, dir string) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	require.Equal(t, dir, cwd)
}
