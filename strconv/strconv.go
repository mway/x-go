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

// Package strconv provides helpers for parsing strings into other data types.
package strconv

import (
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Parsable describes any type that this package's parsing supports.
type Parsable interface {
	constraints.Integer | constraints.Float | ~bool | ~string
}

// Parse attempts to parse str into a T, returning any error(s) encountered.
//
//nolint:gocyclo
func Parse[T Parsable](str string) (T, error) {
	var (
		value T
		err   error
	)

	switch v := any(&value).(type) {
	case *bool:
		*v, err = ParseBool(str)
	case *int:
		*v, err = ParseInt[int](str)
	case *int8:
		*v, err = ParseInt[int8](str)
	case *int16:
		*v, err = ParseInt[int16](str)
	case *int32:
		*v, err = ParseInt[int32](str)
	case *int64:
		*v, err = ParseInt[int64](str)
	case *uint:
		*v, err = ParseUint[uint](str)
	case *uint8:
		*v, err = ParseUint[uint8](str)
	case *uint16:
		*v, err = ParseUint[uint16](str)
	case *uint32:
		*v, err = ParseUint[uint32](str)
	case *uint64:
		*v, err = ParseUint[uint64](str)
	case *float32:
		*v, err = ParseFloat[float32](str)
	case *float64:
		*v, err = ParseFloat[float64](str)
	case *string:
		*v = str
	default:
	}

	return value, err
}

// ParseBool attempts to parse str into a bool, returning any error(s)
// encountered. The parsable "truthy" values are:
//
//	1 t T true True TRUE  y Y yes Yes YES
//	--------------------  ---------------
//	 stdlib                    extension
//
// Similarly, the parsable "falsy" values are:
//
//	0 f F false False FALSE  n N no No NO
//	-----------------------  ------------
//	 stdlib                    extension
//
// Any other value, including incorrectly-cased variants of allowed values
// (e.g. "tRuE") and empty strings, will return an error.
func ParseBool(str string) (bool, error) {
	str = strings.TrimSpace(str)

	x, err := strconv.ParseBool(str)
	if err == nil {
		return x, nil
	}

	switch str {
	case "y", "yes", "Y", "Yes", "YES":
		return true, nil
	case "n", "no", "N", "No", "NO":
		return false, nil
	default:
		return false, err
	}
}

// ParseInt attempts to parse str into a T, returning any error(s) encountered.
func ParseInt[T constraints.Signed](str string) (T, error) {
	x, err := strconv.ParseInt(strings.TrimSpace(str), 10, bitsize[T]())
	return T(x), err
}

// ParseUint attempts to parse str into a T, returning any error(s)
// encountered.
func ParseUint[T constraints.Unsigned](str string) (T, error) {
	x, err := strconv.ParseUint(strings.TrimSpace(str), 10, bitsize[T]())
	return T(x), err
}

// ParseFloat attempts to parse str into a T, returning any error(s)
// encountered.
func ParseFloat[T constraints.Float](str string) (T, error) {
	x, err := strconv.ParseFloat(strings.TrimSpace(str), bitsize[T]())
	return T(x), err
}

func bitsize[T constraints.Integer | constraints.Float]() int {
	var x T
	return int(unsafe.Sizeof(x) * 8)
}
