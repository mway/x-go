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

package tree_test

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/container/tree"
)

func TestBasicNode(t *testing.T) {
	var zero tree.BasicNode[string, string]
	require.Zero(t, zero.Key())
	require.Zero(t, zero.Value())
	zero.SetValue(t.Name())
	require.Equal(t, t.Name(), zero.Value())

	require.NotPanics(t, func() {
		var node *tree.BasicNode[string, string]
		node.SetParent(nil)
		node.SetValue("foo")
		require.Equal(t, 0, node.Len())
	})

	root := tree.NewBasicNode(t.Name(), t.Name())
	require.Equal(t, t.Name(), root.Key())
	require.Equal(t, t.Name(), root.Value())
	require.Nil(t, root.Parent())
	require.Nil(t, root.Children())
	require.Nil(t, root.Child("child1"))
	require.Equal(t, 1, root.Len())

	child1 := root.Add("child1", "child1")
	require.Equal(t, root, child1.Parent())
	require.Equal(t, child1, root.Child("child1"))
	child1.Add("grandchild11", "grandchild11")
	child1.Add("grandchild12", "grandchild12")
	require.Len(t, child1.Children(), 2)
	require.Equal(t, 3, child1.Len())
	require.Equal(t, 4, root.Len())

	child2 := root.Add("child2", "child2")
	require.Equal(t, root, child2.Parent())
	require.Equal(t, child2, root.Child("child2"))
	require.Len(t, root.Children(), 2)
	child2.Add("grandchild21", "grandchild21")
	child2.Add("grandchild22", "grandchild22")
	require.Len(t, child2.Children(), 2)
	require.Equal(t, 3, child2.Len())
	require.Equal(t, 7, root.Len())

	var (
		wantKeys = []string{
			t.Name(),
			"child1",
			"grandchild11",
			"grandchild12",
			"child2",
			"grandchild21",
			"grandchild22",
		}
		haveKeys []string
	)
	err := root.Walk(func(node *tree.BasicNode[string, string]) error {
		haveKeys = append(haveKeys, node.Key())
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, wantKeys, haveKeys)

	var (
		wantRevKeys = []string{
			"grandchild11",
			"grandchild12",
			"child1",
			"grandchild21",
			"grandchild22",
			"child2",
			t.Name(),
		}
		haveRevKeys []string
	)
	err = root.WalkRev(func(node *tree.BasicNode[string, string]) error {
		haveRevKeys = append(haveRevKeys, node.Key())
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, wantRevKeys, haveRevKeys)

	var calls int
	err = root.Walk(func(*tree.BasicNode[string, string]) error {
		calls++
		return io.EOF
	})
	require.ErrorIs(t, err, io.EOF)
	err = root.WalkRev(func(*tree.BasicNode[string, string]) error {
		calls++
		return io.EOF
	})
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, 2, calls)
}

func TestBasicNode_EmptyOrNil(t *testing.T) {
	cases := map[string]struct {
		node      *tree.BasicNode[string, string]
		wantCalls int
	}{
		"nil": {
			node:      nil,
			wantCalls: 0,
		},
		"empty": {
			node:      &tree.BasicNode[string, string]{},
			wantCalls: 2,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.NotPanics(t, func() {
				require.Zero(t, tt.node.Key())
				require.Zero(t, tt.node.Value())
				require.Nil(t, tt.node.Parent())
				require.Nil(t, tt.node.Children())
				require.Nil(t, tt.node.Child("foo"))
				require.Equal(t, tt.node == nil, tt.node.Len() == 0)

				child := tt.node.Add("foo", "bar")
				require.NotNil(t, child)
				require.Equal(t, tt.node, child.Parent())
				require.Equal(t, 1, child.Len())

				var calls int
				require.NoError(t, tt.node.Walk(
					func(*tree.BasicNode[string, string]) error {
						calls++
						return nil
					},
				))
				require.NoError(t, tt.node.WalkRev(
					func(*tree.BasicNode[string, string]) error {
						calls++
						return nil
					},
				))
				require.Equal(t, tt.wantCalls*2, calls)
			})
		})
	}
}

func TestBasicNode_Path(t *testing.T) {
	a := tree.NewBasicNode("a", "a")
	b := a.Add("b", "b")
	b.Add("c", "c")

	var (
		expect = [][]string{
			{"a"},
			{"a", "b"},
			{"a", "b", "c"},
		}
		expectRev = [][]string{
			{"a"},
			{"b", "a"},
			{"c", "b", "a"},
		}
	)

	var i int
	a.Walk(func(n *tree.BasicNode[string, string]) error { //nolint:errcheck
		require.Equal(t, expect[i], n.Path())
		require.Equal(t, expectRev[i], n.PathRev())
		i++
		return nil
	})
}

func TestBasicNode_SetParent(t *testing.T) {
	a := tree.NewBasicNode("a", "b")
	b := tree.NewBasicNode("b", "b")
	c := tree.NewBasicNode("c", "c")
	c.SetParent(b)
	c.SetParent(a)

	var (
		expect = map[string]int{
			"a": 1,
			"c": 1,
		}
		actual = make(map[string]int)
	)

	a.Walk(func(n *tree.BasicNode[string, string]) error { //nolint:errcheck
		actual[n.Key()]++
		return nil
	})

	require.Equal(t, expect, actual)
}

func TestBasicNode_Remove(t *testing.T) {
	a := tree.NewBasicNode("a", "b")
	b := tree.NewBasicNode("b", "b")
	c := tree.NewBasicNode("c", "c")
	c.SetParent(b)

	removedC, ok := b.Remove(c.Key())
	require.True(t, ok)
	require.Equal(t, c, removedC)

	var (
		expect = map[string]int{
			"a": 1,
		}
		actual = make(map[string]int)
	)

	a.Walk(func(n *tree.BasicNode[string, string]) error { //nolint:errcheck
		actual[n.Key()]++
		return nil
	})

	require.Equal(t, expect, actual)
}

func TestBasicNode_WalkErrors(t *testing.T) {
	root := tree.NewBasicNode(t.Name(), t.Name()+"value")
	require.Equal(t, t.Name(), root.Key())
	require.Equal(t, t.Name()+"value", root.Value())
	require.Nil(t, root.Parent())
	require.Nil(t, root.Children())
	require.Nil(t, root.Child("child1"))
	require.Equal(t, 1, root.Len())

	child1 := root.Add("child1", "value")
	require.Equal(t, root, child1.Parent())
	require.Equal(t, child1, root.Child("child1"))
	child1.Add("grandchild11", "value")
	child1.Add("grandchild12", "value")
	require.Len(t, child1.Children(), 2)
	require.Equal(t, 3, child1.Len())
	require.Equal(t, 4, root.Len())

	child2 := root.Add("child2", "childv2")
	require.Equal(t, root, child2.Parent())
	require.Equal(t, child2, root.Child("child2"))
	require.Len(t, root.Children(), 2)
	child2.Add("grandchild21", "value")
	child2.Add("grandchild22", "value")
	require.Len(t, child2.Children(), 2)
	require.Equal(t, 3, child2.Len())
	require.Equal(t, 7, root.Len())

	rootKey := t.Name()

	cases := map[string]struct { //nolint:govet
		pred    func(key string) error
		want    []string
		wantRev []string
		wantErr error
	}{
		"abort at child2": {
			pred: func(key string) error {
				if key == "child2" {
					return io.EOF
				}
				return nil
			},
			want: []string{
				rootKey,
				"child1",
				"grandchild11",
				"grandchild12",
			},
			wantRev: []string{
				"grandchild11",
				"grandchild12",
				"child1",
				"grandchild21",
				"grandchild22",
			},
			wantErr: io.EOF,
		},
		"no grandchildren": {
			pred: func(key string) error {
				if strings.HasPrefix(key, "grand") {
					return tree.ErrSkipSubtree
				}
				return nil
			},
			want: []string{
				rootKey,
				"child1",
				"child2",
			},
			wantRev: []string{
				"child1",
				"child2",
				rootKey,
			},
			wantErr: nil,
		},
		"root only": {
			pred: func(key string) error {
				if key != rootKey {
					return tree.ErrSkipSubtree
				}
				return nil
			},
			want: []string{
				rootKey,
			},
			wantRev: []string{
				rootKey,
			},
			wantErr: nil,
		},
		"no child1": {
			pred: func(key string) error {
				if key == "child1" {
					return tree.ErrSkipSubtree
				}
				return nil
			},
			want: []string{
				rootKey,
				"child2",
				"grandchild21",
				"grandchild22",
			},
			wantRev: []string{
				"grandchild11",
				"grandchild12",
				"grandchild21",
				"grandchild22",
				"child2",
				rootKey,
			},
			wantErr: nil,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				seen   []string
				walker = func(node *tree.BasicNode[string, string]) (err error) {
					if err = tt.pred(node.Key()); err == nil {
						seen = append(seen, node.Key())
					}
					return err
				}
			)

			err := root.Walk(walker)
			require.Equal(t, tt.want, seen)
			require.ErrorIs(t, err, tt.wantErr)

			seen = seen[:0]
			err = root.WalkRev(walker)
			require.Equal(t, tt.wantRev, seen)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
