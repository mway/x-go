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

package net

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

type listenFunc = func(string, string) (net.Listener, error)

var (
	errListenSentinel = errors.New("listen error")

	_ listenFunc = net.Listen
)

func TestListenRandom(t *testing.T) {
	var (
		cases = map[string]struct {
			listen      listenFunc
			expectError error
		}{
			"nominal": {
				listen:      nil,
				expectError: nil,
			},
			"error": {
				listen: func(string, string) (net.Listener, error) {
					return nil, errListenSentinel
				},
				expectError: errListenSentinel,
			},
		}
		validate = func(t *testing.T, expect error, l net.Listener, err error) {
			require.ErrorIs(t, err, expect)

			if expect == nil {
				require.NotNil(t, l)
			} else {
				require.Nil(t, l)
			}
		}
	)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Run("ListenRandom", func(t *testing.T) {
				withListen(tt.listen, func() {
					listener, err := ListenRandom(_localAddr)
					validate(t, tt.expectError, listener, err)
				})
			})

			t.Run("ListenLocalRandom", func(t *testing.T) {
				withListen(tt.listen, func() {
					listener, err := ListenLocalRandom()
					validate(t, tt.expectError, listener, err)
				})
			})
		})
	}
}

func TestMustListen(t *testing.T) {
	require.Panics(t, func() {
		MustListen(nil, errors.New("error"))
	})

	require.NotPanics(t, func() {
		MustListen(nopListener{}, nil)
	})
}

func withListen(listen listenFunc, fn func()) {
	if listen != nil {
		prev := _listen
		defer func() {
			_listen = prev
		}()
		_listen = listen
	}

	fn()
}

type nopListener struct{}

func (nopListener) Accept() (net.Conn, error) { return nil, nil } //nolint:nilnil
func (nopListener) Close() error              { return nil }
func (nopListener) Addr() net.Addr            { return nil }
