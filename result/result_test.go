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

package result_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/result"
)

func TestResult_ZeroValue(t *testing.T) {
	var r result.Result[int]
	require.False(t, r.HasValue())
	v, ok := r.Value()
	require.False(t, ok)
	require.Zero(t, v)
	require.Equal(t, 987, r.ValueOr(987))
	require.Equal(t, 654, r.ValueOrElse(func() int { return 654 }))

	wantErr := errors.New(t.Name())
	require.False(t, r.HasErr())
	e, ok := r.Err()
	require.False(t, ok)
	require.Nil(t, e)
	require.ErrorIs(t, r.ErrOr(wantErr), wantErr)
	require.ErrorIs(t, r.ErrOrElse(func() error { return wantErr }), wantErr)
}

func TestOk(t *testing.T) {
	r := result.Ok(123)
	require.True(t, r.HasValue())
	require.False(t, r.HasErr())

	v, ok := r.Value()
	require.True(t, ok)
	require.Equal(t, 123, v)
	require.Equal(t, 123, r.ValueOr(456))
	require.Equal(t, 123, r.ValueOrElse(func() int { return 456 }))

	e, ok := r.Err()
	require.False(t, ok)
	require.Nil(t, e)
}

func TestErr(t *testing.T) {
	wantErr := errors.New(t.Name())
	r := result.Err[int](wantErr)
	require.False(t, r.HasValue())
	require.True(t, r.HasErr())

	e, ok := r.Err()
	require.True(t, ok)
	require.ErrorIs(t, e, wantErr)
	require.ErrorIs(t, r.ErrOr(errors.New("nope")), wantErr)
	require.ErrorIs(t,
		r.ErrOrElse(func() error { return errors.New("nope") }),
		wantErr,
	)

	v, ok := r.Value()
	require.False(t, ok)
	require.Zero(t, v)
}
