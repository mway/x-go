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
