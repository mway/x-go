// Copyright (c) 2022 Matt Way
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

package ptr_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/ptr"
)

func TestOf(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		expect := int(123)
		require.Equal(t, expect, *ptr.Of(expect))
	})

	t.Run("uint", func(t *testing.T) {
		expect := uint(123)
		require.Equal(t, expect, *ptr.Of(expect))
	})

	t.Run("struct", func(t *testing.T) {
		expect := struct {
			value string
		}{
			value: t.Name(),
		}
		require.Equal(t, expect, *ptr.Of(expect))
	})

	t.Run("pointer", func(t *testing.T) {
		expect := bytes.NewBuffer(nil)
		require.Equal(t, expect, *ptr.Of(expect))
	})
}

func TestFrom(t *testing.T) {
	require.Equal(t, int(123), ptr.From(ptr.Of(int(123)), 0))
	require.Equal(t, int(123), ptr.From(nil, int(123)))
}
