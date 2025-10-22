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

// Package tempdir provides helpers for temporary directories.
package tempdir

import (
	"os"
	"runtime/debug"
	"strings"
	"sync"

	"go.mway.dev/errors"

	xos "go.mway.dev/x/os"
)

var (
	// ErrInvalidFuncType indicates that an invalid type was provided.
	ErrInvalidFuncType = errors.New("invalid func type")

	_osMkdirTemp = os.MkdirTemp
	_osRemoveAll = os.RemoveAll
)

type (
	// A PathFunc is a function that accepts a temporary path.
	PathFunc = func(path string)
	// A PathErrorFunc is a function that accepts a temporary path and returns
	// an error.
	PathErrorFunc = func(path string) error
)

// Func constrains the types of functions allowed by [With] (and opaqely,
// [Dir.With]).
type Func interface {
	PathFunc | PathErrorFunc
}

// With creates a new [Dir], and calls [Dir.With](fn), returning any errors.
func With[T Func](fn T) (err error) {
	var d *Dir
	if d, err = New(); err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, d.Close())
	}()

	return d.With(fn)
}

// A Dir represents a temporary directory. Dir must be constructed with [New]
// or [NewWithBase].
type Dir struct {
	path string
	done bool
	mu   sync.Mutex
}

// New creates a new [Dir] using the system's default temporary directory as
// the base for this temporary directory.
func New() (*Dir, error) {
	return NewWithBase(os.TempDir())
}

// NewWithBase creates a new [Dir] using the provided dir as the base for this
// temporary directory.
func NewWithBase(dir string) (*Dir, error) {
	pattern := "x-tempdir"
	build, ok := debug.ReadBuildInfo()
	if ok {
		pattern = strings.ToLower(strings.ReplaceAll(
			build.Path,
			string(os.PathSeparator),
			"_",
		))
	}

	var err error
	if dir, err = _osMkdirTemp(dir, pattern); err != nil {
		return nil, errors.Wrapf(
			err,
			"failed to MkdirTemp(%q, %q)",
			dir,
			pattern,
		)
	}

	return &Dir{
		path: dir,
	}, nil
}

// With changes the working directory to d's temporary directory and invokes
// fn, returning any error.
func (d *Dir) With(fn any) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.done {
		return nil
	}

	var wrapper func() error
	switch f := fn.(type) {
	case PathFunc:
		wrapper = func() error {
			f(d.path)
			return nil
		}
	case PathErrorFunc:
		wrapper = func() error {
			return f(d.path)
		}
	default:
		return errors.Wrapf(ErrInvalidFuncType, "%T", fn)
	}

	return xos.WithCwd(d.path, wrapper)
}

// Close cleans up d's temporary directory. No further calls to d.With will be
// executed.
func (d *Dir) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	var done bool
	if done, d.done = d.done, true; done {
		return nil
	}

	return errors.Wrapf(
		_osRemoveAll(d.path),
		"failed to remove tempdir %q",
		d.path,
	)
}
