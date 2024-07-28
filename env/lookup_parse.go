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
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

func parseBool(str string) (bool, error) {
	str = strings.TrimSpace(str)
	x, err := strconv.ParseBool(str)
	if err != nil {
		switch str {
		case "y", "yes", "Y", "Yes", "YES":
			return true, nil
		case "n", "no", "N", "No", "NO":
			return false, nil
		default:
			return false, wrapParsingFailedError(err)
		}
	}
	return x, nil
}

func parseInt[T constraints.Signed](str string) (T, error) {
	var (
		bits   int
		result T
	)

	switch any(result).(type) {
	case int8:
		bits = 8
	case int16:
		bits = 16
	case int32:
		bits = 32
	case int64:
		bits = 64
	default:
		// fine, e.g. int
	}

	x, err := strconv.ParseInt(strings.TrimSpace(str), 10, bits)
	return T(x), wrapParsingFailedError(err)
}

func parseUint[T constraints.Unsigned](str string) (T, error) {
	var (
		bits   int
		result T
	)

	switch any(result).(type) {
	case uint8:
		bits = 8
	case uint16:
		bits = 16
	case uint32:
		bits = 32
	case uint64:
		bits = 64
	default:
		// fine, e.g. uint
	}

	x, err := strconv.ParseUint(strings.TrimSpace(str), 10, bits)
	return T(x), wrapParsingFailedError(err)
}

func parseFloat[T constraints.Float](str string) (T, error) {
	var (
		bits   int
		result T
	)

	switch any(result).(type) {
	case float32:
		bits = 32
	case float64:
		bits = 64
	}

	x, err := strconv.ParseFloat(strings.TrimSpace(str), bits)
	return T(x), wrapParsingFailedError(err)
}
