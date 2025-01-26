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

package clamp_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/math/clamp"
)

func TestRange(t *testing.T) {
	require.Panics(t, func() {
		clamp.NewRange(3, 1)
	})

	r := clamp.NewRange(1, 3)
	require.Equal(t, 1, r.Min)
	require.Equal(t, 3, r.Max)
	require.Equal(t, r.Min, r.Clamp(0))
	require.Equal(t, r.Min, r.Clamp(r.Min))
	require.Equal(t, r.Min+1, r.Clamp(r.Min+1))
	require.Equal(t, r.Max, r.Clamp(3))
	require.Equal(t, r.Max, r.Clamp(4))
	require.Equal(t, r.Min, r.Add(0, 0))
	require.Equal(t, r.Min, r.Add(0, 1))
	require.Equal(t, r.Min, r.Add(r.Min, 0))
	require.Equal(t, r.Min+1, r.Add(r.Min, 1))
	require.Equal(t, r.Max, r.Add(r.Max-1, 1))
	require.Equal(t, r.Max, r.Add(r.Max-1, 2))
	require.Equal(t, r.Max, r.Add(r.Max, 0))
	require.Equal(t, r.Max, r.Add(r.Max, 1))
	require.Equal(t, r.Min, r.Sub(r.Min, 1))
	require.Equal(t, r.Min, r.Sub(r.Min, 0))
	require.Equal(t, r.Max-1, r.Sub(r.Max, 1))
	require.Equal(t, r.Max, r.Sub(r.Max, 0))
	require.Equal(t, r.Min, r.Sub(2, 2))
	require.Equal(t, r.Max-2, r.Sub(r.Max, 2))
}
