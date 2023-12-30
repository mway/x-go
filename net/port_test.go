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

var errCloseSentinel = errors.New("close error")

func TestRandomPorts(t *testing.T) {
	var (
		cases = map[string]struct {
			listen      listenFunc
			expectError error
		}{
			"nominal": {
				listen:      nil,
				expectError: nil,
			},
			"listen error": {
				listen:      newUncloseableListener(errListenSentinel, true).Listen,
				expectError: errListenSentinel,
			},
			"close error": {
				listen:      newUncloseableListener(errCloseSentinel, false).Listen,
				expectError: errCloseSentinel,
			},
		}
		validate = func(t *testing.T, expect error, port int, err error) {
			require.ErrorIs(t, err, expect)

			if expect == nil {
				require.NotEqual(t, 0, port)
			} else {
				require.Equal(t, 0, port)
			}
		}
	)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Run("RandomPort", func(t *testing.T) {
				withListen(tt.listen, func() {
					port, err := RandomPort(_localAddr)
					validate(t, tt.expectError, port, err)
				})
			})

			t.Run("RandomLocalPort", func(t *testing.T) {
				withListen(tt.listen, func() {
					port, err := RandomLocalPort()
					validate(t, tt.expectError, port, err)
				})
			})
		})
	}
}

func TestMustPort(t *testing.T) {
	require.Panics(t, func() {
		MustPort(0, errors.New("error"))
	})

	require.NotPanics(t, func() {
		MustPort(123, nil)
	})
}

type uncloseableListener struct {
	addr          *net.TCPAddr
	err           error
	errorOnListen bool
}

func newUncloseableListener(
	err error,
	errorOnListen bool,
) uncloseableListener {
	return uncloseableListener{
		err:           err,
		errorOnListen: errorOnListen,
	}
}

func (l uncloseableListener) Listen(
	network string,
	hostport string,
) (net.Listener, error) {
	if l.errorOnListen && l.err != nil {
		return nil, l.err
	}

	x, err := net.ResolveTCPAddr(network, hostport)
	if err != nil {
		return nil, err
	}

	l.addr = x
	return l, nil
}

func (l uncloseableListener) Accept() (net.Conn, error) {
	return nil, nil
}

func (l uncloseableListener) Close() error {
	return l.err
}

func (l uncloseableListener) Addr() net.Addr {
	return l.addr
}
