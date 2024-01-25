// Copyright (c) 2024 Matt Way
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

// Package tree provides tree structure related types and utilities.
package tree

import (
	"errors"
	"maps"
	"slices"
	"sort"

	"golang.org/x/exp/constraints"
	xmaps "golang.org/x/exp/maps"
)

// ErrSkipSubtree is a sentinel value that can be returned by a [NodeWalker] to
// skip the remainder of a subtree. This error halts iteration of that subtree
// at the first handling parent.
var ErrSkipSubtree = errors.New("skip subtree")

// A NodeWalker is a function that controls how nodes are walked.
type NodeWalker[K constraints.Ordered, V comparable] func(node *BasicNode[K, V]) error

// BasicNode is a basic, arbitrarily-ordered, non-balancing, key/value tree.
type BasicNode[K constraints.Ordered, V comparable] struct {
	key      K
	value    V
	parent   *BasicNode[K, V]
	children map[K]*BasicNode[K, V]
}

// NewBasicNode creates a new [BasicNode].
func NewBasicNode[K constraints.Ordered, V comparable](
	key K,
	value V,
) *BasicNode[K, V] {
	return &BasicNode[K, V]{
		key:   key,
		value: value,
	}
}

// Key returns the node's key.
func (n *BasicNode[K, V]) Key() (key K) {
	if n == nil {
		return
	}
	return n.key
}

// Path returns all keys in the path from the root node to this node.
func (n *BasicNode[K, V]) Path() []K {
	keys := n.PathRev()
	slices.Reverse(keys)
	return keys
}

// PathRev returns all keys in the path from this node to the root.
func (n *BasicNode[K, V]) PathRev() []K {
	var (
		keys = []K{n.key}
		cur  = n
	)

	for cur.parent != nil {
		cur = cur.parent
		keys = append(keys, cur.key)
	}

	return keys
}

// Value returns the node's value.
func (n *BasicNode[K, V]) Value() (value V) {
	if n == nil {
		return
	}
	return n.value
}

// SetValue sets the node's value.
func (n *BasicNode[K, V]) SetValue(value V) {
	if n == nil {
		return
	}
	n.value = value
}

// Parent returns the node's parent node.
func (n *BasicNode[K, V]) Parent() *BasicNode[K, V] {
	if n == nil {
		return nil
	}
	return n.parent
}

// Child returns the child of with the given key, if one exists. If no child
// with the given key is found, nil is returned.
func (n *BasicNode[K, V]) Child(key K) *BasicNode[K, V] {
	if n == nil {
		return nil
	}

	c, ok := n.children[key]
	if !ok {
		return nil
	}
	return c
}

// Children returns the node's children.
func (n *BasicNode[K, V]) Children() map[K]*BasicNode[K, V] {
	switch {
	case n == nil:
		return nil
	case len(n.children) == 0:
		return nil
	default:
		return maps.Clone(n.children)
	}
}

// SetParent sets the node's parent to the given parent node.
func (n *BasicNode[K, V]) SetParent(parent *BasicNode[K, V]) {
	if n == nil {
		return
	}

	if n.parent != nil {
		delete(n.parent.children, n.key)
	}

	if parent.children == nil {
		parent.children = make(map[K]*BasicNode[K, V])
	}
	parent.children[n.key] = n
	n.parent = parent
}

// Add adds a new child with the given key and value to this node, returning
// the new node.
func (n *BasicNode[K, V]) Add(key K, value V) *BasicNode[K, V] {
	node := &BasicNode[K, V]{
		key:    key,
		value:  value,
		parent: n,
	}

	if n != nil {
		if n.children == nil {
			n.children = make(map[K]*BasicNode[K, V])
		}
		n.children[key] = node
	}

	return node
}

// Remove removes the child with the given key, if one exists.
func (n *BasicNode[K, V]) Remove(
	key K,
) (child *BasicNode[K, V], removed bool) {
	if child, removed = n.children[key]; removed {
		child.parent = nil
		delete(n.children, key)
	}
	return
}

// Len returns the recursive length of the tree relative to the node.
func (n *BasicNode[K, V]) Len() (total int) {
	switch {
	case n == nil:
		return 0
	case len(n.children) == 0:
		return 1
	default:
		for _, child := range n.children {
			total += child.Len()
		}
		total++ // for n itself
		return
	}
}

func handleWalkError(err error) (stop bool, unhandled error) {
	switch {
	case err == nil:
		return false, nil
	case errors.Is(err, ErrSkipSubtree):
		return true, nil
	default:
		return true, err
	}
}

// Walk walks through the tree depth-first.
func (n *BasicNode[K, V]) Walk(fn NodeWalker[K, V]) error {
	if n == nil {
		return nil
	}

	stop, err := handleWalkError(fn(n))
	if stop || len(n.children) == 0 {
		return err
	}

	keys := xmaps.Keys(n.children)
	sort.Slice(keys, func(i int, j int) bool {
		return keys[i] < keys[j]
	})

	for _, key := range keys {
		if stop, err = handleWalkError(n.children[key].Walk(fn)); stop {
			return err
		}
	}

	return nil
}

// WalkRev walks through the tree depth-first in reverse.
func (n *BasicNode[K, V]) WalkRev(fn NodeWalker[K, V]) error {
	if n == nil {
		return nil
	}

	var (
		stop bool
		err  error
	)

	if len(n.children) > 0 {
		keys := xmaps.Keys(n.children)
		sort.Slice(keys, func(i int, j int) bool {
			return keys[i] < keys[j]
		})

		for _, key := range keys {
			if stop, err = handleWalkError(n.children[key].WalkRev(fn)); stop {
				return err
			}
		}
	}

	_, err = handleWalkError(fn(n))
	return err
}
