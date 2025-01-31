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
	"fmt"
	"net"
	"strconv"
)

var (
	// ErrParsingFailed is returned by parsing functions when parsing fails.
	ErrParsingFailed = errors.New("parsing failed")
	// ErrNoPortFound is returned by [ParsePort] when no port is found.
	ErrNoPortFound = fmt.Errorf("%w: no port found", ErrParsingFailed)
	// ErrMalformedAddr is returned by [ParsePort] when a malformed port is found.
	ErrMalformedAddr = fmt.Errorf("%w: malformed addr", ErrParsingFailed)
	// ErrInvalidPort is returned by [ParsePort] when an invalid port is found.
	ErrInvalidPort = fmt.Errorf("%w: invalid port", ErrParsingFailed)
)

// RandomPort returns a random, OS-assigned port bindable to the given addr.
func RandomPort(addr string) (int, error) {
	l, err := _listen("tcp", net.JoinHostPort(addr, "0"))
	if err != nil {
		return 0, err
	}

	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}

	return port, nil
}

// RandomLocalPort returns a random, OS-assigned port bindable to localhost.
func RandomLocalPort() (int, error) {
	return RandomPort(_localAddr)
}

// MustPort panics if the given error is not nil; otherwise, it returns the
// given port.
func MustPort(port int, err error) int {
	if err != nil {
		panic(err)
	}

	return port
}

// ParsePort parses any port value in addr, returning the numeric value or any
// error encountered while parsing.
func ParsePort(addr string) (int, error) {
	_, portstr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, fmt.Errorf("%w: %q: %w", ErrMalformedAddr, addr, err)
	}
	if len(portstr) == 0 {
		return 0, fmt.Errorf("%w: %q", ErrNoPortFound, addr)
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to parse port: %w: %w",
			ErrInvalidPort,
			err,
		)
	}

	return port, nil
}
