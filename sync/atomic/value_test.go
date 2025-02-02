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

package atomic_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/container/ptr"
	"go.mway.dev/x/sync/atomic"
)

func TestValue_Load(t *testing.T) {
	cases := map[string]struct {
		give     any
		validate func(any) bool
	}{
		"string": {
			give: "hello world",
			validate: func(x any) bool {
				_, ok := x.(string)
				return ok
			},
		},
		"bytes.Buffer": {
			give: bytes.NewBuffer(nil),
			validate: func(x any) bool {
				_, ok := x.(*bytes.Buffer)
				return ok
			},
		},
		"struct pointer": {
			give: ptr.To(testStruct{t.Name()}),
			validate: func(x any) bool {
				_, ok := x.(*testStruct)
				return ok
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				value  = atomic.NewValue(tt.give)
				actual = value.Load()
			)

			require.Equal(t, tt.give, actual)
			require.True(t, tt.validate(actual))
		})
	}
}

func TestValue_CompareAndSwap(t *testing.T) {
	var (
		structA = &testStruct{"hello"}
		structB = &testStruct{123}
		structC = &testStruct{}
		structD = &testStruct{}
		value   = atomic.NewValue(structA)
	)

	require.False(t, value.CompareAndSwap(structB, structA))
	require.False(t, value.CompareAndSwap(structB, structB))
	require.True(t, value.CompareAndSwap(structA, structA))
	require.True(t, value.CompareAndSwap(structA, structB))
	require.True(t, value.CompareAndSwap(structB, structC))
	require.False(t, value.CompareAndSwap(structD, structC))
	require.True(t, value.CompareAndSwap(structC, structD))
}

func TestValue_Store(t *testing.T) {
	value := atomic.NewValue(0)

	for i := 0; i < 100; i++ {
		value.Store(i)
		require.Equal(t, i, value.Load())
	}
}

func TestValue_Swap(t *testing.T) {
	value := atomic.NewValue(0)

	for i := 1; i < 100; i++ {
		require.Equal(t, i-1, value.Swap(i))
	}
}

type testStruct struct {
	val any
}
