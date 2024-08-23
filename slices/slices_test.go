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

package slices_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/container/set"
	"go.mway.dev/x/slices"
	"go.mway.dev/x/unsafe"
)

func TestHasPrefix(t *testing.T) {
	cases := map[string]struct {
		giveSlice     []int
		givePrefix    []int
		wantHasPrefix bool
	}{
		"nil slice nil prefix": {
			giveSlice:     nil,
			givePrefix:    nil,
			wantHasPrefix: true,
		},
		"non-nil slice nil prefix": {
			giveSlice:     []int{0},
			givePrefix:    nil,
			wantHasPrefix: true,
		},
		"not prefixed": {
			giveSlice:     []int{0, 1, 2},
			givePrefix:    []int{1, 2},
			wantHasPrefix: false,
		},
		"prefixed": {
			giveSlice:     []int{0, 1, 2},
			givePrefix:    []int{0, 1},
			wantHasPrefix: true,
		},
		"equal": {
			giveSlice:     []int{0, 1, 2},
			givePrefix:    []int{0, 1, 2},
			wantHasPrefix: true,
		},
		"prefix too long": {
			giveSlice:     []int{0, 1},
			givePrefix:    []int{0, 1, 2},
			wantHasPrefix: false,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tt.wantHasPrefix,
				slices.HasPrefix(tt.giveSlice, tt.givePrefix),
			)
		})
	}
}

func TestFilter(t *testing.T) {
	var (
		allowAll   = func(int) bool { return true }
		allowNone  = func(int) bool { return false }
		allowOdds  = func(i int) bool { return i%2 != 0 }
		allowEvens = func(i int) bool { return !allowOdds(i) }
	)

	cases := map[string]struct {
		give []int
		pred func(int) bool
		want []int
	}{
		"nil": {
			give: nil,
			pred: allowAll,
			want: nil,
		},
		"empty": {
			give: []int{},
			pred: allowAll,
			want: nil,
		},
		"all": {
			give: []int{1, 2, 3, 4},
			pred: allowAll,
			want: []int{1, 2, 3, 4},
		},
		"none": {
			give: []int{1, 2, 3, 4},
			pred: allowNone,
			want: nil,
		},
		"odds": {
			give: []int{1, 2, 3, 4},
			pred: allowOdds,
			want: []int{1, 3},
		},
		"evens": {
			give: []int{1, 2, 3, 4},
			pred: allowEvens,
			want: []int{2, 4},
		},
		"only incrementing": {
			give: []int{1, 1, 0, 1, 2, 1, 4, 5, 2, 6, 5},
			pred: func() func(i int) bool {
				var highest int
				return func(i int) bool {
					if i > highest {
						highest = i
						return true
					}
					return false
				}
			}(),
			want: []int{1, 2, 4, 5, 6},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, slices.Filter(tt.give, tt.pred))
		})
	}
}

func TestTransform(t *testing.T) {
	cases := map[string]struct {
		give   []string
		mapper func(string) int
		want   []int
	}{
		"empty source": {
			give:   nil,
			mapper: func(string) int { return 0 },
			want:   nil,
		},
		"nil mapper": {
			give:   []string{"a", "ab", "abc"},
			mapper: nil,
			want:   nil,
		},
		"lengths": {
			give:   []string{"a", "ab", "abc"},
			mapper: func(str string) int { return len(str) },
			want:   []int{1, 2, 3},
		},
		"unique chars": {
			give: []string{"aa", "aabb", "aabbcc"},
			mapper: func(str string) int {
				x := set.New(unsafe.StringToBytes(str)...)
				return x.Len()
			},
			want: []int{1, 2, 3},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, slices.Transform(tt.give, tt.mapper))
		})
	}
}
