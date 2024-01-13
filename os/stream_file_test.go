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

package os_test

import (
	"io"
	"io/fs"
	goos "os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/os"
	"go.mway.dev/x/os/tempdir"
)

func TestWithFileReader(t *testing.T) {
	err := tempdir.With(func(_ string) {
		wantContent := []byte("hello")
		require.NoError(t, goos.WriteFile(t.Name(), wantContent, 0o744))

		var haveContent []byte
		require.NoError(
			t,
			os.WithFileReader(
				t.Name(),
				func(r io.Reader) (err error) {
					haveContent, err = io.ReadAll(r)
					return
				},
			),
		)
		require.Equal(t, wantContent, haveContent)
	})
	require.NoError(t, err)
}

func TestWithFileReader_NotExist(t *testing.T) {
	err := tempdir.With(func(_ string) {
		var haveContent []byte
		require.ErrorIs(
			t,
			os.WithFileReader(
				t.Name(),
				func(r io.Reader) (err error) {
					haveContent, err = io.ReadAll(r)
					return
				},
			),
			goos.ErrNotExist,
		)
		require.Nil(t, haveContent)
	})
	require.NoError(t, err)
}

func TestWithFileWriter(t *testing.T) {
	err := tempdir.With(func(_ string) {
		wantContent := []byte("hello")

		require.NoError(
			t,
			os.WithFileWriter(
				t.Name(),
				func(w io.Writer) error {
					_, err := w.Write(wantContent)
					return err
				},
			),
		)

		haveContent, err := goos.ReadFile(t.Name())
		require.NoError(t, err)
		require.Equal(t, wantContent, haveContent)
	})
	require.NoError(t, err)
}

func TestWithFileWriter_NotExist(t *testing.T) {
	err := tempdir.With(func(_ string) {
		require.ErrorIs(
			t,
			os.WithFileWriter(
				filepath.Join("dne", t.Name()),
				func(w io.Writer) error {
					_, err := w.Write([]byte("won't be written"))
					return err
				},
			),
			goos.ErrNotExist,
		)

		_, err := goos.ReadFile(t.Name())
		require.ErrorIs(t, err, goos.ErrNotExist)
	})
	require.NoError(t, err)
}

func TestWithFileModeWriter(t *testing.T) {
	err := tempdir.With(func(_ string) {
		var (
			wantContent = []byte("hello")
			wantMode    = fs.FileMode(0o741) // intentionally unusual permission
		)

		require.NoError(
			t,
			os.WithFileModeWriter(
				t.Name(),
				wantMode,
				func(w io.Writer) error {
					_, err := w.Write(wantContent)
					return err
				},
			),
		)

		haveContent, err := goos.ReadFile(t.Name())
		require.NoError(t, err)
		require.Equal(t, wantContent, haveContent)

		haveStat, err := goos.Stat(t.Name())
		require.NoError(t, err)
		require.Equal(t, wantMode.String(), (haveStat.Mode() & fs.ModePerm).String())
	})
	require.NoError(t, err)
}
