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
	"io/fs"
	"os"
	gofilepath "path/filepath"

	"go.mway.dev/errors"
	"go.mway.dev/x/path/filepath"
)

var (
	_osStat     = os.Stat
	_osMkdirAll = os.MkdirAll
)

// MkdirAllInherit calls [os.MkdirAll] on each ancestor in path, with each
// descendant inheriting the mode of its parent.
func MkdirAllInherit(path string) error {
	paths := append([]string{path}, filepath.Ancestors(path)...)
	for i := len(paths) - 1; i >= 0; i-- {
		mode, err := getDirPerm(paths[i])
		if err != nil {
			return errors.Wrapf(
				err,
				"failed to inherit permissions for %q",
				paths[i],
			)
		}

		if err = _osMkdirAll(paths[i], mode); err != nil {
			return errors.Wrapf(
				err,
				"failed to create directory %q with mode %s",
				paths[i],
				mode,
			)
		}
	}

	return nil
}

func getDirPerm(path string) (fs.FileMode, error) {
	parent := gofilepath.Dir(path)
	stat, err := _osStat(parent)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get mode for %q", parent)
	}
	return stat.Mode() & fs.ModePerm, nil
}
