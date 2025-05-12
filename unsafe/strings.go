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

package unsafe

import (
	"unsafe"
)

// StringToBytes returns the underlying byte storage of x. The returned storage
// must not be modified: in the best case, it will cause undefined behavior; in
// the worst case, the program will segfault.
//
// The lifetime of the resulting storage is inextricable from that of x: they
// are both guaranteed to live for as long as there is a live reference to
// either.
//
// Do not use this function lightly, and never leak the returned storage
// outside of the lifetime of the caller.
func StringToBytes(x string) []byte {
	if len(x) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(x), len(x))
}

// BytesToString returns a string that uses x as its underlying storage
// verbatim. The provided bytes must not be modified: in the best case, it will
// create subtle bugs; in the worst case, the program will segfault.
//
// The lifetime of the resulting string is inextricable from that of x: they
// are both guaranteed to live for as long as there is a live reference to
// either.
//
// Do not use this function lightly, and never leak the returned string outside
// of the lifetime of the caller.
func BytesToString(x []byte) string {
	if len(x) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(x), len(x))
}
