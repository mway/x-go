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

package http

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/errors"
	"go.mway.dev/x/os/tempdir"
	"go.mway.dev/x/stub"
)

func TestGetFile_Nominal(t *testing.T) {
	var (
		giveURL     = "https://foo.bar/baz.bat"
		wantFile    = "baz.bat"
		wantContent = t.Name()
	)

	err := tempdir.With(func(dst string) {
		stub.With(&_httpGet, newHTTPGetFunc(t, giveURL, t.Name(), nil), func() {
			dst = filepath.Join(dst, wantFile)
			require.NoError(t, GetFile(giveURL, dst))

			raw, err := os.ReadFile(dst)
			require.NoError(t, err)
			require.Equal(t, wantContent, string(raw))
		})
	})
	require.NoError(t, err)
}

func TestGetFile_HTTPGetError(t *testing.T) {
	wantErr := errors.New(t.Name())
	stub.With(&_httpGet, newHTTPGetFunc(t, "http://foo", "", wantErr), func() {
		require.ErrorIs(t, GetFile("http://foo", ""), wantErr)
	})
}

func TestGetFile_HTTPResponseBodyCloseError(t *testing.T) {
	var (
		giveURL     = "https://foo.bar/baz.bat"
		wantFile    = "baz.bat"
		wantContent = t.Name()
		wantErr     = errors.New(t.Name())
		newRequest  = func(string) (*http.Response, error) { //nolint:unparam
			resp, err := newHTTPGetFunc(t, giveURL, t.Name(), nil)(giveURL)
			require.NoError(t, err)

			resp.Body = testReader{
				reader:   resp.Body,
				closeErr: wantErr,
			}
			return resp, nil
		}
	)

	err := tempdir.With(func(dst string) {
		stub.With(&_httpGet, newRequest, func() {
			dst = filepath.Join(dst, wantFile)
			require.ErrorIs(t, GetFile(giveURL, dst), wantErr)

			raw, err := os.ReadFile(dst)
			require.NoError(t, err)
			require.Equal(t, wantContent, string(raw))
		})
	})
	require.NoError(t, err)
}

func TestGetFile_OSStatError(t *testing.T) {
	var (
		wantStatErr = errors.New(t.Name())
		osStat      = func(string) (fs.FileInfo, error) { //nolint:unparam
			return nil, wantStatErr
		}
		giveURL  = "https://foo.bar/baz.bat"
		wantFile = "baz.bat"
	)

	err := tempdir.With(func(dst string) {
		stub.With(&_httpGet, newHTTPGetFunc(t, giveURL, t.Name(), nil), func() {
			stub.With(&_osStat, osStat, func() {
				dst = filepath.Join(dst, wantFile)
				require.ErrorIs(t, GetFile(giveURL, dst), wantStatErr)
			})
		})
	})
	require.NoError(t, err)
}

func TestGetFile_OSMkdirAllError(t *testing.T) {
	var (
		wantMkdirError = errors.New(t.Name())
		mkdirAll       = func(string, fs.FileMode) error {
			return wantMkdirError
		}
		giveURL  = "https://foo.bar/baz.bat"
		wantFile = "baz.bat"
	)

	err := tempdir.With(func(dst string) {
		stub.With(&_httpGet, newHTTPGetFunc(t, giveURL, t.Name(), nil), func() {
			stub.With(&_osMkdirAll, mkdirAll, func() {
				dst = filepath.Join(dst, "foo", wantFile)
				require.ErrorIs(t, GetFile(giveURL, dst), wantMkdirError)
			})
		})
	})
	require.NoError(t, err)
}

func TestGetFile_DestNotWritable(t *testing.T) {
	var (
		unixAccess = func(string, uint32) error {
			return errors.New(t.Name())
		}
		giveURL  = "https://foo.bar/baz.bat"
		wantFile = "baz.bat"
	)

	err := tempdir.With(func(dst string) {
		stub.With(&_httpGet, newHTTPGetFunc(t, giveURL, t.Name(), nil), func() {
			stub.With(&_unixAccess, unixAccess, func() {
				dst = filepath.Join(dst, "foo", wantFile)
				require.ErrorIs(t, GetFile(giveURL, dst), ErrDestNotWritable)
			})
		})
	})
	require.NoError(t, err)
}

func TestGetFile_DestIsDir(t *testing.T) {
	var (
		giveURL  = "https://foo.bar/baz.bat"
		wantFile = "baz.bat"
	)

	err := tempdir.With(func(dst string) {
		stub.With(&_httpGet, newHTTPGetFunc(t, giveURL, t.Name(), nil), func() {
			dst = filepath.Join(dst, "foo", wantFile)
			dir := filepath.Dir(dst)
			require.NoError(t, os.Mkdir(dir, 0o755))
			require.NoError(t, GetFile(giveURL, dir))

			_, err := os.ReadFile(dst)
			require.NoError(t, err)
		})
	})
	require.NoError(t, err)
}

func newHTTPGetFunc(
	t *testing.T,
	wantURL string,
	contents string,
	err error,
) func(string) (*http.Response, error) {
	return func(url string) (*http.Response, error) {
		require.Equal(t, wantURL, url)

		if err != nil {
			return nil, err
		}

		buf := bytes.NewBufferString(contents)
		return &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Proto:         "https",
			Body:          io.NopCloser(buf),
			ContentLength: int64(buf.Len()),
		}, nil
	}
}

type testReader struct {
	reader   io.ReadCloser
	readErr  error
	closeErr error
}

func (e testReader) Read(p []byte) (int, error) {
	switch {
	case e.readErr != nil:
		return 0, e.readErr
	case e.reader == nil:
		return 0, io.EOF
	default:
		return e.reader.Read(p)
	}
}

func (e testReader) Close() error {
	err := e.reader.Close()
	if e.closeErr != nil {
		err = e.closeErr
	}
	return err
}
