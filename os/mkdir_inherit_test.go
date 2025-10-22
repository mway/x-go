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
	goos "os"
	gofilepath "path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"go.mway.dev/x/os"
)

var (
	_modeTree = map[string]goos.FileMode{
		"foo":                 0o755,
		"foo/bar":             0o765,
		"foo/bar/baz":         0o766,
		"foo/bar/baz/bat":     0o777,
		"foo/bar/baz/bat/qux": 0o777,
	}
	_modePaths = func() []string {
		paths := maps.Keys(_modeTree)
		sort.Slice(paths, func(i int, j int) bool {
			return len(paths[i]) < len(paths[j])
		})
		return paths
	}()
)

func TestMkdirAllInherit(t *testing.T) {
	tempdir := createModeTree(t)

	for _, path := range _modePaths {
		t.Run(path, func(t *testing.T) {
			path = gofilepath.Join(tempdir, path)

			newdir := gofilepath.Join(path, "newdir")
			require.NoError(t, os.MkdirAllInherit(newdir))
			requireModePermsMatch(t, path, newdir)
		})
	}
}

func TestMkdirAllInheritRelative(t *testing.T) {
	tempdir := createModeTree(t)

	for _, path := range _modePaths {
		t.Run(path, func(t *testing.T) {
			err := os.WithCwd(tempdir, func() {
				newdir := gofilepath.Join(path, "newdir")
				require.NoError(t, os.MkdirAllInherit(newdir))
				requireModePermsMatch(t, path, newdir)
			})
			require.NoError(t, err)
		})
	}
}

func createModeTree(t *testing.T) string {
	tempdir := t.TempDir()

	for _, path := range _modePaths {
		fqpath := gofilepath.Join(tempdir, path)
		require.NoError(t, goos.MkdirAll(fqpath, _modeTree[path]))
	}

	return tempdir
}

func requireModePermsMatch(t *testing.T, a string, b string) {
	t.Helper()

	astat, aerr := goos.Stat(a)
	require.NoErrorf(t, aerr, "failed to stat %q", a)

	bstat, berr := goos.Stat(b)
	require.NoErrorf(t, berr, "failed to stat %q", b)

	require.Equal(t, astat.Mode()&goos.ModePerm, bstat.Mode()&goos.ModePerm)
}
