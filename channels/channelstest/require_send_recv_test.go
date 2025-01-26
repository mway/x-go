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

package channelstest_test

import (
	"testing"
	"time"

	"go.mway.dev/x/channels/channelstest"
	"go.mway.dev/x/testing/testingmock"
	"go.uber.org/mock/gomock"
)

func TestRequireSend(t *testing.T) {
	ch := make(chan struct{}, 1)
	defer safeClose(t, ch)

	mockT := testingmock.NewMockT(gomock.NewController(t))
	for range cap(ch) {
		channelstest.RequireSend(mockT, ch, struct{}{}, 100*time.Millisecond)
	}

	mockT.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
	mockT.EXPECT().FailNow().Times(1)
	channelstest.RequireSend(mockT, ch, struct{}{}, 100*time.Millisecond)
}

func TestRequireNoSend(t *testing.T) {
	ch := make(chan struct{})
	defer safeClose(t, ch)

	mockT := testingmock.NewMockT(gomock.NewController(t))
	for range 3 {
		channelstest.RequireNoSend(mockT, ch, struct{}{}, 100*time.Millisecond)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		<-ch
	}()

	mockT.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
	mockT.EXPECT().FailNow().Times(1)
	channelstest.RequireNoSend(mockT, ch, struct{}{}, 100*time.Millisecond)
	<-done
}

func TestRequireRecv(t *testing.T) {
	ch := make(chan struct{}, 1)
	defer safeClose(t, ch)

	ch <- struct{}{}
	mockT := testingmock.NewMockT(gomock.NewController(t))
	channelstest.RequireRecv(mockT, ch, 100*time.Millisecond)

	mockT.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
	mockT.EXPECT().FailNow().Times(1)
	channelstest.RequireRecv(mockT, ch, 100*time.Millisecond)
}

func TestRequireNoRecv(t *testing.T) {
	ch := make(chan struct{})
	defer safeClose(t, ch)

	mockT := testingmock.NewMockT(gomock.NewController(t))
	channelstest.RequireNoRecv(mockT, ch, 100*time.Millisecond)

	done := make(chan struct{})
	defer wait(done)
	go func() {
		defer close(done)
		ch <- struct{}{}
	}()

	mockT.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
	mockT.EXPECT().FailNow().Times(1)
	channelstest.RequireNoRecv(mockT, ch, 100*time.Millisecond)

	<-done
}

func safeClose[T any](t *testing.T, ch chan T) {
	t.Helper()
	select {
	case <-ch:
	default:
		close(ch)
	}
}

func wait(ch <-chan struct{}) {
	<-ch
}
