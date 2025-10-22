// Copyright (c) 2022 Matt Way
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
	"errors"
	"os"
)

// n.b. See with_cwd_internals_test.go.
var (
	_getwd = os.Getwd
	_chdir = os.Chdir
)

// CwdFuncConstraint constrains the function types allowed for [WithCwd].
type CwdFuncConstraint interface {
	func() | func() error
}

// WithCwd attempts to change the working directory to dir and, if successful,
// calls f. WithCwd expects dir to exist already; if dir does not exist, or is
// removed or renamed during f's execution, an error will be returned.
//
// WithCwd will attempt to restore the original working directory, even if the
// given function panics.
func WithCwd[T CwdFuncConstraint](dir string, fn T) (err error) {
	var orig string
	if orig, err = _getwd(); err != nil {
		return err
	}

	// If we're not already in the target directory, change directories.
	if dir != orig {
		if chdirErr := _chdir(dir); chdirErr != nil {
			return chdirErr
		}

		defer func() {
			err = errors.Join(err, os.Chdir(orig))
		}()
	}

	switch fn := any(fn).(type) {
	case func():
		fn()
	case func() error:
		err = fn()
	}

	return err
}
