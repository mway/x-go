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

// Package filepath provides filepath-related utilities.
package filepath

import (
	"os"
	"path/filepath"
	"strings"
)

const _sep = string(os.PathSeparator)

// Ancestors returns a list of all ancestor paths contained within path.
func Ancestors(path string) []string {
	path = filepath.Clean(path)
	switch path {
	case "", ".", "..", _sep:
		return nil
	}

	var offset int
	if path[0] == _sep[0] {
		offset++
	}

	size := strings.Count(path, "/") - offset
	if size == 0 {
		return nil
	}

	paths := make([]string, 0, strings.Count(path, "/"))
	for {
		if path = filepath.Dir(path); path == "." || path == _sep {
			break
		}
		paths = append(paths, path)
	}
	return paths
}
