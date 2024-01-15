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

package io

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNop(t *testing.T) {
	var (
		x   nop
		b   []byte
		n   int
		err error
	)

	n, err = x.Read(make([]byte, 8))
	require.Equal(t, 0, n)
	require.ErrorIs(t, err, io.EOF)

	b, err = io.ReadAll(x)
	require.Len(t, b, 0)
	require.NoError(t, err)

	n, err = x.Write([]byte("hello"))
	require.Equal(t, 5, n)
	require.NoError(t, err)

	n, err = x.WriteString("hello")
	require.Equal(t, 5, n)
	require.NoError(t, err)

	require.NoError(t, x.Close())
	require.NoError(t, x.Close())
}
