package os_test

import (
	"bytes"
	"errors"
	"io"
	goos "os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/os"
)

func TestWriteReaderToFile(t *testing.T) {
	var (
		errReader = errors.New("reader error")
		errRead   = errors.New("read error")
		errClose  = errors.New("close error")
	)

	cases := map[string]struct {
		newReader  func(t *testing.T) any
		wantErr    error
		wantCreate bool
		wantWrite  bool
	}{
		"io.Reader nominal": {
			newReader: func(t *testing.T) any {
				return bytes.NewBuffer([]byte(t.Name()))
			},
			wantErr:    nil,
			wantCreate: true,
			wantWrite:  true,
		},
		"io.Reader read error": {
			newReader: func(t *testing.T) any {
				return erroringReader{
					err: errRead,
				}
			},
			wantErr:    errRead,
			wantCreate: true,
			wantWrite:  false,
		},
		"io.ReadCloser nominal": {
			newReader: func(t *testing.T) any {
				return io.NopCloser(bytes.NewBuffer([]byte(t.Name())))
			},
			wantErr:    nil,
			wantCreate: true,
			wantWrite:  true,
		},
		"io.ReadCloser read error": {
			newReader: func(t *testing.T) any {
				return erroringReadCloser{
					Reader: erroringReader{
						err: errRead,
					},
					Closer: nopReadCloser{},
				}
			},
			wantErr:    errRead,
			wantCreate: true,
			wantWrite:  false,
		},
		"io.ReadCloser close error": {
			newReader: func(t *testing.T) any {
				return erroringReadCloser{
					Reader: bytes.NewBuffer([]byte(t.Name())),
					Closer: erroringCloser{
						err: errClose,
					},
				}
			},
			wantErr:    errClose,
			wantCreate: true,
			wantWrite:  true,
		},
		"ReaderFunc nominal": {
			newReader: func(t *testing.T) any {
				return func() (io.Reader, error) {
					return bytes.NewBuffer([]byte(t.Name())), nil
				}
			},
			wantErr:    nil,
			wantCreate: true,
			wantWrite:  true,
		},
		"ReaderFunc error": {
			newReader: func(t *testing.T) any {
				return func() (io.Reader, error) {
					return nil, errReader
				}
			},
			wantErr:    errReader,
			wantCreate: false,
			wantWrite:  false,
		},
		"ReaderFunc read error": {
			newReader: func(t *testing.T) any {
				return func() (io.Reader, error) {
					return erroringReader{
						err: errRead,
					}, nil
				}
			},
			wantErr:    errRead,
			wantCreate: true,
			wantWrite:  false,
		},
		"ReadCloserFunc nominal": {
			newReader: func(t *testing.T) any {
				return func() (io.ReadCloser, error) {
					return io.NopCloser(bytes.NewBuffer([]byte(t.Name()))), nil
				}
			},
			wantErr:    nil,
			wantCreate: true,
			wantWrite:  true,
		},
		"ReadCloserFunc error": {
			newReader: func(t *testing.T) any {
				return func() (io.ReadCloser, error) {
					return nil, errReader
				}
			},
			wantErr:    errReader,
			wantCreate: false,
			wantWrite:  false,
		},
		"ReadCloserFunc read error": {
			newReader: func(t *testing.T) any {
				return func() (io.ReadCloser, error) {
					return erroringReadCloser{
						Reader: erroringReader{
							err: errRead,
						},
						Closer: nopReadCloser{},
					}, nil
				}
			},
			wantErr:    errRead,
			wantCreate: true,
			wantWrite:  false,
		},
		"ReadCloserFunc close error": {
			newReader: func(t *testing.T) any {
				return func() (io.ReadCloser, error) {
					return erroringReadCloser{
						Reader: bytes.NewBuffer([]byte(t.Name())),
						Closer: erroringCloser{
							err: errClose,
						},
					}, nil
				}
			},
			wantErr:    errClose,
			wantCreate: true,
			wantWrite:  true,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				tempdir = t.TempDir()
				dst     = filepath.Join(tempdir, "testfile")
				size    = len(t.Name())
			)

			n, err := os.WriteReaderToFile(dst, tt.newReader(t))
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.wantWrite, size == n)

			stat, statErr := goos.Stat(dst)
			if !tt.wantCreate {
				require.ErrorIs(t, statErr, goos.ErrNotExist)
				return
			}
			require.NoError(t, statErr)
			require.EqualValues(t, 0o644, stat.Mode()&goos.ModePerm)

			raw, readErr := goos.ReadFile(dst)
			require.NoError(t, readErr)
			if tt.wantWrite {
				require.Equal(t, t.Name(), string(raw))
			} else {
				require.Empty(t, raw)
			}
		})
	}
}

func TestWriteReaderToFile_BadPath(t *testing.T) {
	_, err := os.WriteReaderToFile(
		filepath.Join("-", t.Name(), "-", "does", "not", "exist"),
		nopReadCloser{},
	)
	require.ErrorIs(t, err, goos.ErrNotExist)
}

func TestWriteReaderToFile_BadReader(t *testing.T) {
	readers := []any{nil, "foo", 123, io.Writer(nil)}
	for _, reader := range readers {
		_, err := os.WriteReaderToFile("", reader)
		require.ErrorIs(t, err, os.ErrUnsupportedWriteReader)
	}
}

type nopReadCloser struct {
	erroringReader
	erroringCloser
}

type erroringReader struct {
	err error
}

func (e erroringReader) Read([]byte) (n int, err error) {
	if err = e.err; err == nil {
		err = io.EOF
	}
	return
}

type erroringCloser struct {
	err error
}

func (c erroringCloser) Close() error {
	return c.err
}

type erroringReadCloser struct {
	io.Reader
	io.Closer
}
