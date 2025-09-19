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

// Package heap provides heap-related types and utilities.
package heap

import (
	"cmp"
	"slices"
)

// n.b. Most of this functionality was ported (essentially verbatim) from the
//      Go standard library for parity.

// MinHeap is a min heap (P<=C).
type MinHeap[T cmp.Ordered] struct {
	heap[T, heapTypeMin]
}

// NewMinHeap creates a new [MinHeap] with the given initial values.
func NewMinHeap[T cmp.Ordered](values ...T) *MinHeap[T] {
	return &MinHeap[T]{
		heap: newHeap[T, heapTypeMin](values...),
	}
}

// Min returns the current minimum value on the heap.
func (h *MinHeap[T]) Min() T {
	return h.top()
}

// MaxHeap is a max heap (P>=C).
type MaxHeap[T cmp.Ordered] struct {
	heap[T, heapTypeMax]
}

// NewMaxHeap creates a new [MaxHeap] with the given initial values.
func NewMaxHeap[T cmp.Ordered](values ...T) *MaxHeap[T] {
	return &MaxHeap[T]{
		heap: newHeap[T, heapTypeMax](values...),
	}
}

// Max returns the current maximum value on the heap.
func (h *MaxHeap[T]) Max() T {
	return h.top()
}

// Types used to do static type switching to disambiguate comparison while
// sharing logic.
type (
	heapTypeMin struct{}
	heapTypeMax struct{}
)

type heapType interface {
	heapTypeMin | heapTypeMax
}

type heap[T cmp.Ordered, H heapType] struct {
	data []T
}

func newHeap[T cmp.Ordered, H heapType](values ...T) heap[T, H] {
	h := heap[T, H]{
		data: slices.Clone(values),
	}
	h.init()
	return h
}

func (h *heap[T, H]) Push(value T) {
	h.data = append(h.data, value)
	h.up(h.Len() - 1)
}

func (h *heap[T, H]) Pop() T {
	n := h.Len() - 1
	h.Swap(0, n)
	h.down(0, n)
	x := h.data[n]
	h.data = h.data[:n]
	return x
}

func (h *heap[T, H]) Len() int {
	return len(h.data)
}

func (h *heap[T, H]) Less(i int, j int) bool {
	var zero H
	if _, ok := any(zero).(heapTypeMin); ok {
		return h.data[i] < h.data[j]
	}
	return h.data[i] > h.data[j]
}

func (h *heap[T, H]) Swap(i int, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *heap[T, H]) Remove(i int) T {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	x := h.data[n]
	h.data = h.data[:n]
	return x
}

func (h *heap[T, H]) Reset() {
	h.data = h.data[:0]
}

func (h *heap[T, H]) down(i0 int, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}

func (h *heap[T, H]) init() {
	// heapify
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

func (h *heap[T, H]) top() T {
	var x T
	if len(h.data) > 0 {
		x = h.data[0]
	}
	return x
}

func (h *heap[T, H]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}
