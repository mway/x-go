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

package env

import (
	"errors"
	"fmt"

	"golang.org/x/exp/constraints"
)

// ErrValueParsingFailed indicates that a string value from the program's
// environment could not be parsed into a requested data type.
var ErrValueParsingFailed = errors.New("failed to parse environment value")

// Parsable describes a type that can be parsed into a type.
type Parsable interface {
	constraints.Integer | constraints.Float | ~string | ~byte | ~bool
}

// Get gets the value of the environment variable with the specified name. Note
// that an empty string will be returned if either the environment itself
// contains an empty string or if the name is not found.
func Get(name string, opts ...LookupOption) string {
	x, _ := GetAs[string](name, opts...) //nolint:errcheck
	return x
}

// GetAs gets the value of the environment variable with the specified name,
// attempting to parse it into the specified type. If the parsed value cannot
// be reasonably parsed into the given type, an error will be returned. Note
// that a zero value will be returned if either the environment itself contains
// an empty string or if the name is not found.
func GetAs[T Parsable](name string, opts ...LookupOption) (T, error) {
	x, _, err := LookupAs[T](name, opts...)
	return x, err
}

// Lookup gets the value of the environment variable with the specified name,
// returning its value and whether it was present.
func Lookup(name string, opts ...LookupOption) (string, bool) {
	value, found, _ := LookupAs[string](name, opts...) //nolint:errcheck
	return value, found
}

// LookupP gets the value of the environment variable with the specified
// name, returning a pointer holding its value and whether it was present. If
// the value was present but an empty string, the returned pointer will be
// non-nil and hold an allocated empty string; if the value was not present,
// the returned pointer will be nil.
func LookupP(name string, opts ...LookupOption) (*string, bool) {
	x, found := Lookup(name, opts...)
	if !found {
		return nil, false
	}
	return &x, true
}

// LookupAs gets the value of the environment variable with the specified name,
// attempts to parse it into the specified type, and returns the value, whether
// or not the value was found, and any associated errors (for example, if the
// parsed value cannot be reasonably parsed into the given type). Note that
// the bool return is usable even when the returned error is nil; it indicates
// whether the error was encountered when interacting with the environment or
// when parsing.
//
//nolint:gocyclo
func LookupAs[T Parsable](
	name string,
	opts ...LookupOption,
) (value T, found bool, err error) {
	options := defaultLookupOptions().With(opts...)
	if options.SanitizeNames {
		name = SanitizeName(name)
	}

	var str string
	if len(options.LookupFuncs) == 0 {
		str, found = _osLookupEnv(name)
	} else {
		for _, lookup := range options.LookupFuncs {
			if str, found = lookup(name); found {
				break
			}
		}
	}
	if !found {
		return value, found, err
	}

	switch any(value).(type) {
	case bool:
		x, parseErr := parseBool(str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case int:
		x, parseErr := parseInt[int](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case int8:
		x, parseErr := parseInt[int8](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case int16:
		x, parseErr := parseInt[int16](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case int32:
		x, parseErr := parseInt[int32](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case int64:
		x, parseErr := parseInt[int64](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case uint:
		x, parseErr := parseUint[uint](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case uint8:
		x, parseErr := parseUint[uint8](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case uint16:
		x, parseErr := parseUint[uint16](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case uint32:
		x, parseErr := parseUint[uint32](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case uint64:
		x, parseErr := parseUint[uint64](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case float32:
		x, parseErr := parseFloat[float32](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	case float64:
		x, parseErr := parseFloat[float64](str)
		return any(x).(T), true, parseErr //nolint:errcheck
	default:
		// n.b. This case is only selected for string, because it is the only
		//      remaining unenumerated type that matches the Parsable type
		//      constraint. Unfortunately, the compiler doesn't see the default
		//      case as redundant (i.e., while adding a string case makes this
		//      switch exhaustive, the type switch doesn't understand the type
		//      constraint, and so is seen as inexhaustible). To work around
		//      this (and to make coverage fairer), we use default for the
		//      string case.
		return any(str).(T), true, nil //nolint:errcheck
	}
}

// LookupAsP gets the value of the environment variable with the specified
// name, attempts to parse it into the specified type, returns a pointer
// containing the parsed value, whether or not the value was found, and any
// associated errors (for example, if the parsed value cannot be reasonably
// parsed into the given type). Note that the bool return is usable even when
// the returned error is nil; it indicates whether the error was encountered
// when interacting with the environment or when parsing.
func LookupAsP[T Parsable](
	name string,
	opts ...LookupOption,
) (*T, bool, error) {
	x, found, err := LookupAs[T](name, opts...)
	if err != nil || !found {
		return nil, found, err
	}
	return &x, true, nil
}

func wrapParsingFailedError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %w", ErrValueParsingFailed, err)
}
