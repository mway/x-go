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
	"go.mway.dev/x/slices"
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
