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

// With executes fn with the given name=value set in the environment for the
// duration of the call and then restores the environment. If an error is
// encountered when setting the environment variable, it will be ignored.
func With(fn func(), name string, value string) {
	//nolint:errcheck
	WithError(func() error {
		fn()
		return nil
	}, name, value)
}

// WithError executes fn with the given name=value set in the environment for
// the duration of the call, returns the result, and restores the environment.
// If an error is encountered when setting the environment variable, it will be
// returned prior to invoking fn.
func WithError(fn func() error, name string, value string) error {
	v, err := NewVarWithValue(name, value)
	if err != nil {
		return err
	}
	defer v.Restore() //nolint:errcheck
	return fn()
}

// Without executes fn with any named variable(s) removed from the environment
// for the duration of the call. If an error is encountered when unsetting an
// environment variable, it will be ignored.
func Without(fn func(), names ...string) {
	//nolint:errcheck
	WithoutError(func() error {
		fn()
		return nil
	}, names...)
}

// WithoutError executes fn with any named variable(s) removed from the
// environment for the duration of the call, and returns the result. If an
// error is encountered when unsetting the environment variable, it will be
// returned prior to invoking fn.
func WithoutError(fn func() error, names ...string) error {
	for _, name := range names {
		v := NewVar(name)
		if err := v.Unset(); err != nil {
			return err
		}
		defer v.Restore() //nolint:errcheck,gocritic
	}
	return fn()
}

// WithAll sets any name=value pairs in the environment, executes fn, and
// restores the environment.
func WithAll(fn func(), vars map[string]string) {
	//nolint:errcheck
	WithAllError(func() error {
		fn()
		return nil
	}, vars)
}

// WithAllError sets any name=value pairs in the environment, executes fn and
// returns the result, and restores the environment.
func WithAllError(fn func() error, vars map[string]string) error {
	for name, value := range vars {
		v, err := NewVarWithValue(name, value)
		if err != nil {
			return err
		}
		defer v.Restore() //nolint:errcheck,gocritic
	}
	return fn()
}
