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

func TestFirstOf(t *testing.T) {
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

// // FirstOf returns the first T in values that is not a zero value of T.
// func FirstOf[T comparable](values ...T) T {
// 	var zero T
// 	for _, value := range values {
// 		if value != zero {
// 			return value
// 		}
// 	}
// 	return zero
// }

// // LastOf returns the last T in values that is not a zero value of T.
// func LastOf[T comparable](values ...T) T {
// 	var zero T
// 	for i := len(values) - 1; i >= 0; i-- {
// 		if value := values[i]; value != zero {
// 			return value
// 		}
// 	}
// 	return zero
// }
