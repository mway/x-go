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

package context

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mway.dev/chrono/clock"
	"go.uber.org/multierr"
)

var (
	// ErrNilClock indicates that a nil clock was provided.
	ErrNilClock = errors.New("nil clock provided")
	// ErrNilContext indicates that a nil parent context was provided.
	ErrNilContext = errors.New("nil context provided")
	// ErrInvalidDebounce indicates that the debounce value is invalid.
	ErrInvalidDebounce = errors.New("invalid debounce value provided")
	// ErrNilContextFunc indicates that the given TimeoutFunc is nil.
	ErrNilContextFunc = errors.New("nil context func provided")

	_defaultDebouncedOptions = DebouncedOptions{
		Clock:       clock.NewMonotonicClock(),
		Context:     context.Background(),
		Debounce:    time.Second,
		ContextFunc: context.WithTimeout,
	}

	_ DebouncedOption = DebouncedOptions{}
)

// A TimeoutFunc is a function that is analogous to context.WithTimeout.
type TimeoutFunc = func(
	context.Context,
	time.Duration,
) (context.Context, context.CancelFunc)

// DebouncedOptions configure a DebouncedFactory.
type DebouncedOptions struct {
	// Clock is used to tell time. Defaults to a monotonic clock.
	Clock clock.Clock
	// Context is the parent context used by a DebouncedFactory. Defaults
	// to context.Background().
	Context context.Context
	// Debounce is the debounce duration used by a DebouncedFactory. Defaults
	// to 1s.
	Debounce time.Duration
	// ContextFunc is the function used to generate new a context.Context.
	// Defaults to context.WithTimeout.
	ContextFunc TimeoutFunc
}

// DefaultDebouncedOptions returns default, sane DebouncedOptions.
func DefaultDebouncedOptions() DebouncedOptions {
	return _defaultDebouncedOptions
}

// Validate returns an error if DebouncedOptions contains invalid data.
func (o DebouncedOptions) Validate() (err error) {
	if o.Clock == nil {
		err = multierr.Append(err, ErrNilClock)
	}

	if o.Context == nil {
		err = multierr.Append(err, ErrNilContext)
	}

	if o.Debounce < 0 {
		err = multierr.Append(
			err,
			errors.WithMessage(ErrInvalidDebounce, o.Debounce.String()),
		)
	}

	if o.ContextFunc == nil {
		err = multierr.Append(err, ErrNilContextFunc)
	}

	return
}

// With returns a new DebouncedOptions based on o with opts merged down.
func (o DebouncedOptions) With(opts ...DebouncedOption) DebouncedOptions {
	for _, opt := range opts {
		opt.apply(&o)
	}
	return o
}

func (o DebouncedOptions) apply(other *DebouncedOptions) {
	if o.Clock != nil {
		other.Clock = o.Clock
	}

	if o.Context != nil {
		other.Context = o.Context
	}

	if o.Debounce != 0 {
		other.Debounce = o.Debounce
	}

	if o.ContextFunc != nil {
		other.ContextFunc = o.ContextFunc
	}
}

// A DebouncedOption configures a DebouncedFactory.
type DebouncedOption interface {
	apply(*DebouncedOptions)
}

// WithClock configures a DebouncedFactory to use the given clock.Clock.
func WithClock(clk clock.Clock) DebouncedOption {
	return optionFunc(func(o *DebouncedOptions) {
		o.Clock = clk
	})
}

// WithContext configures a DebouncedFactory to use ctx as its parent context.
func WithContext(ctx context.Context) DebouncedOption {
	return optionFunc(func(o *DebouncedOptions) {
		o.Context = ctx
	})
}

// WithContextFunc configures a DebouncedFactory to use ctx as its parent context.
func WithContextFunc(fn TimeoutFunc) DebouncedOption {
	return optionFunc(func(o *DebouncedOptions) {
		o.ContextFunc = fn
	})
}

// WithDebounce configures a DebouncedFactory to use dur as its debounce
// period.
func WithDebounce(dur time.Duration) DebouncedOption {
	return optionFunc(func(o *DebouncedOptions) {
		o.Debounce = dur
	})
}

type optionFunc func(*DebouncedOptions)

func (f optionFunc) apply(o *DebouncedOptions) {
	f(o)
}
