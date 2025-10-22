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

func TestValue(t *testing.T) {
	require.Panics(t, func() {
		clamp.NewValue(0, 3, 1)
	})

	vs := []clamp.Value[int]{
		clamp.NewValue(0, 1, 3),
		clamp.NewValueWithRange(0, clamp.NewRange(1, 3)),
	}
	for _, v := range vs {
		require.Equal(t, v.Min(), v.Value())
		v.Sub(1)
		require.Equal(t, v.Min(), v.Value())
		v.Add(1)
		require.Equal(t, v.Min()+1, v.Value())
		v.Set(v.Max() + 1)
		require.Equal(t, v.Max(), v.Value())
		v.Set(v.Max() - 1)
		require.Equal(t, v.Max()-1, v.Value())
		require.True(t, v.InRange(v.Min()))
		require.True(t, v.InRange(v.Max()))

		v.SetMin(v.Min() + 1)
		require.Equal(t, 2, v.Min())
		require.Equal(t, v.Min(), v.Value())

		v.SetMax(v.Max() - 1)
		require.Equal(t, 2, v.Max())
		require.Equal(t, v.Max(), v.Value())
	}
}
