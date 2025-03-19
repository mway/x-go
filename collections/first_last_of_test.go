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
	"go.mway.dev/x/slices"
)

func TestFirstOf(t *testing.T) {
	cases := map[string]struct {
		give      []int
		wantIndex int
		wantValue int
	}{
		"nil slice": {
			give:      nil,
			wantIndex: 0,
			wantValue: 0,
		},
		"empty slice": {
			give:      []int{},
			wantIndex: 0,
			wantValue: 0,
		},
		"single 1": {
			give:      []int{123},
			wantIndex: 0,
			wantValue: 123,
		},
		"single 2": {
			give:      []int{0},
			wantIndex: 0,
			wantValue: 0,
		},
		"double 1": {
			give:      []int{123, 456},
			wantIndex: 0,
			wantValue: 123,
		},
		"double 2": {
			give:      []int{0, 456},
			wantIndex: 1,
			wantValue: 456,
		},
		"double 3": {
			give:      []int{123, 0},
			wantIndex: 0,
			wantValue: 123,
		},
		"double 4": {
			give:      []int{0, 0},
			wantIndex: 0,
			wantValue: 0,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.wantValue, collections.FirstOf(tt.give...))
			require.Equal(t, tt.wantValue, collections.FirstOfSeq(slices.Iter(tt.give)))
			haveIndex, haveValue := collections.FirstOfSeq2(slices.Iter2(tt.give))
			require.Equal(t, tt.wantIndex, haveIndex)
			require.Equal(t, tt.wantValue, haveValue)
			require.Equal(t, tt.wantValue, collections.FirstOfFuncs(toFuncs(tt.give)...))
		})
	}
}

func TestFirstOfLastOfFuncs_Nil(t *testing.T) {
	const want = 123

	cases := map[string]struct {
		funcs     []func() int
		wantFirst int
		wantLast  int
	}{
		"no funcs": {
			funcs:     nil,
			wantFirst: 0,
			wantLast:  0,
		},
		"single nil": {
			funcs:     []func() int{nil},
			wantFirst: 0,
			wantLast:  0,
		},
		"single ok": {
			funcs: []func() int{
				func() int { return want },
			},
			wantFirst: want,
			wantLast:  want,
		},
		"double nil": {
			funcs:     []func() int{nil, nil},
			wantFirst: 0,
			wantLast:  0,
		},
		"double ok": {
			funcs: []func() int{
				func() int { return want },
				func() int { return want * 2 },
			},
			wantFirst: want,
			wantLast:  want * 2,
		},
		"double first nil": {
			funcs: []func() int{
				nil,
				func() int { return want * 2 },
			},
			wantFirst: want * 2,
			wantLast:  want * 2,
		},
		"double second nil": {
			funcs: []func() int{
				func() int { return want },
				nil,
			},
			wantFirst: want,
			wantLast:  want,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.wantFirst, collections.FirstOfFuncs(tt.funcs...))
			require.Equal(t, tt.wantLast, collections.LastOfFuncs(tt.funcs...))
		})
	}
}

func TestFirstOfOr(t *testing.T) {
	const (
		fallbackKey   = -1
		fallbackValue = -42
	)
	cases := map[string]struct {
		give      []int
		wantIndex int
		wantValue int
		fallback  bool
	}{
		"nil slice": {
			give:      nil,
			wantIndex: fallbackKey,
			wantValue: fallbackValue,
			fallback:  true,
		},
		"empty slice": {
			give:      []int{},
			wantIndex: fallbackKey,
			wantValue: fallbackValue,
			fallback:  true,
		},
		"single 1": {
			give:      []int{123},
			wantIndex: 0,
			wantValue: 123,
			fallback:  false,
		},
		"single 2": {
			give:      []int{0},
			wantIndex: fallbackKey,
			wantValue: fallbackValue,
			fallback:  true,
		},
		"double 1": {
			give:      []int{123, 456},
			wantIndex: 0,
			wantValue: 123,
			fallback:  false,
		},
		"double 2": {
			give:      []int{0, 456},
			wantIndex: 1,
			wantValue: 456,
			fallback:  false,
		},
		"double 3": {
			give:      []int{123, 0},
			wantIndex: 0,
			wantValue: 123,
			fallback:  false,
		},
		"double 4": {
			give:      []int{0, 0},
			wantIndex: fallbackKey,
			wantValue: fallbackValue,
			fallback:  true,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tt.wantValue,
				collections.FirstOfOr(fallbackValue, tt.give...),
			)
			require.Equal(
				t,
				tt.wantValue,
				collections.FirstOfSeqOr(fallbackValue, slices.Iter(tt.give)),
			)
			haveIndex, haveValue := collections.FirstOfSeq2Or(
				fallbackKey,
				fallbackValue,
				slices.Iter2(tt.give),
			)
			require.Equal(t, tt.wantIndex, haveIndex)
			require.Equal(t, tt.wantValue, haveValue)
			require.Equal(t, tt.wantValue, collections.FirstOfFuncsOr(
				fallbackValue,
				toFuncs(tt.give)...,
			))

			var (
				fallbackFunc  = toFunc(fallbackValue)
				fallbackFunc2 = func() (int, int) {
					return fallbackKey, fallbackValue
				}
			)

			if tt.fallback {
				require.Zero(t, collections.FirstOfOrElse(nil, tt.give...))
				require.Zero(t, collections.FirstOfSeqOrElse(nil, slices.Iter(tt.give)))
				haveIndex, haveValue = collections.FirstOfSeq2OrElse(nil, slices.Iter2(tt.give))
				require.Zero(t, haveIndex)
				require.Zero(t, haveValue)
				require.Zero(t, collections.FirstOfFuncsOrElse(
					nil,
					toFuncs(tt.give)...,
				))
			} else {
				require.Equal(t, tt.wantValue, collections.FirstOfOrElse(nil, tt.give...))
				require.Equal(t, tt.wantValue, collections.FirstOfSeqOrElse(nil, slices.Iter(tt.give)))
				haveIndex, haveValue = collections.FirstOfSeq2OrElse(
					nil,
					slices.Iter2(tt.give),
				)
				require.Equal(t, tt.wantIndex, haveIndex)
				require.Equal(t, tt.wantValue, haveValue)
				require.Equal(t, tt.wantValue, collections.FirstOfFuncsOrElse(
					nil,
					toFuncs(tt.give)...,
				))
			}

			require.Equal(
				t,
				tt.wantValue,
				collections.FirstOfOrElse(fallbackFunc, tt.give...),
			)
			require.Equal(
				t,
				tt.wantValue,
				collections.FirstOfSeqOrElse(fallbackFunc, slices.Iter(tt.give)),
			)
			haveIndex, haveValue = collections.FirstOfSeq2OrElse(
				fallbackFunc2,
				slices.Iter2(tt.give),
			)
			require.Equal(t, tt.wantIndex, haveIndex)
			require.Equal(t, tt.wantValue, haveValue)
			require.Equal(t, tt.wantValue, collections.FirstOfFuncsOrElse(
				fallbackFunc,
				toFuncs(tt.give)...,
			))
		})
	}
}

func TestLastOf(t *testing.T) {
	cases := map[string]struct {
		give      []int
		wantIndex int
		wantValue int
	}{
		"nil slice": {
			give:      nil,
			wantIndex: 0,
			wantValue: 0,
		},
		"empty slice": {
			give:      []int{},
			wantIndex: 0,
			wantValue: 0,
		},
		"single 1": {
			give:      []int{123},
			wantIndex: 0,
			wantValue: 123,
		},
		"single 2": {
			give:      []int{0},
			wantIndex: 0,
			wantValue: 0,
		},
		"double 1": {
			give:      []int{123, 456},
			wantIndex: 1,
			wantValue: 456,
		},
		"double 2": {
			give:      []int{0, 456},
			wantIndex: 1,
			wantValue: 456,
		},
		"double 3": {
			give:      []int{123, 0},
			wantIndex: 0,
			wantValue: 123,
		},
		"double 4": {
			give:      []int{0, 0},
			wantIndex: 0,
			wantValue: 0,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.wantValue, collections.LastOf(tt.give...))
			require.Equal(t, tt.wantValue, collections.LastOfSeq(slices.Iter(tt.give)))
			haveIndex, haveValue := collections.LastOfSeq2(slices.Iter2(tt.give))
			require.Equal(t, tt.wantIndex, haveIndex)
			require.Equal(t, tt.wantValue, haveValue)
			require.Equal(t, tt.wantValue, collections.LastOfFuncs(toFuncs(tt.give)...))
		})
	}
}

func TestLastOfOr(t *testing.T) {
	const (
		fallbackKey   = -1
		fallbackValue = -42
	)
	cases := map[string]struct {
		give      []int
		wantIndex int
		wantValue int
		fallback  bool
	}{
		"nil slice": {
			give:      nil,
			wantIndex: fallbackKey,
			wantValue: fallbackValue,
			fallback:  true,
		},
		"empty slice": {
			give:      []int{},
			wantIndex: fallbackKey,
			wantValue: fallbackValue,
			fallback:  true,
		},
		"single 1": {
			give:      []int{123},
			wantIndex: 0,
			wantValue: 123,
			fallback:  false,
		},
		"single 2": {
			give:      []int{0},
			wantIndex: fallbackKey,
			wantValue: fallbackValue,
			fallback:  true,
		},
		"double 1": {
			give:      []int{123, 456},
			wantIndex: 1,
			wantValue: 456,
			fallback:  false,
		},
		"double 2": {
			give:      []int{0, 456},
			wantIndex: 1,
			wantValue: 456,
			fallback:  false,
		},
		"double 3": {
			give:      []int{123, 0},
			wantIndex: 0,
			wantValue: 123,
			fallback:  false,
		},
		"double 4": {
			give:      []int{0, 0},
			wantIndex: fallbackKey,
			wantValue: fallbackValue,
			fallback:  true,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tt.wantValue,
				collections.LastOfOr(fallbackValue, tt.give...),
			)
			require.Equal(
				t,
				tt.wantValue,
				collections.LastOfSeqOr(fallbackValue, slices.Iter(tt.give)),
			)
			haveIndex, haveValue := collections.LastOfSeq2Or(
				fallbackKey,
				fallbackValue,
				slices.Iter2(tt.give),
			)
			require.Equal(t, tt.wantIndex, haveIndex)
			require.Equal(t, tt.wantValue, haveValue)
			require.Equal(t, tt.wantValue, collections.LastOfFuncsOr(
				fallbackValue,
				toFuncs(tt.give)...,
			))

			var (
				fallbackFunc  = toFunc(fallbackValue)
				fallbackFunc2 = func() (int, int) {
					return fallbackKey, fallbackValue
				}
			)

			if tt.fallback {
				require.Zero(t, collections.LastOfOrElse(nil, tt.give...))
				require.Zero(t, collections.LastOfSeqOrElse(nil, slices.Iter(tt.give)))
				haveIndex, haveValue = collections.LastOfSeq2OrElse(nil, slices.Iter2(tt.give))
				require.Zero(t, haveIndex)
				require.Zero(t, haveValue)
				require.Zero(t, collections.LastOfFuncsOrElse(
					nil,
					toFuncs(tt.give)...,
				))
			} else {
				require.Equal(t, tt.wantValue, collections.LastOfOrElse(nil, tt.give...))
				require.Equal(t, tt.wantValue, collections.LastOfSeqOrElse(nil, slices.Iter(tt.give)))
				haveIndex, haveValue = collections.LastOfSeq2OrElse(
					nil,
					slices.Iter2(tt.give),
				)
				require.Equal(t, tt.wantIndex, haveIndex)
				require.Equal(t, tt.wantValue, haveValue)
				require.Equal(t, tt.wantValue, collections.LastOfFuncsOrElse(
					nil,
					toFuncs(tt.give)...,
				))
			}

			require.Equal(
				t,
				tt.wantValue,
				collections.LastOfOrElse(fallbackFunc, tt.give...),
			)
			require.Equal(
				t,
				tt.wantValue,
				collections.LastOfSeqOrElse(fallbackFunc, slices.Iter(tt.give)),
			)
			haveIndex, haveValue = collections.LastOfSeq2OrElse(
				fallbackFunc2,
				slices.Iter2(tt.give),
			)
			require.Equal(t, tt.wantIndex, haveIndex)
			require.Equal(t, tt.wantValue, haveValue)
			require.Equal(t, tt.wantValue, collections.LastOfFuncsOrElse(
				fallbackFunc,
				toFuncs(tt.give)...,
			))
		})
	}
}

func toFunc[T any](value T) func() T {
	return func() T {
		return value
	}
}

func toFuncs[T any, S ~[]T](x S) []func() T {
	fns := make([]func() T, len(x))
	for i, val := range x {
		fns[i] = toFunc(val)
	}
	return fns
}
