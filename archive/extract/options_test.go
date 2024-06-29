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

package extract_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/archive/extract"
)

func TestOptions(t *testing.T) {
	requireEqualOptions := func(t *testing.T, want extract.Options, have extract.Options) {
		t.Helper()
		require.Equal(t, have.Callback == nil, want.Callback == nil)
		require.Equal(t, want.Output, have.Output)
		require.Equal(t, want.StripPrefix, have.StripPrefix)
		require.Equal(t, want.IncludePaths, have.IncludePaths)
		require.Equal(t, want.ExcludePaths, have.ExcludePaths)
		require.Equal(t, want.Delete, have.Delete)
	}

	t.Run("full", func(t *testing.T) {
		want := extract.Options{
			Callback: func(_ context.Context, _ string) error {
				return nil
			},
			Output:      bytes.NewBuffer(nil),
			StripPrefix: "abc123",
			IncludePaths: map[string]string{
				"foo": "bar",
			},
			ExcludePaths: []string{
				"baz",
			},
			Delete: true,
		}
		requireEqualOptions(t, want, extract.Options{}.With(want))
	})

	t.Run("with", func(t *testing.T) {
		var (
			writer = bytes.NewBuffer(nil)
			cases  = map[string]struct {
				give extract.Option
				want extract.Options
			}{
				"output": {
					give: extract.Output(writer),
					want: extract.Options{
						Output: writer,
					},
				},
				"strip prefix": {
					give: extract.StripPrefix(t.Name()),
					want: extract.Options{
						StripPrefix: t.Name(),
					},
				},
				"include paths": {
					give: extract.IncludePaths(map[string]string{"a": "b"}),
					want: extract.Options{
						IncludePaths: map[string]string{"a": "b"},
					},
				},
				"exclude paths": {
					give: extract.ExcludePaths([]string{t.Name()}),
					want: extract.Options{
						ExcludePaths: []string{t.Name()},
					},
				},
				"delete paths": {
					give: extract.Delete(true),
					want: extract.Options{
						Delete: true,
					},
				},
			}
		)

		for name, tt := range cases {
			t.Run(name, func(t *testing.T) {
				have := extract.Options{}.With(tt.give)
				require.Equal(t, tt.want, have)
			})
		}
	})
}
