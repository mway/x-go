package os

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithCwdGetwdError(t *testing.T) {
	var (
		expectErr = errors.New("os.Getwd error")
		prevwd    = _getwd
	)
	defer func() {
		_getwd = prevwd
	}()
	_getwd = func() (string, error) { return "", expectErr }

	require.NotPanics(t, func() {
		err := WithCwd(".", func() {
			panic("bad")
		})
		require.ErrorIs(t, err, expectErr)
	})
}
