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
	"bytes"
	"strings"

	"go.mway.dev/pool"
)

var _bufpool = pool.NewWithReleaser(
	func() *bytes.Buffer {
		return bytes.NewBuffer(make([]byte, 0, 64))
	},
	func(x *bytes.Buffer) {
		x.Reset()
	},
)

// SanitizeName sanitizes str to be appropriate for using as a shell variable
// name. Characters outside of [A-Za-z0-9_] are replaced with underscores, and
// lowercase characters are converted to uppercase.
//
//nolint:gocyclo
func SanitizeName(str string) string {
	buf := _bufpool.Get()
	defer _bufpool.Put(buf)

	// Conditionally write a rune to the buffer; supports backfilling.
	maybeWriteRune := func(pos int, r rune, replace bool) {
		switch {
		case buf.Len() > 0:
			// If the buffer has anything in it, we'll need to write the given
			// rune below, regardless of whether it's a replacement character
			// or not.
		case !replace:
			// If we're not replacing anything and the buffer is empty, there's
			// no reason (yet) to store this rune. If the entire string is
			// processed this way, no allocations will occur.
			return
		case pos > 0:
			// If we're replacing a character for the first time and it's not
			// the first character seen, we need to backfill the portions of
			// the string that we skipped.
			buf.WriteString(str[:pos])
		}

		_, _ = buf.WriteRune(r) // nolint_errcheck
	}

	str = strings.TrimSpace(str)
	for i, r := range str {
		switch {
		case ('A' <= r && r <= 'Z') || r == '_':
			// This rune is fine.
			maybeWriteRune(i, r, false /* replace */)
		case 'a' <= r && r <= 'z':
			// This rune needs to be converted to uppercase; flip the 5th bit.
			maybeWriteRune(i, r^0x20, true /* replace */)
		case '0' <= r && r <= '9':
			// If this is the first character of the string, prefix the string
			// with an underscore (variables can't start with a digit) so that
			// it's valid. Regardless, write the digit verbatim.
			if i > 0 {
				maybeWriteRune(i, r, false /* replace */)
			} else {
				_ = buf.WriteByte('_')
				maybeWriteRune(i, r, true)
			}
		default:
			// This is a character outside of the allowed set ([A-Za-z0-9_])
			// and should be replaced with a delimiter (`_`).
			maybeWriteRune(i, '_', true /* replace */)
		}
	}

	// If any data's been buffered, use it.
	if buf.Len() > 0 {
		return buf.String()
	}

	return str
}
