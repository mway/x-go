package os_test

import (
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
