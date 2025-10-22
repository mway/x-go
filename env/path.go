// Copyright (c) 2023 Matt Way
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

package env

import (
	"bytes"
	"fmt"
	"os"
)

const _envKeyPath = "PATH"

// Path represents $PATH within the environment.
type Path struct {
	v *Var
}

// NewPath returns a [Path] that holds the current value of $PATH.
func NewPath() Path {
	return Path{
		v: NewVar(_envKeyPath),
	}
}

// Prepend prepends paths to the current value of p.
func (p Path) Prepend(paths ...string) (newpath Path, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to prepend path(s): %w", err)
		}
	}()

	if len(paths) == 0 {
		return p, nil
	}

	// Calculate how much extra storage is needed.
	extra := len(paths) // for delimiters
	for _, path := range paths {
		extra += len(path)
	}

	// Prepend all of the path fragments to the current value.
	var buf bytes.Buffer
	buf.Grow(len(p.v.Value()) + extra)
	for _, path := range paths {
		if len(path) == 0 {
			continue
		}

		buf.WriteString(path)
		buf.WriteByte(os.PathListSeparator)
	}

	// Write the current value after any prepended path(s).
	buf.WriteString(p.Value())

	// Create a copy of this Path and persist the change to the copy (and the
	// environment).
	tmp := Path{
		v: p.v.Clone(),
	}
	if setErr := tmp.v.Set(buf.String()); setErr != nil {
		return newpath, setErr
	}
	return tmp, nil
}

// MustPrepend calls p.Prepend and panics if it returns an error.
func (p Path) MustPrepend(paths ...string) Path {
	newpath, err := p.Prepend(paths...)
	if err != nil {
		panic(err)
	}
	return newpath
}

// Append appends paths to the current value of p.
func (p Path) Append(paths ...string) (newpath Path, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to append path(s): %w", err)
		}
	}()

	if len(paths) == 0 {
		return p, nil
	}

	// Calculate how much extra storage is needed.
	extra := len(paths) // for delimiters
	for _, path := range paths {
		extra += len(path)
	}

	var buf bytes.Buffer
	buf.Grow(len(p.v.Value()) + extra)

	// Write the current value before any appended path(s).
	buf.WriteString(p.v.Value())

	// Append all of the path fragments to the current value.
	for _, path := range paths {
		if len(path) == 0 {
			continue
		}

		buf.WriteByte(os.PathListSeparator)
		buf.WriteString(path)
	}

	// Create a copy of this Path and persist the change to the copy (and the
	// environment).
	tmp := Path{
		v: p.v.Clone(),
	}
	if setErr := tmp.v.Set(buf.String()); setErr != nil {
		return newpath, setErr
	}
	return tmp, nil
}

// MustAppend calls p.Append and panics if it returns an error.
func (p Path) MustAppend(paths ...string) Path {
	newpath, err := p.Append(paths...)
	if err != nil {
		panic(err)
	}
	return newpath
}

// Restore PATH to its original state within the environment.
func (p Path) Restore() error {
	return p.v.Restore()
}

// MustRestore calls p.Restore and panics if it returns an error.
func (p Path) MustRestore() {
	if err := p.Restore(); err != nil {
		panic(err)
	}
}

// Value returns the current value of p.
func (p Path) Value() string {
	return p.v.Value()
}

// String returns p as a string.
func (p Path) String() string {
	return p.Value()
}
