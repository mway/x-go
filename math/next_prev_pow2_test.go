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

package math_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/math"
)

func TestNextPow2(t *testing.T) {
	cases := map[int]int{
		1:    1,
		2:    2,
		3:    4,
		4:    4,
		5:    8,
		6:    8,
		7:    8,
		8:    8,
		9:    16,
		10:   16,
		11:   16,
		12:   16,
		13:   16,
		14:   16,
		15:   16,
		16:   16,
		1917: 2048,
	}
	for give, want := range cases {
		t.Run(strconv.Itoa(give), func(t *testing.T) {
			require.Equal(t, want, math.NextPow2(give))
		})
	}
}

func TestPrevPow2(t *testing.T) {
	cases := map[int]int{
		1:    1,
		2:    2,
		3:    2,
		4:    4,
		5:    4,
		6:    4,
		7:    4,
		8:    8,
		9:    8,
		10:   8,
		11:   8,
		12:   8,
		13:   8,
		14:   8,
		15:   8,
		16:   16,
		1917: 1024,
	}
	for give, want := range cases {
		t.Run(strconv.Itoa(give), func(t *testing.T) {
			require.Equal(t, want, math.PrevPow2(give))
		})
	}
}

func BenchmarkNextPow2(b *testing.B) {
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			_ = math.NextPow2(1917)
		}
	})
}

func BenchmarkPrevPow2(b *testing.B) {
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			_ = math.PrevPow2(1917)
		}
	})
}
