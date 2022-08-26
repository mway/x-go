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

package sampling

import (
	_ "unsafe" // go:linkname
)

const (
	// On is a Gate that always succeeds when tried.
	On = Gate(1.0)
	// Off is a Gate that always fails when tried.
	Off = Gate(0.0)
	// Coin is a Gate that flips a coin (i.e., 50/50).
	Coin = Gate(0.5)
)

// Gate is a simple sampling gate in the range [0.0, 1.0].
type Gate float64

// Try makes an attempt against g's inherent probability.
func (g Gate) Try() bool {
	const _max = 1 << 24
	return fastrandn(_max) >= _max-uint32(g*_max)
}

//go:linkname fastrandn runtime.fastrandn
func fastrandn(uint32) uint32
