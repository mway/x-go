// Copyright (c) 2023 Matt Way
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

package env

import (
	"slices"
)

// A LookupOption configures an environment lookup through [Get], [GetAs],
// [GetAsP], [Lookup], [LookupAs], and [LookupAsP].
type LookupOption interface {
	apply(*lookupOptions)
}

// A LookupFunc is a function that looks up a given key (typically an
// environment variable name) and returns its value and whether it was found.
type LookupFunc func(key string) (value string, found bool)

func (f LookupFunc) apply(dst *lookupOptions) {
	if f != nil {
		dst.LookupFuncs = append(slices.Clone(dst.LookupFuncs), f)
	}
}

// SanitizeNames controls
type SanitizeNames bool

func (s SanitizeNames) apply(dst *lookupOptions) {
	dst.SanitizeNames = bool(s)
}

type lookupOptions struct {
	LookupFuncs   []LookupFunc
	SanitizeNames bool
}

func defaultLookupOptions() lookupOptions {
	return lookupOptions{}
}

func (o lookupOptions) With(opts ...LookupOption) lookupOptions {
	for _, opt := range opts {
		opt.apply(&o)
	}
	return o
}
