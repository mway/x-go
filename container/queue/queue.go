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

// Package queue provides queue-based types and helpers.
package queue

// A Queue is a FIFO queue that holds values of type T.
type Queue[T any] struct {
	data []T
}

// NewQueue creates a new [Queue[T]] with the given initial capacity.
func NewQueue[T any](size int) *Queue[T] {
	return &Queue[T]{
		data: make([]T, 0, size),
	}
}

// Push pushes x to the back of the queue.
func (q *Queue[T]) Push(x T) {
	q.data = append(q.data, x)
}

// Front returns the value at the front of the queue.
func (q *Queue[T]) Front() T {
	return q.data[0]
}

// Pop pops the front value off of the queue and returns it.
func (q *Queue[T]) Pop() T {
	x := q.data[0]
	copy(q.data, q.data[1:])
	q.data = q.data[:len(q.data)-1]
	return x
}

// Len returns the number of values held by the queue.
func (q *Queue[T]) Len() int {
	return len(q.data)
}

// Cap returns the queue's current capacity.
func (q *Queue[T]) Cap() int {
	return cap(q.data)
}

// MaybeFront returns the value at the front of the queue if there is one. The
// boolean return indicates whether the T value is valid.
func (q *Queue[T]) MaybeFront() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}
	return q.Front(), true
}

// MaybePop pops the top value off of the front of the queue and returns it, if
// there is one. The boolean return indicates whether the T is valid.
func (q *Queue[T]) MaybePop() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}
	return q.Pop(), true
}

// PeekEach calls fn for each value in the queue. If fn returns false,
// iteration will stop and the function will return immediately. If there are
// no values in the queue, fn will not be called.
func (q *Queue[T]) PeekEach(fn func(T) bool) {
	if len(q.data) == 0 {
		return
	}
	for i := 0; i < len(q.data); i++ {
		if !fn(q.data[i]) {
			return
		}
	}
}

// PopEach pops the front value off of the front of the queue and passes it to
// fn while the queue is not empty. If fn returns false, no more values will be
// popped and the function will return immediately. If there are no values in
// the queue, fn will not be called.
func (q *Queue[T]) PopEach(fn func(T) bool) {
	for {
		if x, ok := q.MaybePop(); !ok || !fn(x) {
			break
		}
	}
}
