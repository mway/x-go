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

package channels_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.mway.dev/x/channels"
)

func TestDrain(t *testing.T) {
	ch := make(chan struct{}, 3)
	for range 3 {
		ch <- struct{}{}
	}
	require.Equal(t, 3, channels.Drain(ch))
}

func TestDrain_Closed(t *testing.T) {
	ch := make(chan struct{}, 3)
	for range 3 {
		ch <- struct{}{}
	}
	close(ch)
	require.Equal(t, 3, channels.Drain(ch))
}

func TestDrainContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		100*time.Millisecond,
	)
	defer cancel()

	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	require.Equal(t, 1, channels.DrainContext(ctx, ch))
}

func TestDrainContext_Closed(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		100*time.Millisecond,
	)
	defer cancel()

	ch := make(chan struct{})
	close(ch)
	require.Equal(t, 0, channels.DrainContext(ctx, ch))
}

func TestDrainContext_Canceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var ch chan int
	require.Equal(t, 0, channels.DrainContext(ctx, ch))
}
