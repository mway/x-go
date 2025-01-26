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

package collections_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/collections"
)

func TestFirstLastOf(t *testing.T) {
	cases := map[string]struct {
		give      []int
		wantFirst int
		wantLast  int
	}{
		"nil": {
			give:      []int(nil),
			wantFirst: 0,
			wantLast:  0,
		},
		"empty": {
			give:      []int{},
			wantFirst: 0,
			wantLast:  0,
		},
		"single": {
			give:      []int{123},
			wantFirst: 123,
			wantLast:  123,
		},
		"negative": {
			give:      []int{-123},
			wantFirst: -123,
			wantLast:  -123,
		},
		"multiple": {
			give:      []int{123, 456},
			wantFirst: 123,
			wantLast:  456,
		},
		"nested": {
			give:      []int{0, 123, 0},
			wantFirst: 123,
			wantLast:  123,
		},
		"sparse": {
			give:      []int{0, 123, 0, 456, 0},
			wantFirst: 123,
			wantLast:  456,
		},
		"zeroes": {
			give:      []int{0, 0, 0},
			wantFirst: 0,
			wantLast:  0,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.wantFirst, collections.FirstOf(tt.give...))
			require.Equal(t, tt.wantLast, collections.LastOf(tt.give...))
		})
	}
}

func TestFirstLastOfFuncs(t *testing.T) {
	makeFuncs := func(values []int) []func() int {
		switch {
		case values == nil:
			return nil
		case len(values) == 0:
			return []func() int{}
		default:
			fns := make([]func() int, len(values))
			for i := range len(values) {
				fns[i] = func() int {
					return values[i]
				}
			}
			return fns
		}
	}

	cases := map[string]struct {
		give      []func() int
		wantFirst int
		wantLast  int
	}{
		"nil": {
			give:      makeFuncs(nil),
			wantFirst: 0,
			wantLast:  0,
		},
		"empty": {
			give:      makeFuncs([]int{}),
			wantFirst: 0,
			wantLast:  0,
		},
		"single": {
			give:      makeFuncs([]int{123}),
			wantFirst: 123,
			wantLast:  123,
		},
		"negative": {
			give:      makeFuncs([]int{-123}),
			wantFirst: -123,
			wantLast:  -123,
		},
		"multiple": {
			give:      makeFuncs([]int{123, 456}),
			wantFirst: 123,
			wantLast:  456,
		},
		"nested": {
			give:      makeFuncs([]int{0, 123, 0}),
			wantFirst: 123,
			wantLast:  123,
		},
		"sparse": {
			give:      makeFuncs([]int{0, 123, 0, 456, 0}),
			wantFirst: 123,
			wantLast:  456,
		},
		"zeroes": {
			give:      makeFuncs([]int{0, 0, 0}),
			wantFirst: 0,
			wantLast:  0,
		},
		"nil func": {
			give:      []func() int{nil},
			wantFirst: 0,
			wantLast:  0,
		},
		"sparse nil funcs": {
			give: []func() int{
				nil,
				func() int { return 123 },
				nil,
			},
			wantFirst: 123,
			wantLast:  123,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			haveFirst := collections.FirstOfFuncs(tt.give...)
			require.Equal(t, tt.wantFirst, haveFirst)
			haveLast := collections.LastOfFuncs(tt.give...)
			require.Equal(t, tt.wantLast, haveLast)
		})
	}
}

func TestFirstLastOfOr(t *testing.T) {
	cases := map[string]struct {
		fn        func() int
		give      []int
		wantFirst int
		wantLast  int
	}{
		"nil": {
			give:      nil,
			fn:        func() int { return 123 },
			wantFirst: 123,
			wantLast:  123,
		},
		"nil with nil func": {
			give:      nil,
			fn:        nil,
			wantFirst: 0,
			wantLast:  0,
		},
		"nil with func": {
			give:      nil,
			fn:        func() int { return 123 },
			wantFirst: 123,
			wantLast:  123,
		},
		"single": {
			give:      []int{123},
			fn:        nil,
			wantFirst: 123,
			wantLast:  123,
		},
		"negative": {
			give:      []int{-123},
			fn:        nil,
			wantFirst: -123,
			wantLast:  -123,
		},
		"multiple": {
			give:      []int{123, 456},
			fn:        nil,
			wantFirst: 123,
			wantLast:  456,
		},
		"nested": {
			give:      []int{0, 123, 0},
			fn:        nil,
			wantFirst: 123,
			wantLast:  123,
		},
		"sparse": {
			give:      []int{0, 123, 0, 456, 0},
			fn:        nil,
			wantFirst: 123,
			wantLast:  456,
		},
		"zeroes": {
			give:      []int{0, 0, 0},
			fn:        nil,
			wantFirst: 0,
			wantLast:  0,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			haveFirst := collections.FirstOfOr(tt.fn, tt.give...)
			require.Equal(t, tt.wantFirst, haveFirst)
			haveLast := collections.LastOfOr(tt.fn, tt.give...)
			require.Equal(t, tt.wantLast, haveLast)
		})
	}
}
