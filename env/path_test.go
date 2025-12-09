// Copyright (c) 2025 Matt Way
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
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mway.dev/x/stub"
)

func TestPath(t *testing.T) {
	withDummyPath(t.Name(), func() {
		path := NewPath()
		require.Equal(t, t.Name(), path.Value())

		path, err := path.Prepend("a", "b", "c")
		require.NoError(t, err)
		require.Equal(t, "a:b:c:"+t.Name(), path.String())

		path, err = path.Append("d", "e", "f")
		require.NoError(t, err)
		require.Equal(t, "a:b:c:"+t.Name()+":d:e:f", path.String())
	})
}

func TestPath_Append_Empty(t *testing.T) {
	withDummyPath(t.Name(), func() {
		path := NewPath().MustAppend().MustAppend("", "foo")
		require.Equal(t, t.Name()+":foo", path.String())
	})
}

func TestPath_Prepend_Error(t *testing.T) {
	withDummyPath(t.Name(), func() {
		wantErr := errors.New("setenv error")
		stub.With(&_osSetenv, osSetenvReturning(wantErr), func() {
			path := NewPath()
			_, err := path.Prepend("foo")
			require.ErrorIs(t, err, wantErr)
			require.Panics(t, func() {
				path.MustPrepend("foo")
			})
		})
	})
}

func TestPath_Append_Error(t *testing.T) {
	withDummyPath(t.Name(), func() {
		wantErr := errors.New("setenv error")
		stub.With(&_osSetenv, osSetenvReturning(wantErr), func() {
			path := NewPath()
			_, err := path.Append("foo")
			require.ErrorIs(t, err, wantErr)
			require.Panics(t, func() {
				path.MustAppend("foo")
			})
		})
	})
}

func TestPath_Prepend_Empty(t *testing.T) {
	withDummyPath(t.Name(), func() {
		path := NewPath().MustPrepend().MustPrepend("", "foo")
		require.Equal(t, "foo:"+t.Name(), path.String())
	})
}

func TestPath_Restore_Error(t *testing.T) {
	withDummyPath(t.Name(), func() {
		wantErr := errors.New("setenv error")
		stub.With(&_osSetenv, osSetenvReturning(wantErr), func() {
			path := NewPath()
			require.ErrorIs(t, path.Restore(), wantErr)
			require.Panics(t, path.MustRestore)
		})
	})
}
