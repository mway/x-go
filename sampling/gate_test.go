// Copyright (c) 2022 Matt Way
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

package sampling_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/sampling"
)

func TestGate_On_Try(t *testing.T) {
	var heads int
	for i := 0; i < 1_000_000; i++ {
		if sampling.On.Try() {
			heads++
		}
	}

	require.Equal(t, 1_000_000, heads)
}

func TestGate_Off_Try(t *testing.T) {
	var heads int
	for i := 0; i < 1_000_000; i++ {
		if sampling.Off.Try() {
			heads++
		}
	}

	require.Equal(t, 0, heads)
}

func TestGate_Coin_Try(t *testing.T) {
	var heads int
	for i := 0; i < 1_000_000; i++ {
		if sampling.Coin.Try() {
			heads++
		}
	}

	require.InDelta(t, 500_000, heads, 100_000)
}
