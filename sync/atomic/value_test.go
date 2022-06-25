package atomic_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/ptr"
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
			give: ptr.Of(testStruct{t.Name()}),
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
