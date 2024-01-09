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
package os

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/errors"
	"go.mway.dev/x/stub"
)

func TestMkdirAllInheritStatError(t *testing.T) {
	var (
		tempdir  = t.TempDir()
		statErr  = errors.New("stat error")
		statStub = func(string) (os.FileInfo, error) { //nolint:unparam
			return nil, statErr
		}
	)

	err := WithCwd(tempdir, func() {
		stub.With(&_osStat, statStub, func() {
			require.ErrorIs(t, MkdirAllInherit("foo"), statErr)
		})
	})
	require.NoError(t, err)
}

func TestMkdirAllInheritMkdirAllError(t *testing.T) {
	var (
		tempdir      = t.TempDir()
		mkdirAllErr  = errors.New("mkdir error")
		mkdirAllStub = func(string, os.FileMode) error {
			return mkdirAllErr
		}
	)

	err := WithCwd(tempdir, func() {
		stub.With(&_osMkdirAll, mkdirAllStub, func() {
			require.ErrorIs(t, MkdirAllInherit("foo"), mkdirAllErr)
		})
	})
	require.NoError(t, err)
}
