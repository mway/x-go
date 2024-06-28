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
	"errors"
	"os"
	"testing"
	"unsafe"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mway.dev/x/stub"
)

func TestVar_NotExists(t *testing.T) {
	key := uuid.New().String()

	v := NewVar(key)
	require.Equal(t, key, v.Key())
	require.Empty(t, v.Value())
	require.False(t, v.exists)

	envValue, exists := os.LookupEnv(v.Key())
	require.False(t, exists)
	require.Empty(t, envValue)

	require.NoError(t, v.Set(t.Name()))
	envValue, exists = os.LookupEnv(v.Key())
	require.True(t, exists)
	require.Equal(t, t.Name(), envValue)
	require.Equal(t, t.Name(), v.Value())

	require.NoError(t, v.Restore())
	envValue, exists = os.LookupEnv(v.Key())
	require.False(t, exists)
	require.Empty(t, envValue)
	require.Empty(t, v.Value())
}

func TestVar_Exists(t *testing.T) {
	key := uuid.New().String()

	// Ensure there is a value to start.
	require.NoError(t, os.Setenv(key, t.Name()))
	defer func() {
		require.NoError(t, os.Unsetenv(key))
	}()

	v := NewVar(key)
	require.Equal(t, key, v.Key())
	require.Equal(t, t.Name(), v.Value())
	require.True(t, v.exists)

	envValue, exists := os.LookupEnv(v.Key())
	require.True(t, exists)
	require.Equal(t, t.Name(), envValue)

	require.NoError(t, v.Set("override"))
	envValue, exists = os.LookupEnv(v.Key())
	require.True(t, exists)
	require.Equal(t, "override", envValue)
	require.Equal(t, "override", v.Value())

	require.NoError(t, v.Restore())
	envValue, exists = os.LookupEnv(v.Key())
	require.True(t, exists)
	require.Equal(t, t.Name(), envValue)
	require.Equal(t, t.Name(), v.Value())
}

func TestVar_Load(t *testing.T) {
	key := uuid.New().String()

	// Load a previously unset env.
	v := NewVar(key)
	defer v.MustUnset()

	require.NoError(t, os.Setenv(key, t.Name()))
	v.Load()
	require.Equal(t, t.Name(), v.Value())

	// Load a previously set env.
	v = NewVar(key)
	require.Equal(t, t.Name(), v.Value())
	require.NoError(t, os.Setenv(key, "override"))
	v.Load()
	require.Equal(t, "override", v.Value())
}

func TestVar_Clone(t *testing.T) {
	var (
		v1 = NewVar("foo")
		v2 = v1.Clone()
	)

	require.Equal(t, unsafe.Pointer(v1), unsafe.Pointer(v1))
	require.NotEqual(t, unsafe.Pointer(v1), unsafe.Pointer(v2))
}

func TestVar_Set_Error(t *testing.T) {
	err := errors.New("setenv error")
	stub.With(&_osSetenv, osSetenvReturning(err), func() {
		v := NewVar("test")
		require.ErrorIs(t, v.Set("this will error"), err)
		require.Panics(t, func() {
			v.MustSet("this will panic")
		})
	})
}

func TestVar_Restore_Error(t *testing.T) {
	err := errors.New("setenv/unsetenv error")

	// Does not exist; uses os.Unsetenv
	stub.With(&_osUnsetenv, osUnsetenvReturning(err), func() {
		v := NewVar("test")
		require.ErrorIs(t, v.Restore(), err)
		require.Panics(t, v.MustRestore)
	})

	// Exists; uses os.Setenv
	stub.With(&_osSetenv, osSetenvReturning(err), func() {
		require.NoError(t, os.Setenv("test", "test"))
		v := NewVar("test")
		require.ErrorIs(t, v.Restore(), err)
		require.Panics(t, v.MustRestore)
	})
}

func TestVar_Unset_Error(t *testing.T) {
	err := errors.New("setenv error")
	stub.With(&_osUnsetenv, osUnsetenvReturning(err), func() {
		v := NewVar("test")
		require.ErrorIs(t, v.Unset(), err)
		require.Panics(t, v.MustUnset)
	})
}
