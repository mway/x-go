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

// Package channelstest provides channel-related test helpers.
package channelstest

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/channels"
	"go.mway.dev/x/testing"
)

// RequireSend fails t if it fails to send value to ch within timeout.
func RequireSend[T any](
	t testing.T,
	ch chan<- T,
	value T,
	timeout time.Duration,
) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	require.True(t, channels.Send(ctx, ch, value))
}

// RequireNoSend fails to if it sends value to ch within timeout.
func RequireNoSend[T any](
	t testing.T,
	ch chan<- T,
	value T,
	timeout time.Duration,
) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	require.False(t, channels.Send(ctx, ch, value))
}

// RequireRecv fails t if it receives from ch within timeout. Otherwise, it
// returns the received value.
func RequireRecv[T any](
	t testing.T,
	ch <-chan T,
	timeout time.Duration,
) T {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	value, ok := channels.Recv(ctx, ch)
	require.True(t, ok)
	return value
}

// RequireNoRecv fails t if it receives from ch within timeout.
func RequireNoRecv[T any](
	t testing.T,
	ch <-chan T,
	timeout time.Duration,
) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, ok := channels.Recv(ctx, ch)
	require.False(t, ok)
}
