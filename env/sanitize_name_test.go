// Copyright (c) 2023 Matt Way
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

package env_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.mway.dev/x/env"
)

func TestSanitizeName(t *testing.T) {
	// give -> want
	cases := map[string]string{
		"foo":     "FOO",
		"foo_bar": "FOO_BAR",
		"FOO_BAR": "FOO_BAR",
		"foo bar": "FOO_BAR",
		"foo$bar": "FOO_BAR",
		"__fOo__": "__FOO__",
		"  foo  ": "FOO",
		"$$FOO!!": "__FOO__",
		"0foo":    "_0FOO",
		"foo0":    "FOO0",
		"Foo":     "FOO",
	}

	for give, want := range cases {
		require.Equal(t, want, env.SanitizeName(give))
	}
}
