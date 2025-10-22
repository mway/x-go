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

// Package deque provides deque-based types and helpers.
package deque

import (
	"slices"

	"go.mway.dev/pool"

	"go.mway.dev/x/container/list"
)

// A Deque is a double-ended (FIFO and LIFO) queue that holds values of type T.
type Deque[T any] struct {
	data []T
}

// New creates a new [Deque[T]] with the given initial capacity.
func New[T any](size int) *Deque[T] {
	return &Deque[T]{
		data: make([]T, 0, size),
	}
}

// NewWithValues creates a new [Deque[T]] with the given initial values.
func NewWithValues[T any](values ...T) *Deque[T] {
	return &Deque[T]{
		data: slices.Clone(values),
	}
}

// PushFront pushes x to the front of the deque.
func (d *Deque[T]) PushFront(x T) {
	d.data = append(d.data, x)
	copy(d.data[1:], d.data)
	d.data[0] = x
}

// Front returns the value at the front of the deque.
func (d *Deque[T]) Front() T {
	x, _ := d.MaybeFront()
	return x
}

// MaybeFront returns the value at the front of the deque if there is one. The
// boolean return indicates whether the T value is valid.
func (d *Deque[T]) MaybeFront() (T, bool) {
	if len(d.data) == 0 {
		var zero T
		return zero, false
	}
	return d.data[0], true
}

// PopFront pops the value off of the front of the deque and returns it.
func (d *Deque[T]) PopFront() T {
	x, _ := d.MaybePopFront()
	return x
}

// MaybePopFront pops the value off of the front of the deque and returns it,
// if there is one. The boolean return indicates whether the T is valid.
func (d *Deque[T]) MaybePopFront() (T, bool) {
	if len(d.data) == 0 {
		var zero T
		return zero, false
	}

	x := d.data[0]
	d.data = d.data[1:]
	return x, true
}

// PeekEachFront yields each value in the deque to fn, working from the front
// of the deque to the back. If fn returns false, iteration will stop and
// control will return to the caller immediately. If there are no values in the
// deque, fn will not be called.
func (d *Deque[T]) PeekEachFront(fn func(T) bool) {
	if len(d.data) == 0 {
		return
	}

	for i := 0; i < len(d.data); i++ {
		if !fn(d.data[i]) {
			return
		}
	}
}

// PopEachFront pops and yields each value in the deque to fn, working from the
// front of the deque to the back. If fn returns false, iteration will stop and
// control will return to the caller immediately. If there are no values in the
// deque, fn will not be called.
func (d *Deque[T]) PopEachFront(fn func(T) bool) {
	for {
		if x, ok := d.MaybePopFront(); !ok || !fn(x) {
			break
		}
	}
}

// PushBack pushes x to the back of the deque.
func (d *Deque[T]) PushBack(x T) {
	d.data = append(d.data, x)
}

// Back returns the value at the back of the deque.
func (d *Deque[T]) Back() T {
	x, _ := d.MaybeBack()
	return x
}

// MaybeBack returns the value at the back of the deque if there is one. The
// boolean return indicates whether the T value is valid.
func (d *Deque[T]) MaybeBack() (T, bool) {
	if len(d.data) == 0 {
		var zero T
		return zero, false
	}

	return d.data[len(d.data)-1], true
}

// PopBack pops the back value off of the deque and returns it.
func (d *Deque[T]) PopBack() T {
	x, _ := d.MaybePopBack()
	return x
}

// MaybePopBack pops the top value off of the front of the deque and returns
// it, if there is one. The boolean return indicates whether the T is valid.
func (d *Deque[T]) MaybePopBack() (T, bool) {
	if len(d.data) == 0 {
		var zero T
		return zero, false
	}

	var (
		n = len(d.data) - 1
		x = d.data[n]
	)
	d.data = d.data[:n]
	return x, true
}

// PeekEachBack yields each value in the deque to fn, working from the back of
// the deque to the front. If fn returns false, iteration will stop and control
// will return to the caller immediately. If there are no values in the deque,
// fn will not be called.
func (d *Deque[T]) PeekEachBack(fn func(T) bool) {
	if len(d.data) == 0 {
		return
	}

	for i := len(d.data) - 1; i >= 0; i-- {
		if !fn(d.data[i]) {
			return
		}
	}
}

// PopEachBack pops and yields each value in the deque to fn, working from the
// back of the deque to the front. If fn returns false, iteration will stop and
// control will return to the caller immediately. If there are no values in the
// deque, fn will not be called.
func (d *Deque[T]) PopEachBack(fn func(T) bool) {
	for {
		if x, ok := d.MaybePopBack(); !ok || !fn(x) {
			break
		}
	}
}

// Len returns the number of values held by the deque.
func (d *Deque[T]) Len() int {
	return len(d.data)
}

// A LinkedDeque is a double-ended (FIFO and LIFO) queue that holds values of
// type T.
type LinkedDeque[T any] struct {
	head *list.DoubleNode[T]
	tail *list.DoubleNode[T]
	pool *pool.Pool[*list.DoubleNode[T]]
	len  int
}

// NewLinked creates a new [LinkedDeque[T]] with the given initial capacity.
func NewLinked[T any]() *LinkedDeque[T] {
	return &LinkedDeque[T]{
		head: nil,
		tail: nil,
		pool: pool.NewWithReleaser(
			func() *list.DoubleNode[T] {
				return &list.DoubleNode[T]{}
			},
			func(x *list.DoubleNode[T]) {
				x.Next = nil
				x.Prev = nil
				x.Unset()
			},
		),
		len: 0,
	}
}

// NewLinkedWithValues creates a new [LinkedDeque[T]] with the given initial
// values.
func NewLinkedWithValues[T any](values ...T) *LinkedDeque[T] {
	if len(values) == 0 {
		return nil
	}

	head, tail := list.LinkDoublyWithTail(values[0], values[1:]...)
	return &LinkedDeque[T]{
		head: head,
		tail: tail,
		len:  len(values),
	}
}

// PushFront pushes x to the front of the deque.
func (d *LinkedDeque[T]) PushFront(value T) {
	d.len++

	node := d.pool.Get()
	node.Set(value)

	if d.head == nil {
		d.head = node
		d.tail = d.head
		return
	}

	d.head.InsertBefore(node)
	d.head = d.head.Prev
}

// Front returns the value at the front of the deque.
func (d *LinkedDeque[T]) Front() T {
	x, _ := d.MaybeFront()
	return x
}

// MaybeFront returns the value at the front of the deque if there is one. The
// boolean return indicates whether the T value is valid.
func (d *LinkedDeque[T]) MaybeFront() (T, bool) {
	if d.head == nil {
		var zero T
		return zero, false
	}
	return d.head.Value(), true
}

// PeekEachFront yields each value in the deque to fn, working from the front
// of the deque to the back. If fn returns false, iteration will stop and
// control will return to the caller immediately. If there are no values in the
// deque, fn will not be called.
func (d *LinkedDeque[T]) PeekEachFront(fn func(T) bool) {
	cur := d.head
	for cur != nil {
		if !fn(cur.Value()) {
			return
		}
		cur = cur.Next
	}
}

// PopFront pops the value off of the front of the deque and returns it.
func (d *LinkedDeque[T]) PopFront() T {
	x, _ := d.MaybePopFront()
	return x
}

// MaybePopFront pops the value off of the front of the deque and returns it,
// if there is one. The boolean return indicates whether the T is valid.
func (d *LinkedDeque[T]) MaybePopFront() (T, bool) {
	if d.head == nil {
		var zero T
		return zero, false
	}

	x := d.head
	defer d.pool.Put(x)

	if d.head = d.head.Next.DetachPrev(); d.head == nil {
		d.tail = nil
	}
	d.len--
	return x.Value(), true
}

// PopEachFront pops and yields each value in the deque to fn, working from the
// front of the deque to the back. If fn returns false, iteration will stop and
// control will return to the caller immediately. If there are no values in the
// deque, fn will not be called.
func (d *LinkedDeque[T]) PopEachFront(fn func(T) bool) {
	for {
		if x, ok := d.MaybePopFront(); !ok || !fn(x) {
			break
		}
	}
}

// PushBack pushes x to the back of the deque.
func (d *LinkedDeque[T]) PushBack(value T) {
	d.len++

	node := d.pool.Get()
	node.Set(value)

	if d.head == nil {
		d.head = node
		d.tail = d.head
		return
	}

	d.tail.InsertAfter(node)
	d.tail = d.tail.Next
}

// Back returns the value at the back of the deque.
func (d *LinkedDeque[T]) Back() T {
	x, _ := d.MaybeBack()
	return x
}

// MaybeBack returns the value at the back of the deque if there is one. The
// boolean return indicates whether the T value is valid.
func (d *LinkedDeque[T]) MaybeBack() (T, bool) {
	if d.tail == nil {
		var zero T
		return zero, false
	}
	return d.tail.Value(), true
}

// PopBack pops the back value off of the deque and returns it.
func (d *LinkedDeque[T]) PopBack() T {
	x, _ := d.MaybePopBack()
	return x
}

// MaybePopBack pops the top value off of the front of the deque and returns
// it, if there is one. The boolean return indicates whether the T is valid.
func (d *LinkedDeque[T]) MaybePopBack() (T, bool) {
	if d.tail == nil {
		var zero T
		return zero, false
	}

	x := d.tail
	defer d.pool.Put(x)

	if d.tail = d.tail.Prev.DetachNext(); d.tail == nil {
		d.head = nil
	}
	d.len--
	return x.Value(), true
}

// PeekEachBack yields each value in the deque to fn, working from the back of
// the deque to the front. If fn returns false, iteration will stop and control
// will return to the caller immediately. If there are no values in the deque,
// fn will not be called.
func (d *LinkedDeque[T]) PeekEachBack(fn func(T) bool) {
	cur := d.tail
	for cur != nil {
		if !fn(cur.Value()) {
			return
		}
		cur = cur.Prev
	}
}

// PopEachBack pops and yields each value in the deque to fn, working from the
// back of the deque to the front. If fn returns false, iteration will stop and
// control will return to the caller immediately. If there are no values in the
// deque, fn will not be called.
func (d *LinkedDeque[T]) PopEachBack(fn func(T) bool) {
	for {
		if x, ok := d.MaybePopBack(); !ok || !fn(x) {
			break
		}
	}
}

// Len returns the number of values held by the deque.
func (d *LinkedDeque[T]) Len() int {
	return d.len
}
