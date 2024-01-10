package tempdir

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/errors"
	"go.mway.dev/x/stub"
)

func TestWith_Error(t *testing.T) {
	var (
		wantErr   = errors.New(t.Name())
		mkdirTemp = func(_ string, _ string) (string, error) {
			return "", wantErr
		}
	)

	stub.With(&_osMkdirTemp, mkdirTemp, func() {
		require.ErrorIs(t, With(func(string) { /* nop */ }), wantErr)
	})
}

func TestDir_CloseError(t *testing.T) {
	var (
		wantErr   = errors.New(t.Name())
		removeAll = func(string) error {
			return wantErr
		}
	)

	stub.With(&_osRemoveAll, removeAll, func() {
		d, err := New()
		require.NoError(t, err)
		require.ErrorIs(t, d.Close(), wantErr)
	})
}
