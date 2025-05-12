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

package list

import (
	"iter"
)

// A DoubleNode is a doubly-linked list node that holds a T value.
type DoubleNode[T any] struct {
	Prev  *DoubleNode[T]
	Next  *DoubleNode[T]
	value T
	isset bool
}

// NewDoubleNode creates a new [DoubleNode[T]] with the given value.
func NewDoubleNode[T any](value T) *DoubleNode[T] {
	return &DoubleNode[T]{
		value: value,
		isset: true,
	}
}

// LinkDoubly doubly-links value and any extra values in order, returning the
// head of the list.
func LinkDoubly[T any](value T, extra ...T) *DoubleNode[T] {
	x, _ := LinkDoublyWithTail(value, extra...)
	return x
}

// LinkDoublyWithTail doubly-links value and any extra values in order,
// returning the head and tail of the list.
func LinkDoublyWithTail[T any](
	value T,
	extra ...T,
) (*DoubleNode[T], *DoubleNode[T]) {
	var (
		head = NewDoubleNode(value)
		tail = head
		cur  = head
	)

	for _, val := range extra {
		cur.Next = NewDoubleNode(val).WithPrev(cur)
		cur = cur.Next
		tail = cur
	}
	return head, tail
}

// WithPrev doubly links prev<>n, returning n. It does not adjust the links of
// any original prev.Next or prev.Prev nodes.
func (n *DoubleNode[T]) WithPrev(prev *DoubleNode[T]) *DoubleNode[T] {
	n.Prev = prev
	prev.Next = n
	return n
}

// WithNext doubly links n<>next, returning n. It does not adjust the links of
// any original next.Next or next.Prev nodes.
func (n *DoubleNode[T]) WithNext(next *DoubleNode[T]) *DoubleNode[T] {
	n.Next = next
	next.Prev = n
	return n
}

// Value returns the node's value.
func (n *DoubleNode[T]) Value() T {
	x, _ := n.Get()
	return x
}

// Get returns the node's value. The boolean return indicates if the node has a
// value set.
func (n *DoubleNode[T]) Get() (T, bool) {
	if !n.isset {
		var zero T
		return zero, false
	}

	return n.value, true
}

// Set sets the node's value to the given value.
func (n *DoubleNode[T]) Set(value T) {
	n.value = value
	n.isset = true
}

// Unset unsets the node.
func (n *DoubleNode[T]) Unset() {
	var zero T
	n.value = zero
	n.isset = false
}

// IsSet indicates if the node has been set.
func (n *DoubleNode[T]) IsSet() bool {
	return n != nil && n.isset
}

// DetachPrev detches the node from its previous node, if one exists.
func (n *DoubleNode[T]) DetachPrev() *DoubleNode[T] {
	if n == nil || n.Prev == nil {
		return n
	}
	n.Prev.Next = nil
	n.Prev = nil
	return n
}

// DetachNext detaches the node from its next node, if one exists.
func (n *DoubleNode[T]) DetachNext() *DoubleNode[T] {
	if n == nil || n.Next == nil {
		return n
	}
	n.Next.Prev = nil
	n.Next = nil
	return n
}

// InsertBefore inserts the given node or list before n.
func (n *DoubleNode[T]) InsertBefore(node *DoubleNode[T]) {
	end := node
	for end.Next != nil {
		end = end.Next
	}
	end.Next = n

	if n.Prev != nil {
		n.Prev.Next = node
	}
	node.Prev = n.Prev
	n.Prev = end
}

// InsertValueBefore inserts a new node for value before n and returns it.
func (n *DoubleNode[T]) InsertValueBefore(value T) *DoubleNode[T] {
	node := NewDoubleNode(value)
	n.InsertBefore(node)
	return node
}

// InsertAfter inserts the given node or list after n.
func (n *DoubleNode[T]) InsertAfter(node *DoubleNode[T]) {
	if n.Next != nil {
		end := node
		for end.Next != nil {
			end = end.Next
		}
		end.Next = n.Next
		n.Next.Prev = end
	}

	n.Next = node
	node.Prev = n
}

// InsertValueAfter inserts a new node for value after n and returns it.
func (n *DoubleNode[T]) InsertValueAfter(value T) *DoubleNode[T] {
	node := NewDoubleNode(value)
	n.InsertAfter(node)
	return node
}

// Append appends node to the end of n's list.
func (n *DoubleNode[T]) Append(node *DoubleNode[T]) {
	end := n
	for end.Next != nil {
		end = end.Next
	}
	end.Next = node
	node.Prev = end
}

// AppendValue appends a new node for value to the end of n's list.
func (n *DoubleNode[T]) AppendValue(value T) *DoubleNode[T] {
	node := NewDoubleNode(value)
	n.Append(node)
	return node
}

// Sort sorts the list with head n according to pred.
func (n *DoubleNode[T]) Sort(pred func(T, T) int) {
	cur := n
	for cur != nil {
		peek := cur.Next
		for peek != nil {
			if pred(cur.value, peek.value) > 0 {
				cur.value, peek.value = peek.value, cur.value
			}
			peek = peek.Next
		}
		cur = cur.Next
	}
}

// ToSlice returns a slice of all of values contained in the list with head n.
func (n *DoubleNode[T]) ToSlice() []T {
	if !n.isset {
		return nil
	}

	var (
		values = []T{n.value}
		cur    = n.Next
	)
	for cur != nil {
		x, ok := cur.Get()
		if !ok {
			break
		}
		values = append(values, x)
		cur = cur.Next
	}
	return values
}

// ToSliceRev returns a slice of all of values contained in the list with n as
// the tail in reverse order.
func (n *DoubleNode[T]) ToSliceRev() []T {
	if !n.isset {
		return nil
	}

	var (
		values = []T{n.value}
		cur    = n.Prev
	)
	for cur != nil {
		x, ok := cur.Get()
		if !ok {
			break
		}
		values = append(values, x)
		cur = cur.Prev
	}
	return values
}

// Iter returns a single-value iterator for all values contained in the list
// with head n.
func (n *DoubleNode[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		cur := n
		for cur != nil {
			if x, ok := cur.Get(); !ok || !yield(x) {
				break
			}
			cur = cur.Next
		}
	}
}

// ForEach walks the nodes in the list with head n, starting with n.
func (n *DoubleNode[T]) ForEach(walk func(*DoubleNode[T]) bool) {
	cur := n
	for cur != nil {
		if !cur.isset || !walk(cur) {
			break
		}
		cur = cur.Next
	}
}

// ForEachValue walks the values in the list with head n, starting with n.
func (n *DoubleNode[T]) ForEachValue(walk func(T) bool) {
	n.ForEach(func(x *DoubleNode[T]) bool {
		return walk(x.Value())
	})
}
