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

package strconv_test

import (
	gostrconv "strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/strconv"
)

func TestParseBool(t *testing.T) {
	t.Run("truthy", func(t *testing.T) {
		strs := []string{
			"1", "t", "true", "True", "TRUE", // stdlib
			"y", "yes", "Y", "Yes", "YES", // this lib
		}
		for _, str := range strs {
			have, err := strconv.ParseBool(str)
			require.NoError(t, err)
			require.True(t, have)

			have, err = strconv.Parse[bool](str)
			require.NoError(t, err)
			require.True(t, have)
		}
	})

	t.Run("falsy", func(t *testing.T) {
		strs := []string{
			"0", "f", "false", "False", "FALSE", // stdlib
			"n", "no", "N", "No", "NO", // this lib
		}
		for _, str := range strs {
			have, err := strconv.ParseBool(str)
			require.NoError(t, err)
			require.False(t, have)

			have, err = strconv.Parse[bool](str)
			require.NoError(t, err)
			require.False(t, have)
		}
	})

	t.Run("empty", func(t *testing.T) {
		have, err := strconv.ParseBool("")
		haveErr := &gostrconv.NumError{}
		require.ErrorAs(t, err, &haveErr)
		require.ErrorIs(t, haveErr.Err, gostrconv.ErrSyntax)
		require.False(t, have)

		have, err = strconv.Parse[bool]("")
		haveErr = &gostrconv.NumError{}
		require.ErrorAs(t, err, &haveErr)
		require.ErrorIs(t, haveErr.Err, gostrconv.ErrSyntax)
		require.False(t, have)
	})

	t.Run("invalid", func(t *testing.T) {
		have, err := strconv.ParseBool("invalid")
		haveErr := &gostrconv.NumError{}
		require.ErrorAs(t, err, &haveErr)
		require.ErrorIs(t, haveErr.Err, gostrconv.ErrSyntax)
		require.False(t, have)

		have, err = strconv.Parse[bool]("invalid")
		haveErr = &gostrconv.NumError{}
		require.ErrorAs(t, err, &haveErr)
		require.ErrorIs(t, haveErr.Err, gostrconv.ErrSyntax)
		require.False(t, have)
	})
}

func TestParseInt(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		must(t, 123)(strconv.ParseInt[int]("123"))
		must(t, 123)(strconv.Parse[int]("123"))
		must(t, 123)(strconv.ParseInt[int]("   123   "))
		must[int8](t, 123)(strconv.ParseInt[int8]("123"))
		must[int8](t, 123)(strconv.Parse[int8]("123"))
		must[int16](t, 123)(strconv.ParseInt[int16]("123"))
		must[int16](t, 123)(strconv.Parse[int16]("123"))
		must[int32](t, 123)(strconv.ParseInt[int32]("123"))
		must[int32](t, 123)(strconv.Parse[int32]("123"))
		must[int64](t, 123)(strconv.ParseInt[int64]("123"))
		must[int64](t, 123)(strconv.Parse[int64]("123"))
	})

	t.Run("invalid", func(t *testing.T) {
		want := gostrconv.ErrSyntax
		mustError[int](t, want)(strconv.ParseInt[int]("invalid"))
		mustError[int](t, want)(strconv.Parse[int]("invalid"))
		mustError[int](t, want)(strconv.ParseInt[int]("123.0"))
		mustError[int](t, want)(strconv.Parse[int]("123.0"))
		mustError[int](t, want)(strconv.ParseInt[int]("true"))
		mustError[int](t, want)(strconv.Parse[int]("true"))
	})
}

func TestParseUint(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		must[uint](t, 123)(strconv.ParseUint[uint]("123"))
		must[uint](t, 123)(strconv.Parse[uint]("123"))
		must[uint](t, 123)(strconv.ParseUint[uint]("   123   "))
		must[uint8](t, 123)(strconv.ParseUint[uint8]("123"))
		must[uint8](t, 123)(strconv.Parse[uint8]("123"))
		must[uint16](t, 123)(strconv.ParseUint[uint16]("123"))
		must[uint16](t, 123)(strconv.Parse[uint16]("123"))
		must[uint32](t, 123)(strconv.ParseUint[uint32]("123"))
		must[uint32](t, 123)(strconv.Parse[uint32]("123"))
		must[uint64](t, 123)(strconv.ParseUint[uint64]("123"))
		must[uint64](t, 123)(strconv.Parse[uint64]("123"))
	})

	t.Run("invalid", func(t *testing.T) {
		want := gostrconv.ErrSyntax
		mustError[uint](t, want)(strconv.ParseUint[uint]("invalid"))
		mustError[uint](t, want)(strconv.Parse[uint]("invalid"))
		mustError[uint](t, want)(strconv.ParseUint[uint]("123.0"))
		mustError[uint](t, want)(strconv.Parse[uint]("123.0"))
		mustError[uint](t, want)(strconv.ParseUint[uint]("true"))
		mustError[uint](t, want)(strconv.Parse[uint]("true"))
	})
}

func TestParseFloat(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		must[float32](t, 123)(strconv.ParseFloat[float32]("123"))
		must[float32](t, 123)(strconv.Parse[float32]("123"))
		must[float32](t, 123.456)(strconv.ParseFloat[float32]("123.456"))
		must[float32](t, 123.456)(strconv.Parse[float32]("123.456"))
		must[float32](t, 123)(strconv.ParseFloat[float32]("   123   "))
		must[float32](t, 123)(strconv.Parse[float32]("   123   "))
		must[float32](t, 123.456)(strconv.ParseFloat[float32]("   123.456   "))
		must[float32](t, 123.456)(strconv.Parse[float32]("   123.456   "))
		must[float64](t, 123)(strconv.ParseFloat[float64]("123"))
		must[float64](t, 123)(strconv.Parse[float64]("123"))
		must(t, 123.456)(strconv.ParseFloat[float64]("123.456"))
		must(t, 123.456)(strconv.Parse[float64]("123.456"))
	})

	t.Run("invalid", func(t *testing.T) {
		want := gostrconv.ErrSyntax
		mustError[float64](t, want)(strconv.ParseFloat[float64]("invalid"))
		mustError[float64](t, want)(strconv.Parse[float64]("invalid"))
		mustError[float64](t, want)(strconv.ParseFloat[float64]("true"))
		mustError[float64](t, want)(strconv.Parse[float64]("true"))
		mustError[float32](t, want)(strconv.ParseFloat[float32]("   123 .456   "))
		mustError[float32](t, want)(strconv.Parse[float32]("   123 .456   "))
	})
}

func TestParse_String(t *testing.T) {
	must(t, t.Name())(strconv.Parse[string](t.Name()))
}

func must[T any](t *testing.T, want T) func(T, error) T {
	t.Helper()
	return func(have T, err error) T {
		t.Helper()
		require.NoError(t, err)
		require.Equal(t, want, have)
		return have
	}
}

func mustError[T any](t *testing.T, want error) func(T, error) {
	return func(x T, err error) {
		require.ErrorIs(t, err, want)
		require.Zero(t, x)
	}
}
