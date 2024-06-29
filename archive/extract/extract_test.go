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

package extract

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/stub"
)

var _fullTree = map[string]string{
	"foo/bar/bat/testfile": "foo/bar/bat/",
	"foo/bar/qux/testfile": "foo/bar/qux/",
	"foo/bar/testfile":     "foo/bar/",
	"foo/baz/bat/testfile": "foo/baz/bat/",
	"foo/baz/qux/testfile": "foo/baz/qux/",
	"foo/baz/testfile":     "foo/baz/",
	"foo/testfile":         "foo/",
}

func TestExtract(t *testing.T) {
	cases := map[string]struct { //nolint:govet
		giveArchives []string
		giveOptions  []Option
		wantTree     map[string]string
		wantMissing  []string
		wantErr      error
	}{
		"gzip": {
			giveArchives: []string{
				"testdata/foo.tar.gz",
				"testdata/foo.tgz",
			},
			giveOptions: nil,
			wantTree:    _fullTree,
			wantMissing: nil,
			wantErr:     nil,
		},
		"bzip2": {
			giveArchives: []string{
				"testdata/foo.tar.bz",
				"testdata/foo.tar.bz2",
				"testdata/foo.tbz",
				"testdata/foo.tbz2",
			},
			giveOptions: nil,
			wantTree:    _fullTree,
			wantMissing: nil,
			wantErr:     nil,
		},
		"xz": {
			giveArchives: []string{
				"testdata/foo.tar.xz",
			},
			giveOptions: nil,
			wantTree:    _fullTree,
			wantMissing: nil,
			wantErr:     nil,
		},
		"zip": {
			giveArchives: []string{
				"testdata/foo.zip",
			},
			giveOptions: nil,
			wantTree:    _fullTree,
			wantMissing: nil,
			wantErr:     nil,
		},
		"dne": {
			giveArchives: []string{
				"testdata/does-not-exist.tgz",
			},
			giveOptions: nil,
			wantTree:    nil,
			wantMissing: nil,
			wantErr:     os.ErrNotExist,
		},
		"include paths": {
			giveArchives: []string{
				"testdata/foo.tar.bz",
				"testdata/foo.tar.bz2",
				"testdata/foo.tar.gz",
				"testdata/foo.tar.xz",
				"testdata/foo.tbz",
				"testdata/foo.tbz2",
				"testdata/foo.tgz",
				"testdata/foo.zip",
			},
			giveOptions: []Option{
				IncludePaths(map[string]string{
					"foo/*ar/bat/testfile": "overridden-bat-dst",
					"foo/baz/*/testfile":   "",
					"foo/testfile":         "",
				}),
			},
			wantTree: map[string]string{
				"overridden-bat-dst":   _fullTree["foo/bar/bat/testfile"],
				"foo/baz/bat/testfile": _fullTree["foo/baz/bat/testfile"],
				"foo/baz/qux/testfile": _fullTree["foo/baz/qux/testfile"],
				"foo/testfile":         _fullTree["foo/testfile"],
			},
			wantMissing: []string{
				"foo/bar/bat/testfile",
				"foo/bar/qux/testfile",
				"foo/bar/testfile",
				"foo/baz/testfile",
			},
			wantErr: nil,
		},
		"exclude paths": {
			giveArchives: []string{
				"testdata/foo.tar.bz",
				"testdata/foo.tar.bz2",
				"testdata/foo.tar.gz",
				"testdata/foo.tar.xz",
				"testdata/foo.tbz",
				"testdata/foo.tbz2",
				"testdata/foo.tgz",
				"testdata/foo.zip",
			},
			giveOptions: []Option{
				ExcludePaths([]string{
					"foo/bar/*",
					"foo/bar/*/*",
					"foo/baz/*",
				}),
			},
			wantTree: map[string]string{
				"foo/testfile":         _fullTree["foo/testfile"],
				"foo/baz/bat/testfile": _fullTree["foo/baz/bat/testfile"],
				"foo/baz/qux/testfile": _fullTree["foo/baz/qux/testfile"],
			},
			wantMissing: []string{
				"foo/bar/bat/testfile",
				"foo/bar/qux/testfile",
				"foo/bar/testfile",
				"foo/baz/testfile",
			},
			wantErr: nil,
		},
		"strip prefix": {
			giveArchives: []string{
				"testdata/foo.tar.bz",
				"testdata/foo.tar.bz2",
				"testdata/foo.tar.gz",
				"testdata/foo.tar.xz",
				"testdata/foo.tbz",
				"testdata/foo.tbz2",
				"testdata/foo.tgz",
				"testdata/foo.zip",
			},
			giveOptions: []Option{
				StripPrefix("foo"),
				IncludePaths(map[string]string{
					"testfile": "",
				}),
			},
			wantTree: map[string]string{
				"testfile": _fullTree["foo/testfile"],
			},
			wantMissing: []string{
				"foo/bar/bat/testfile",
				"foo/bar/bat/testfile",
				"foo/bar/qux/testfile",
				"foo/bar/testfile",
				"foo/baz/bat/testfile",
				"foo/baz/qux/testfile",
				"foo/baz/testfile",
			},
			wantErr: nil,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			for _, archive := range tt.giveArchives {
				t.Run(filepath.Base(archive), func(t *testing.T) {
					var (
						dst  = t.TempDir()
						opts = Options{}.With(tt.giveOptions...)
					)

					err := Extract(context.Background(), dst, archive, opts)
					require.ErrorIs(t, err, tt.wantErr)

					for relpath, contents := range tt.wantTree {
						path := filepath.Join(dst, relpath)

						stat, statErr := os.Stat(path)
						require.NoError(t, statErr)
						require.False(t, stat.IsDir())

						raw, readErr := os.ReadFile(path)
						require.NoError(t, readErr)
						require.Equal(t, contents, string(bytes.TrimSpace(raw)))
					}

					for _, relpath := range tt.wantMissing {
						_, statErr := os.Stat(filepath.Join(dst, relpath))
						require.ErrorIs(t, statErr, os.ErrNotExist, relpath)
					}
				})
			}
		})
	}
}

func TestExtract_Errors(t *testing.T) {
	t.Run("strip prefix", func(t *testing.T) {
		var (
			wantErr          = errors.New(t.Name())
			badFilepathMatch = func(string, string) (bool, error) {
				return false, wantErr
			}
		)
		stub.With(&_filepathMatch, badFilepathMatch, func() {
			const archive = "testdata/foo.tgz"
			err := Extract(
				context.Background(),
				t.TempDir(),
				archive,
				StripPrefix("testdata"),
			)
			require.ErrorIs(t, err, wantErr)
			require.ErrorContains(t, err, "failed to strip prefix")
		})
	})

	t.Run("exclude path", func(t *testing.T) {
		var (
			wantErr          = errors.New(t.Name())
			badFilepathMatch = func(string, string) (bool, error) {
				return false, wantErr
			}
		)
		stub.With(&_filepathMatch, badFilepathMatch, func() {
			const archive = "testdata/foo.tgz"
			err := Extract(
				context.Background(),
				t.TempDir(),
				archive,
				ExcludePaths([]string{"testdata"}),
			)
			require.ErrorIs(t, err, wantErr)
			require.ErrorContains(t, err, "bad exclude path pattern")
		})
	})

	t.Run("include path", func(t *testing.T) {
		var (
			wantErr          = errors.New(t.Name())
			badFilepathMatch = func(string, string) (bool, error) {
				return false, wantErr
			}
		)
		stub.With(&_filepathMatch, badFilepathMatch, func() {
			const archive = "testdata/foo.tgz"
			err := Extract(
				context.Background(),
				t.TempDir(),
				archive,
				IncludePaths(map[string]string{"testdata": "testdata"}),
			)
			require.ErrorIs(t, err, wantErr)
			require.ErrorContains(t, err, "bad include path pattern")
		})
	})

	t.Run("parent removeall", func(t *testing.T) {
		var (
			wantErr      = errors.New(t.Name())
			badRemoveAll = func(string) error {
				return wantErr
			}
		)
		stub.With(&_osRemoveAll, badRemoveAll, func() {
			const archive = "testdata/foo.tgz"
			err := Extract(
				context.Background(),
				t.TempDir(),
				archive,
				Delete(true),
			)
			require.ErrorIs(t, err, wantErr)
			require.ErrorContains(
				t,
				err,
				"failed to remove existing destination directory",
			)
		})
	})

	t.Run("parent mkdirallinherit", func(t *testing.T) {
		var (
			wantErr            = errors.New(t.Name())
			badMkdirAllInherit = func(string) error {
				return wantErr
			}
		)
		stub.With(&_xosMkdirAllInherit, badMkdirAllInherit, func() {
			const archive = "testdata/foo.tgz"
			err := Extract(
				context.Background(),
				t.TempDir(),
				archive,
				Delete(true),
			)
			require.ErrorIs(t, err, wantErr)
			require.ErrorContains(
				t,
				err,
				"failed to create destination parent(s)",
			)
		})
	})

	t.Run("write to file", func(t *testing.T) {
		var (
			wantErr                       = errors.New(t.Name())
			badWriteReaderToFileWithFlags = func(
				_ string,
				_ any,
				_ int,
				_ fs.FileMode,
			) (int, error) {
				return 0, wantErr
			}
		)
		stub.With(
			&_xosWriteReaderToFileWithFlags,
			badWriteReaderToFileWithFlags,
			func() {
				const archive = "testdata/foo.tgz"
				err := Extract(
					context.Background(),
					t.TempDir(),
					archive,
					Delete(true),
				)
				require.ErrorIs(t, err, wantErr)
				require.ErrorContains(t, err, "failed to extract")
			},
		)
	})
}

func TestExtractOutput(t *testing.T) {
	archives := []string{
		"testdata/foo.tar.bz",
		"testdata/foo.tar.bz2",
		"testdata/foo.tar.gz",
		"testdata/foo.tar.xz",
		"testdata/foo.tbz",
		"testdata/foo.tbz2",
		"testdata/foo.tgz",
		"testdata/foo.zip",
	}

	for _, archive := range archives {
		t.Run(filepath.Base(archive), func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			require.NoError(
				t,
				Extract(
					context.Background(),
					ExtractToTempDir,
					archive,
					Output(buf),
				),
			)

			var (
				lines = strings.Split(strings.TrimSpace(buf.String()), "\n")
				count int
			)

			// Ignore ._* files.
			for _, line := range lines {
				if !strings.Contains(line, "/._") {
					count++
				}
			}

			require.Equal(t, len(_fullTree), count)
		})
	}
}

func TestExtractTempDirWithCallback(t *testing.T) {
	archives := []string{
		"testdata/foo.tar.bz",
		"testdata/foo.tar.bz2",
		"testdata/foo.tar.gz",
		"testdata/foo.tar.xz",
		"testdata/foo.tbz",
		"testdata/foo.tbz2",
		"testdata/foo.tgz",
		"testdata/foo.zip",
	}

	for _, archive := range archives {
		t.Run(filepath.Base(archive), func(t *testing.T) {
			var (
				wantErr = errors.New("done")
				err     = Extract(
					context.Background(),
					ExtractToTempDir,
					archive,
					Callback(func(_ context.Context, dir string) error {
						err := filepath.WalkDir(
							dir,
							func(path string, _ fs.DirEntry, err error) error {
								require.NoError(t, err)

								contents, exists := _fullTree[path]
								if !exists {
									return nil
								}

								raw, err := os.ReadFile(path)
								require.NoError(t, err)
								require.Equal(t, contents, string(raw))
								return nil
							},
						)
						require.NoError(t, err)
						return wantErr
					}),
				)
			)
			require.ErrorIs(t, err, wantErr)
		})
	}
}

func TestExtract_MkdirTemp(t *testing.T) {
	var (
		wantErr      = errors.New(t.Name())
		badMkdirTemp = func(_ string, _ string) (string, error) {
			return "", wantErr
		}
	)

	stub.With(&_osMkdirTemp, badMkdirTemp, func() {
		err := Extract(
			context.Background(),
			ExtractToTempDir,
			"foo.tgz",
		)
		require.ErrorIs(t, err, wantErr)
	})
}

func TestStripPrefix(t *testing.T) {
	cases := map[string]struct { //nolint:govet
		givePath     string
		givePrefix   string
		wantPath     string
		wantStripped bool
		wantErr      error
	}{
		"matched exact": {
			givePath:     "foo/bar/baz/bat",
			givePrefix:   "foo/bar",
			wantPath:     "baz/bat",
			wantStripped: true,
			wantErr:      nil,
		},
		"matched exact with trailing slash": {
			givePath:     "foo/bar/baz/bat",
			givePrefix:   "foo/bar/",
			wantPath:     "baz/bat",
			wantStripped: true,
			wantErr:      nil,
		},
		"matched glob": {
			givePath:     "foo/bar/baz/bat",
			givePrefix:   "f*/*",
			wantPath:     "baz/bat",
			wantStripped: true,
			wantErr:      nil,
		},
		"not matched exact": {
			givePath:     "foo/bar",
			givePrefix:   "bar",
			wantPath:     "foo/bar",
			wantStripped: false,
			wantErr:      nil,
		},
		"not matched glob": {
			givePath:     "foo/bar",
			givePrefix:   "o*",
			wantPath:     "foo/bar",
			wantStripped: false,
			wantErr:      nil,
		},
		"empty prefix": {
			givePath:     "foo/bar",
			givePrefix:   "",
			wantPath:     "foo/bar",
			wantStripped: true,
			wantErr:      nil,
		},
		"empty strings": {
			givePath:     "",
			givePrefix:   "",
			wantPath:     "",
			wantStripped: true,
			wantErr:      nil,
		},
		"bad prefix": {
			givePath:     "foo/bar",
			givePrefix:   "\\",
			wantPath:     "",
			wantStripped: false,
			wantErr:      filepath.ErrBadPattern,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			path, stripped, err := stripPrefix(tt.givePath, tt.givePrefix)
			require.Equal(t, tt.wantPath, path)
			require.Equal(t, tt.wantStripped, stripped)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestStripPrefix_Error(t *testing.T) {
	var (
		wantErr          = errors.New(t.Name())
		badFilepathMatch = func(string, string) (bool, error) {
			return false, wantErr
		}
	)
	stub.With(&_filepathMatch, badFilepathMatch, func() {
		path, stripped, err := stripPrefix("foo/bar", "foo")
		require.ErrorIs(t, err, wantErr)
		require.False(t, stripped)
		require.Empty(t, path)
	})
}
