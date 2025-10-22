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

package filepath_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.mway.dev/x/path/filepath"
)

func TestAncestors(t *testing.T) {
	cases := map[string]struct {
		givePath  string
		wantPaths []string
	}{
		"empty": {
			givePath:  "",
			wantPaths: nil,
		},
		"dot": {
			givePath:  ".",
			wantPaths: nil,
		},
		"dotdot": {
			givePath:  "..",
			wantPaths: nil,
		},
		"root": {
			givePath:  "/",
			wantPaths: nil,
		},
		"single": {
			givePath:  "foo",
			wantPaths: nil,
		},
		"absolute": {
			givePath: "/foo/bar/baz/bat",
			wantPaths: []string{
				"/foo/bar/baz",
				"/foo/bar",
				"/foo",
			},
		},
		"relative": {
			givePath: "foo/bar/baz/bat",
			wantPaths: []string{
				"foo/bar/baz",
				"foo/bar",
				"foo",
			},
		},
		"unresolved absolute": {
			givePath: "/foo/bar/../baz/bat",
			wantPaths: []string{
				"/foo/baz",
				"/foo",
			},
		},
		"unresolved relative": {
			givePath: "foo/bar/../baz/bat",
			wantPaths: []string{
				"foo/baz",
				"foo",
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.wantPaths, filepath.Ancestors(tt.givePath))
		})
	}
}
