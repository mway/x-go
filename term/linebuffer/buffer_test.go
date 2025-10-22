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

package linebuffer_test

import (
	"bytes"
	"context"
	"io"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.mway.dev/x/channels"
	"go.mway.dev/x/term/linebuffer"
)

func TestBuffer_Nominal(t *testing.T) {
	buf := linebuffer.New(3)
	for i := range 10 {
		var wantOverflow int
		if i >= buf.Cap() {
			wantOverflow++
		}
		require.Equal(t, wantOverflow, buf.Add(strconv.Itoa(i)))
		requireUpdate(t, context.Background(), buf)
	}

	require.Equal(t, 3, buf.Len())
	requireHasLines(t, []string{"7", "8", "9"}, buf)
}

func TestBuffer_Overflow(t *testing.T) {
	t.Run("burst", func(t *testing.T) {
		buf := linebuffer.New(3)
		require.Equal(t, 0, buf.Add("1", "2"))
		requireUpdate(t, context.Background(), buf)
		require.Equal(t, 1, buf.Add("3", "4"))
		requireUpdate(t, context.Background(), buf)
		requireHasLines(t, []string{"2", "3", "4"}, buf)
	})

	t.Run("overrun", func(t *testing.T) {
		var (
			buf   = linebuffer.New(3)
			tostr = func(i int) string {
				return strconv.Itoa(i)
			}
		)
		for i := range 10 {
			wantOverflow := 1
			if i > 0 {
				wantOverflow += 3
			}
			require.Equal(t, wantOverflow, buf.Add(
				tostr(i),
				tostr(i+1),
				tostr(i+2),
				tostr(i+3),
			))
			requireUpdate(t, context.Background(), buf)
		}
		require.Equal(t, 3, buf.Len())
		requireHasLines(t, []string{"10", "11", "12"}, buf)
	})
}

func TestBuffer_EmptyStrings(t *testing.T) {
	t.Run("all empty", func(t *testing.T) {
		buf := linebuffer.New(3)
		require.Equal(t, 0, buf.Add("", "", ""))
		require.Equal(t, 0, buf.Len())
		requireHasLines(t, nil, buf)
	})
	t.Run("mixed", func(t *testing.T) {
		buf := linebuffer.New(3)
		require.Equal(t, 0, buf.Add("1", "", "2", "", "3", ""))
		require.Equal(t, 3, buf.Len())
		requireUpdate(t, context.Background(), buf)
		requireHasLines(t, []string{"1", "2", "3"}, buf)
	})
}

func TestBuffer_ForEachLine(t *testing.T) {
	type pair struct {
		str string
		idx int
	}

	buf := linebuffer.New(3)
	buf.Add("1", "2", "3")

	t.Run("full", func(t *testing.T) {
		var seen []pair
		buf.ForEachLine(func(str string, idx int) bool {
			seen = append(seen, pair{
				str: str,
				idx: idx,
			})
			return true
		})

		want := []pair{
			{str: "3", idx: 2},
			{str: "2", idx: 1},
			{str: "1", idx: 0},
		}
		require.Equal(t, want, seen)
	})

	t.Run("partial", func(t *testing.T) {
		var seen []pair
		buf.ForEachLine(func(str string, idx int) bool {
			seen = append(seen, pair{
				str: str,
				idx: idx,
			})
			return false
		})

		want := []pair{
			{str: "3", idx: 2},
		}
		require.Equal(t, want, seen)
	})
}

func TestBuffer_WithSource(t *testing.T) {
	buf := linebuffer.New(
		3,
		linebuffer.WithSource(bytes.NewBufferString("1\n2\n3\n4\n")),
	)
	defer func() {
		require.NoError(t, buf.Close())
	}()

	timeout := time.NewTimer(time.Second)
	defer timeout.Stop()

	func() {
		for {
			select {
			case _, ok := <-buf.Updates():
				if !ok || buf.Len() >= 3 {
					return
				}
			case <-timeout.C:
				require.FailNow(t, "timed out waiting for updates")
			}
		}
	}()

	require.Equal(t, 3, buf.Len())
	requireUpdate(t, context.Background(), buf)
	requireHasLines(t, []string{"2", "3", "4"}, buf)
}

func TestBuffer_WithDelayedSource(t *testing.T) {
	var (
		input, output = io.Pipe()
		buf           = linebuffer.New(3, linebuffer.WithSource(input))
	)

	defer func() {
		require.NoError(t, buf.Close())
	}()

	const (
		width = 100 * time.Millisecond
		times = 4
	)
	for i := range 4 {
		time.AfterFunc((width+1)*time.Duration(i), func() {
			_, err := io.WriteString(output, strconv.Itoa(i+1)+"\n")
			require.NoError(t, err)
		})
	}
	time.AfterFunc((width+1)*time.Duration(5), func() {
		require.NoError(t, output.Close())
	})

	timeout := time.NewTimer(time.Second)
	defer timeout.Stop()

	func() {
		for {
			select {
			case _, ok := <-buf.Updates():
				if !ok || buf.Len() >= 3 {
					return
				}
			case <-timeout.C:
				require.FailNow(t, "timed out waiting for updates")
			}
		}
	}()

	require.Equal(t, 3, buf.Len())
	requireUpdate(t, context.Background(), buf)
	requireHasLines(t, []string{"2", "3", "4"}, buf)
}

func TestBuffer_SetCap(t *testing.T) {
	buf := linebuffer.New(3)
	buf.SetCap(3)
	require.Equal(t, 3, buf.Cap())
	require.Equal(t, 0, buf.Len())

	buf.SetCap(4)
	require.Equal(t, 4, buf.Cap())
	requireUpdate(t, context.Background(), buf)
	require.Equal(t, 0, buf.Add("1", "2", "3", "4"))
	require.Equal(t, 4, buf.Len())
	requireUpdate(t, context.Background(), buf)
	requireHasLines(t, []string{"1", "2", "3", "4"}, buf)

	buf.SetCap(2)
	require.Equal(t, 2, buf.Cap())
	require.Equal(t, 2, buf.Len())
	requireUpdate(t, context.Background(), buf)
	requireHasLines(t, []string{"3", "4"}, buf)
}

func requireHasLines(t *testing.T, want []string, buf *linebuffer.Buffer) {
	t.Helper()

	var (
		haveCopy = buf.LinesCopy()
		haveIter []string
	)

	require.Equal(t, want, haveCopy)

	buf.ForEachLine(func(line string, _ int) bool {
		haveIter = append(haveIter, line)
		return true
	})

	if want != nil {
		want = slices.Clone(want)
	}
	slices.Reverse(want)
	require.Equal(t, want, haveIter)
}

//nolint:revive
func requireUpdate(t *testing.T, ctx context.Context, buf *linebuffer.Buffer) {
	t.Helper()
	_, updated := channels.RecvWithTimeout(ctx, buf.Updates(), time.Second)
	require.True(t, updated)
}
