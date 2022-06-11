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

// Package sync provides synchronization tools, primitives, and wrappers.
package sync

import (
	"sync"
	"sync/atomic"
)

// WaitGroup is a wrapper around sync.WaitGroup that adds a length component.
// It is a drop-in replacement that is functionally equivalent in every way
// except that it also tracks the value of the underlying WaitGroup counter.
type WaitGroup struct {
	wg sync.WaitGroup
	n  [1]int64
}

// Add adds delta, which may be negative, to the WaitGroup counter. If the
// counter becomes zero, all goroutines blocked on Wait are released. If the
// counter goes negative, Add panics.
func (g *WaitGroup) Add(delta int) {
	atomic.AddInt64(&g.n[0], int64(delta))
	g.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one.
func (g *WaitGroup) Done() {
	atomic.AddInt64(&g.n[0], -1)
	g.wg.Done()
}

// Inc increments the WaitGroup counter by one.
func (g *WaitGroup) Inc() {
	g.Add(1)
}

// Len returns the current value of the underlying WaitGroup counter.
func (g *WaitGroup) Len() int {
	return int(atomic.LoadInt64(&g.n[0]))
}

// Wait blocks until the WaitGroup counter is zero.
func (g *WaitGroup) Wait() {
	g.wg.Wait()
}
