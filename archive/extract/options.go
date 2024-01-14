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

package extract

import (
	"context"
	"io"
	"maps"
	"slices"
)

var (
	_ Option = Callback(nil)
	_ Option = Options{}
)

// A Callback is a function called by [Extract] once an archive has been
// extracted. When called, the current working directory is the extract
// destination, and is also passed as the path parameter.
type Callback func(ctx context.Context, path string) error

func (c Callback) apply(dst *Options) {
	dst.Callback = c
}

// Options configure the behavior of [Extract].
type Options struct {
	Callback     Callback
	Output       io.Writer
	StripPrefix  string
	IncludePaths map[string]string
	ExcludePaths []string
	Delete       bool
}

// With returns a new [Options] with opts merged on top of o.
func (o Options) With(opts ...Option) Options {
	for _, opt := range opts {
		opt.apply(&o)
	}
	return o
}

func (o Options) apply(dst *Options) {
	if o.Callback != nil {
		dst.Callback = o.Callback
	}

	if o.Output != nil {
		dst.Output = o.Output
	}

	if len(o.StripPrefix) > 0 {
		dst.StripPrefix = o.StripPrefix
	}

	if len(o.IncludePaths) > 0 {
		dst.IncludePaths = maps.Clone(o.IncludePaths)
	}

	if len(o.ExcludePaths) > 0 {
		dst.ExcludePaths = slices.Clone(o.ExcludePaths)
	}

	if o.Delete {
		dst.Delete = true
	}
}

// An Option configures the behavior of [Extract].
type Option interface {
	apply(*Options)
}

// Output returns a new [Option] that configures [Extract] to write any output
// to the given writer.
func Output(output io.Writer) Option {
	return optionFunc(func(dst *Options) {
		dst.Output = output
	})
}

// StripPrefix returns a new [Option] that configures [Extract] to root all
// extraction at the given prefix. Note that paths that are not descendants of
// the given prefix will not be extracted.
func StripPrefix(prefix string) Option {
	return optionFunc(func(dst *Options) {
		dst.StripPrefix = prefix
	})
}

// IncludePaths returns a new [Option] that configures [Extract] to only
// consider the given paths for extraction. The given map's keys should be
// relative to either the archive root or the stripped prefix, and may be
// globs; the map's values are optional and specify non-default destinations
// for any path(s) matched by the corresponding key.
func IncludePaths(paths map[string]string) Option {
	return optionFunc(func(dst *Options) {
		dst.IncludePaths = maps.Clone(paths)
	})
}

// ExcludePaths returns a new [Option] that configures [Extract] to exclude any
// matching paths from extraction. The given paths should be relative to either
// the archive root or the stripped prefix, and may be globs.
func ExcludePaths(paths []string) Option {
	return optionFunc(func(dst *Options) {
		dst.ExcludePaths = slices.Clone(paths)
	})
}

// Delete returns a new [Option] that configures [Extract] to delete any
// archive directories from the destination before extracting them.
func Delete(del bool) Option {
	return optionFunc(func(dst *Options) {
		dst.Delete = del
	})
}

type optionFunc func(*Options)

func (f optionFunc) apply(dst *Options) {
	f(dst)
}
