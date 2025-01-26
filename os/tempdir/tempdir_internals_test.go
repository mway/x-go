// Copyright (c) 2025 Matt Way
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
