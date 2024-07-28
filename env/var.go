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
	"os"
)

var (
	_osLookupEnv = os.LookupEnv
	_osSetenv    = os.Setenv
	_osGetenv    = os.Getenv
	_osUnsetenv  = os.Unsetenv
)

// Var is a helper for managing environment variables.
type Var struct {
	key    string
	value  string
	orig   string
	exists bool
}

// NewVar creates a new [Var] that manages an environment variable with the
// given key.
func NewVar(key string) *Var {
	orig, exists := _osLookupEnv(key)
	return &Var{
		key:    key,
		value:  orig,
		orig:   orig,
		exists: exists,
	}
}

// NewVarWithValue creates a new [Var] that manages an environment variable
// with the given key and sets its value to the given value.
func NewVarWithValue(key string, value string) (*Var, error) {
	orig, exists := _osLookupEnv(key)
	v := &Var{
		key:    key,
		value:  orig,
		orig:   orig,
		exists: exists,
	}
	if err := v.Set(value); err != nil {
		return nil, err
	}

	return v, nil
}

// MustVar returns v, panicking if err is not nil.
func MustVar(v *Var, err error) *Var {
	if err != nil {
		panic(err)
	}
	return v
}

// Key returns the key (environment variable name) for v.
func (v *Var) Key() string {
	return v.key
}

// Value returns the current value of v.
func (v *Var) Value() string {
	return v.value
}

// Clone clones v.
func (v *Var) Clone() *Var {
	tmp := *v
	return &tmp
}

// Set uses value as the current value of both v and its environment variable.
func (v *Var) Set(value string) error {
	if err := _osSetenv(v.key, value); err != nil {
		return err
	}
	v.value = value
	return nil
}

// MustSet calls v.Set and panics if it returns an error.
func (v *Var) MustSet(value string) {
	if err := v.Set(value); err != nil {
		panic(err)
	}
}

// Unset unsets v within the environment, if set.
func (v *Var) Unset() error {
	if err := _osUnsetenv(v.key); err != nil {
		return err
	}

	v.value = ""
	return nil
}

// MustUnset calls v.Unset and panics if it returns an error.
func (v *Var) MustUnset() {
	if err := v.Unset(); err != nil {
		panic(err)
	}
}

// Load loads the current environment value of v and stores it.
func (v *Var) Load() {
	v.value = _osGetenv(v.key)
}

// Restore restores v to its original state. If the variable managed by v did
// not exist previously, it is unset; otherwise, the original value is set in
// both v and the environment.
func (v *Var) Restore() error {
	if v.exists {
		return v.Set(v.orig)
	}

	return v.Unset()
}

// MustRestore calls v.Restore and panics if it returns an error.
func (v *Var) MustRestore() {
	if err := v.Restore(); err != nil {
		panic(err)
	}
}
