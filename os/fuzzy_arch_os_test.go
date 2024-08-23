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

package os_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/os"
)

func TestFuzzyArch_Names(t *testing.T) {
	cases := map[string]struct {
		give string
		want []string
	}{
		"empty": {
			give: "",
			want: nil,
		},
		"unknown": {
			give: "unknown",
			want: []string{"unknown"},
		},
		"amd64": {
			give: "amd64",
			want: []string{"amd64", "x64", "x86_64"},
		},
		"x86_64": {
			give: "x86_64",
			want: []string{"amd64", "x64", "x86_64"},
		},
		"x64": {
			give: "x64",
			want: []string{"amd64", "x64", "x86_64"},
		},
		"arm64": {
			give: "arm64",
			want: []string{"aarch64", "arm64"},
		},
		"aarch64": {
			give: "aarch64",
			want: []string{"aarch64", "arm64"},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				os.FuzzyArch(strings.ToUpper(tt.give)).Names(),
			)
		})
	}
}

func TestFuzzyArch_Matches(t *testing.T) {
	cases := map[string]struct {
		arch os.FuzzyArch
		give string
		want bool
	}{
		"empty arch": {
			arch: os.FuzzyArch(""),
			give: "foo bar baz",
			want: false,
		},
		"glob arch": {
			arch: os.FuzzyArch("*"),
			give: "foo bar baz",
			want: true,
		},
		"prefix match": {
			arch: os.FuzzyArch("foo"),
			give: "foo bar baz",
			want: true,
		},
		"nested match": {
			arch: os.FuzzyArch("bar"),
			give: "foo bar baz",
			want: true,
		},
		"suffix match": {
			arch: os.FuzzyArch("baz"),
			give: "foo bar baz",
			want: true,
		},
		"case insensitive match": {
			arch: os.FuzzyArch("BaZ"),
			give: "foo bar bAz",
			want: true,
		},
		"exact match": {
			arch: os.FuzzyArch("foo bar baz"),
			give: "foo bar baz",
			want: true,
		},
		"partial glob": {
			arch: os.FuzzyArch("foo*"),
			give: "foo bar baz",
			want: false,
		},
		"amd64 alias": {
			arch: os.FuzzyArch("amd64"),
			give: "foo x86_64 baz",
			want: true,
		},
		"arm64 alias": {
			arch: os.FuzzyArch("arm64"),
			give: "foo aarch64 baz",
			want: true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.arch.Matches(tt.give))
		})
	}
}

func TestFuzzyOS_Names(t *testing.T) {
	cases := map[string]struct {
		give string
		want []string
	}{
		"empty": {
			give: "",
			want: nil,
		},
		"unknown": {
			give: "unknown",
			want: []string{"unknown"},
		},
		"darwin": {
			give: "darwin",
			want: []string{"darwin", "macos"},
		},
		"macos": {
			give: "macos",
			want: []string{"darwin", "macos"},
		},
		"debian": {
			give: "debian",
			want: []string{"debian", "linux", "ubuntu"},
		},
		"linux": {
			give: "linux",
			want: []string{"debian", "linux", "ubuntu"},
		},
		"ubuntu": {
			give: "ubuntu",
			want: []string{"debian", "linux", "ubuntu"},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, os.FuzzyOS(tt.give).Names())
		})
	}
}

func TestFuzzyOS_Matches(t *testing.T) {
	cases := map[string]struct {
		arch os.FuzzyOS
		give string
		want bool
	}{
		"empty os": {
			arch: os.FuzzyOS(""),
			give: "foo bar baz",
			want: false,
		},
		"glob os": {
			arch: os.FuzzyOS("*"),
			give: "foo bar baz",
			want: true,
		},
		"prefix match": {
			arch: os.FuzzyOS("foo"),
			give: "foo bar baz",
			want: true,
		},
		"nested match": {
			arch: os.FuzzyOS("bar"),
			give: "foo bar baz",
			want: true,
		},
		"suffix match": {
			arch: os.FuzzyOS("baz"),
			give: "foo bar baz",
			want: true,
		},
		"case insensitive match": {
			arch: os.FuzzyOS("BaZ"),
			give: "foo bar bAz",
			want: true,
		},
		"exact match": {
			arch: os.FuzzyOS("foo bar baz"),
			give: "foo bar baz",
			want: true,
		},
		"partial glob": {
			arch: os.FuzzyOS("foo*"),
			give: "foo bar baz",
			want: false,
		},
		"darwin alias": {
			arch: os.FuzzyOS("darwin"),
			give: "foo MacOS baz",
			want: true,
		},
		"debian alias": {
			arch: os.FuzzyOS("debian"),
			give: "foo Ubuntu baz",
			want: true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.arch.Matches(tt.give))
		})
	}
}
