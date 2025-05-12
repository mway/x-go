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

package env_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/env"
)

type (
	OptionsFunc = func(*testing.T) []env.LookupOption
	SetupFunc   = func(*testing.T) CleanupFunc
	CleanupFunc = func()
)

var (
	_noOptions = func(*testing.T) []env.LookupOption { return nil }
	_noSetup   = func(*testing.T) CleanupFunc { return func() {} }
)

func TestGet(t *testing.T) {
	cases := map[string]struct {
		setup SetupFunc
		opts  OptionsFunc
	}{
		"os.LookupEnv": {
			setup: func(t *testing.T) CleanupFunc {
				return vars{
					env.MustVar(
						env.NewVarWithValue("TEST_VAR_1", t.Name()+"1"),
					),
					env.MustVar(
						env.NewVarWithValue("TEST_VAR_2", t.Name()+"2"),
					),
					env.MustVar(
						env.NewVarWithValue("TEST_VAR_3", t.Name()+"3"),
					),
				}.MustRestore
			},
			opts: func(*testing.T) []env.LookupOption {
				return []env.LookupOption{
					env.SanitizeNames(true),
				}
			},
		},
		"LookupFunc": {
			setup: _noSetup,
			opts: func(t *testing.T) []env.LookupOption {
				return []env.LookupOption{
					env.SanitizeNames(true),
					env.LookupFunc(func(key string) (string, bool) {
						switch key {
						case "TEST_VAR_1":
							return t.Name() + "1", true
						case "TEST_VAR_2":
							return t.Name() + "2", true
						case "TEST_VAR_3":
							return t.Name() + "3", true
						default:
							return "", false
						}
					}),
				}
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			defer tt.setup(t)()
			opts := tt.opts(t)
			require.Equal(t, t.Name()+"1", env.Get("TEST_VAR_1", opts...))
			require.Equal(t, t.Name()+"2", env.Get("TEST_VAR_2", opts...))
			require.Equal(t, t.Name()+"3", env.Get("TEST_VAR_3", opts...))
		})
	}
}

func TestGet_NotFound(t *testing.T) {
	name := env.SanitizeName(t.Name())
	_, found := os.LookupEnv(name)
	require.False(t, found, name)
	require.Zero(t, env.Get(name))
}

func TestLookup(t *testing.T) {
	cases := map[string]struct {
		setup SetupFunc
		opts  OptionsFunc
	}{
		"os.LookupEnv": {
			setup: func(t *testing.T) CleanupFunc {
				return vars{
					env.MustVar(
						env.NewVarWithValue("TEST_VAR_1", t.Name()+"1"),
					),
					env.MustVar(
						env.NewVarWithValue("TEST_VAR_2", t.Name()+"2"),
					),
					env.MustVar(
						env.NewVarWithValue("TEST_VAR_3", t.Name()+"3"),
					),
				}.MustRestore
			},
			opts: _noOptions,
		},
		"LookupFunc": {
			setup: _noSetup,
			opts: func(t *testing.T) []env.LookupOption {
				return []env.LookupOption{
					env.LookupFunc(func(key string) (string, bool) {
						switch key {
						case "TEST_VAR_1":
							return t.Name() + "1", true
						case "TEST_VAR_2":
							return t.Name() + "2", true
						case "TEST_VAR_3":
							return t.Name() + "3", true
						default:
							return "", false
						}
					}),
				}
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			defer tt.setup(t)()

			opts := tt.opts(t)
			requireLookupFound(t, t.Name()+"1", "TEST_VAR_1", opts...)
			requireLookupPFound(t, t.Name()+"1", "TEST_VAR_1", opts...)
			requireLookupFound(t, t.Name()+"2", "TEST_VAR_2", opts...)
			requireLookupPFound(t, t.Name()+"2", "TEST_VAR_2", opts...)
			requireLookupFound(t, t.Name()+"3", "TEST_VAR_3", opts...)
			requireLookupPFound(t, t.Name()+"3", "TEST_VAR_3", opts...)
		})
	}
}

func TestLookup_NotFound(t *testing.T) {
	name := env.SanitizeName(t.Name())
	_, found := os.LookupEnv(name)
	require.False(t, found, name)

	var (
		val, valfound = env.Lookup(name)
		ptr, ptrfound = env.LookupP(name)
	)

	require.False(t, valfound)
	require.Zero(t, val)

	require.False(t, ptrfound)
	require.Nil(t, ptr)
}

func TestGetLookupAs_Integers(t *testing.T) {
	v := env.NewVar("TEST_VAR")
	defer v.MustRestore()
	require.NoError(t, v.Set("123"))

	requireGetAsFound(t, int(123), "TEST_VAR")
	requireLookupAsFound(t, int(123), "TEST_VAR")
	requireLookupAsPFound(t, int(123), "TEST_VAR")
	requireGetAsFound(t, int8(123), "TEST_VAR")
	requireLookupAsFound(t, int8(123), "TEST_VAR")
	requireLookupAsPFound(t, int8(123), "TEST_VAR")
	requireGetAsFound(t, int16(123), "TEST_VAR")
	requireLookupAsFound(t, int16(123), "TEST_VAR")
	requireLookupAsPFound(t, int16(123), "TEST_VAR")
	requireGetAsFound(t, int32(123), "TEST_VAR")
	requireLookupAsFound(t, int32(123), "TEST_VAR")
	requireLookupAsPFound(t, int32(123), "TEST_VAR")
	requireGetAsFound(t, int64(123), "TEST_VAR")
	requireLookupAsFound(t, int64(123), "TEST_VAR")
	requireLookupAsPFound(t, int64(123), "TEST_VAR")
	requireGetAsFound(t, uint(123), "TEST_VAR")
	requireLookupAsFound(t, uint(123), "TEST_VAR")
	requireLookupAsPFound(t, uint(123), "TEST_VAR")
	requireGetAsFound(t, uint8(123), "TEST_VAR")
	requireLookupAsFound(t, uint8(123), "TEST_VAR")
	requireLookupAsPFound(t, uint8(123), "TEST_VAR")
	requireGetAsFound(t, uint16(123), "TEST_VAR")
	requireLookupAsFound(t, uint16(123), "TEST_VAR")
	requireLookupAsPFound(t, uint16(123), "TEST_VAR")
	requireGetAsFound(t, uint32(123), "TEST_VAR")
	requireLookupAsFound(t, uint32(123), "TEST_VAR")
	requireLookupAsPFound(t, uint32(123), "TEST_VAR")
	requireGetAsFound(t, uint64(123), "TEST_VAR")
	requireLookupAsFound(t, uint64(123), "TEST_VAR")
	requireLookupAsPFound(t, uint64(123), "TEST_VAR")
	requireGetAsFound(t, float32(123.0), "TEST_VAR")
	requireLookupAsFound(t, float32(123.0), "TEST_VAR")
	requireLookupAsPFound(t, float32(123.0), "TEST_VAR")
	requireGetAsFound(t, float64(123.0), "TEST_VAR")
	requireLookupAsFound(t, float64(123.0), "TEST_VAR")
	requireLookupAsPFound(t, float64(123.0), "TEST_VAR")
	requireGetAsFound(t, "123", "TEST_VAR")
	requireLookupAsFound(t, "123", "TEST_VAR")
	requireLookupAsPFound(t, "123", "TEST_VAR")
	requireGetAsError[bool](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[bool](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[bool](t, env.ErrValueParsingFailed, "TEST_VAR")
}

func TestGetLookupAs_Floats(t *testing.T) {
	v := env.NewVar("TEST_VAR")
	defer v.MustRestore()
	require.NoError(t, v.Set("123.456"))

	requireGetAsFound(t, float32(123.456), "TEST_VAR")
	requireLookupAsFound(t, float32(123.456), "TEST_VAR")
	requireLookupAsPFound(t, float32(123.456), "TEST_VAR")
	requireGetAsFound(t, float64(123.456), "TEST_VAR")
	requireLookupAsFound(t, float64(123.456), "TEST_VAR")
	requireLookupAsPFound(t, float64(123.456), "TEST_VAR")
	requireGetAsFound(t, "123.456", "TEST_VAR")
	requireLookupAsFound(t, "123.456", "TEST_VAR")
	requireLookupAsPFound(t, "123.456", "TEST_VAR")
	requireGetAsError[int](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[int](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[int](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireGetAsError[int8](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[int8](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[int8](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireGetAsError[int16](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[int16](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[int16](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireGetAsError[int32](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[int32](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[int32](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireGetAsError[int64](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[int64](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[int64](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireGetAsError[uint](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[uint](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[uint](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireGetAsError[uint8](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[uint8](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[uint8](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireGetAsError[uint16](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[uint16](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[uint16](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireGetAsError[uint32](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[uint32](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[uint32](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireGetAsError[uint64](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[uint64](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[uint64](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireGetAsError[bool](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsError[bool](t, env.ErrValueParsingFailed, "TEST_VAR")
	requireLookupAsPError[bool](t, env.ErrValueParsingFailed, "TEST_VAR")
}

func TestGetLookupAs_Bools(t *testing.T) {
	var (
		truthy = []string{
			"1",
			"t",
			"T",
			"True",
			"TRUE",
			"y",
			"Y",
			"Yes",
			"YES",
		}
		falsy = []string{
			"0",
			"f",
			"F",
			"False",
			"FALSE",
			"n",
			"N",
			"No",
			"NO",
		}
	)

	t.Run("truthy", func(t *testing.T) {
		for _, str := range truthy {
			func() {
				v := env.NewVar("TEST_VAR")
				defer v.MustRestore()
				require.NoError(t, v.Set(str))
				requireGetAsFound(t, true, "TEST_VAR")
				requireLookupAsFound(t, true, "TEST_VAR")
				requireLookupAsPFound(t, true, "TEST_VAR")
				requireGetAsFound(t, str, "TEST_VAR")
				requireLookupAsFound(t, str, "TEST_VAR")
				requireLookupAsPFound(t, str, "TEST_VAR")
			}()
		}
	})

	t.Run("falsy", func(t *testing.T) {
		for _, str := range falsy {
			func() {
				v := env.NewVar("TEST_VAR")
				defer v.MustRestore()
				require.NoError(t, v.Set(str))
				requireGetAsFound(t, false, "TEST_VAR")
				requireLookupAsFound(t, false, "TEST_VAR")
				requireLookupAsPFound(t, false, "TEST_VAR")
				requireGetAsFound(t, str, "TEST_VAR")
				requireLookupAsFound(t, str, "TEST_VAR")
				requireLookupAsPFound(t, str, "TEST_VAR")
			}()
		}
	})
}

func TestGetLookupAs_NotFound(t *testing.T) {
	name := env.SanitizeName(t.Name())
	require.Zero(t, env.Get(name))
	requireGetAsNotFound[string](t, name)
	requireLookupNotFound(t, name)
	requireLookupPNotFound(t, name)
	requireLookupAsNotFound[string](t, name)
	requireLookupAsPNotFound[string](t, name)
}

type vars []*env.Var

func (v vars) Restore() error {
	var errs []error
	for _, x := range v {
		if x != nil {
			if err := x.Restore(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

func (v vars) MustRestore() {
	if err := v.Restore(); err != nil {
		panic(err)
	}
}

func requireGetAsFound[T env.Parsable](
	t *testing.T,
	want T,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	have, err := env.GetAs[T](name, opts...)
	require.NoError(t, err)
	require.Equal(t, want, have)
}

func requireGetAsNotFound[T env.Parsable](
	t *testing.T,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	have, err := env.GetAs[T](name, opts...)
	require.NoError(t, err)
	require.Zero(t, have)
}

func requireGetAsError[T env.Parsable](
	t *testing.T,
	want error,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	have, err := env.GetAs[T](name, opts...)
	require.ErrorIs(t, err, want)
	require.Zero(t, have)
}

func requireLookupFound[T env.Parsable](
	t *testing.T,
	want T,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	x, found := env.Lookup(name, opts...)
	require.True(t, found)
	require.Equal(t, want, x)
}

func requireLookupNotFound(
	t *testing.T,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	x, found := env.Lookup(name, opts...)
	require.False(t, found)
	require.Zero(t, x)
}

func requireLookupPFound(
	t *testing.T,
	want string,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	x, found := env.LookupP(name, opts...)
	require.True(t, found)
	require.NotNil(t, x)
	require.Equal(t, want, *x)
}

func requireLookupPNotFound(
	t *testing.T,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	x, found := env.LookupP(name, opts...)
	require.False(t, found)
	require.Nil(t, x)
}

func requireLookupAsFound[T env.Parsable](
	t *testing.T,
	want T,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	x, found, err := env.LookupAs[T](name, opts...)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, want, x)
}

func requireLookupAsNotFound[T env.Parsable](
	t *testing.T,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	x, found, err := env.LookupAs[T](name, opts...)
	require.NoError(t, err)
	require.False(t, found)
	require.Zero(t, x)
}

func requireLookupAsError[T env.Parsable](
	t *testing.T,
	want error,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	x, found, err := env.LookupAs[T](name, opts...)
	require.ErrorIs(t, err, want)
	require.True(t, found)
	require.Zero(t, x)
}

func requireLookupAsPFound[T env.Parsable](
	t *testing.T,
	want T,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	x, found, err := env.LookupAsP[T](name, opts...)
	require.NoError(t, err)
	require.True(t, found)
	require.NotNil(t, x)
	require.Equal(t, want, *x)
}

func requireLookupAsPNotFound[T env.Parsable](
	t *testing.T,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	x, found, err := env.LookupAsP[T](name, opts...)
	require.NoError(t, err)
	require.False(t, found)
	require.Nil(t, x)
}

func requireLookupAsPError[T env.Parsable](
	t *testing.T,
	want error,
	name string,
	opts ...env.LookupOption,
) {
	t.Helper()
	x, found, err := env.LookupAsP[T](name, opts...)
	require.ErrorIs(t, err, want)
	require.True(t, found)
	require.Nil(t, x)
}
