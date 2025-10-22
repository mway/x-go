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

func TestClamp(t *testing.T) {
	t.Run("int8", func(t *testing.T) {
		for i := range 10 {
			want := max(0, min(i, 3))
			require.EqualValues(t, want, clamp.Clamp(int8(i), 0, 3))
		}
	})
	t.Run("uint8", func(t *testing.T) {
		for i := range 10 {
			want := max(0, min(i, 3))
			require.EqualValues(t, want, clamp.Clamp(uint8(i), 0, 3))
		}
	})
}

func TestAdd(t *testing.T) {
	t.Run("smoke", func(t *testing.T) {
		for i := range 5 {
			require.Equal(t, min(i, 3), clamp.Add(i, 0, 0, 3))
			require.Equal(t, min(i+1, 3), clamp.Add(i, 1, 0, 3))
			require.Equal(t, 3, clamp.Add(i, 5, 0, 3))
		}
	})

	t.Run("wrapping", func(t *testing.T) {
		require.EqualValues(t, 3, clamp.Add[uint8](200, 100, 1, 3))
	})
}

func TestSub(t *testing.T) {
	t.Run("smoke", func(t *testing.T) {
		for i := range 5 {
			require.Equal(t, min(max(i, 1), 3), clamp.Sub(i, 0, 1, 3))
			require.Equal(t, max(i-1, 1), clamp.Sub(i, 1, 1, 3))
			require.Equal(t, 1, clamp.Sub(i, 5, 1, 3))
		}
	})

	t.Run("wrapping", func(t *testing.T) {
		require.EqualValues(t, 1, clamp.Sub[uint8](100, 200, 1, 3))
	})
}
