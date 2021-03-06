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
	"net"
)

const _localAddr = "127.0.0.1"

var _listen = net.Listen

// ListenRandom returns a listener bound to a random port on the given addr.
func ListenRandom(addr string) (net.Listener, error) {
	l, err := _listen("tcp", net.JoinHostPort(addr, "0"))
	if err != nil {
		return nil, err
	}

	return l, nil
}

// ListenLocalRandom returns a listener bound to a random port on localhost.
func ListenLocalRandom() (net.Listener, error) {
	return ListenRandom(_localAddr)
}

// MustListen panics if the given error is not nil; otherwise, it returns the
// given listener.
func MustListen(listener net.Listener, err error) net.Listener {
	if err != nil {
		panic(err)
	}

	return listener
}
