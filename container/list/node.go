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

// A Node is a singly-linked list node that holds a T value.
type Node[T any] struct {
	Next  *Node[T]
	value T
	isset bool
}

// NewNode creates a new [Node[T]] with the given value.
func NewNode[T any](value T) *Node[T] {
	return &Node[T]{
		value: value,
		isset: true,
	}
}

// Link singly-links value and any extra values in order, returning the head of
// the list.
func Link[T any](value T, extra ...T) *Node[T] {
	x, _ := LinkWithTail(value, extra...)
	return x
}

// LinkWithTail singly-links value and any extra values in order, returning the
// head and tail of the list.
func LinkWithTail[T any](value T, extra ...T) (*Node[T], *Node[T]) {
	var (
		head = NewNode(value)
		tail = head
		cur  = head
	)

	for _, val := range extra {
		cur.Next = &Node[T]{
			value: val,
			isset: true,
		}
		cur = cur.Next
		tail = cur
	}

	return head, tail
}

// Value returns the node's value.
func (n *Node[T]) Value() T {
	x, _ := n.Get()
	return x
}

// Get returns the node's value. The boolean return indicates if the node has a
// value set.
func (n *Node[T]) Get() (T, bool) {
	if !n.isset {
		var zero T
		return zero, false
	}
	return n.value, true
}

// Set sets the node's value to the given value.
func (n *Node[T]) Set(value T) {
	n.value = value
	n.isset = true
}

// Unset unsets the node.
func (n *Node[T]) Unset() {
	var zero T
	n.value = zero
	n.isset = false
}

// IsSet indicates if the node has been set.
func (n *Node[T]) IsSet() bool {
	return n != nil && n.isset
}

// Pop removes the node if it is set, and replaces it with the next node, if
// any. If the node is not set, Pop does nothing. If there is not a next node,
// the current node becomes unset. The returned value is the popped value, if
// a value was popped, or a zero value T otherwise.
func (n *Node[T]) Pop() T {
	x, _ := n.MaybePop()
	return x
}

// MaybePop removes the node if it is set, and replaces it with the next node,
// if any. If the node is not set, MaybePop does nothing. If there is not a
// next node, the current node becomes unset. The returned value is the popped
// value, and the boolean indicates whether a value was popped.
func (n *Node[T]) MaybePop() (T, bool) {
	var x T
	if !n.isset {
		return x, false
	}

	n.isset = false
	if x, n.value = n.value, x; n.Next != nil {
		*n = *n.Next
	}
	return x, true
}

// InsertBefore inserts the given node or list before n.
func (n *Node[T]) InsertBefore(node *Node[T]) {
	end := node
	for end.Next != nil {
		end = end.Next
	}

	tmp := *n
	end.Next = &tmp

	*n = *node
}

// InsertValueBefore inserts a new node for value before n and returns it.
func (n *Node[T]) InsertValueBefore(value T) *Node[T] {
	node := NewNode(value)
	n.InsertBefore(node)
	return node
}

// InsertAfter inserts the given node or list after n.
func (n *Node[T]) InsertAfter(node *Node[T]) {
	if n.Next != nil {
		end := node
		for end.Next != nil {
			end = end.Next
		}
		end.Next = n.Next
	}

	n.Next = node
}

// InsertValueAfter inserts a new node for value after n and returns it.
func (n *Node[T]) InsertValueAfter(value T) *Node[T] {
	node := NewNode(value)
	n.InsertAfter(node)
	return node
}

// Append appends node to the end of n's list.
func (n *Node[T]) Append(node *Node[T]) {
	end := n
	for end.Next != nil {
		end = end.Next
	}
	end.Next = node
}

// AppendValue appends a new node for value to the end of n's list.
func (n *Node[T]) AppendValue(value T) *Node[T] {
	node := NewNode(value)
	n.Append(node)
	return node
}

// Sort sorts the list with head n according to pred.
func (n *Node[T]) Sort(pred func(T, T) int) {
	if !n.isset {
		return
	}

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
func (n *Node[T]) ToSlice() []T {
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

// Iter returns a single-value iterator for all values contained in the list
// with head n.
func (n *Node[T]) Iter() iter.Seq[T] {
	if !n.isset {
		return func(func(T) bool) {}
	}

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
func (n *Node[T]) ForEach(walk func(*Node[T]) bool) {
	if !n.isset {
		return
	}

	cur := n
	for cur != nil {
		if !cur.isset || !walk(cur) {
			break
		}
		cur = cur.Next
	}
}

// ForEachValue walks the values in the list with head n, starting with n.
func (n *Node[T]) ForEachValue(walk func(T) bool) {
	n.ForEach(func(x *Node[T]) bool {
		return walk(x.value)
	})
}
