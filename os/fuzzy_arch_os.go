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

package os

import (
	"runtime"
	"sort"
	"strings"

	"go.mway.dev/x/container/set"
)

const (
	// LocalFuzzyArch is a [FuzzyArch] that corresponds to the local system.
	LocalFuzzyArch = FuzzyArch(runtime.GOARCH)
	// LocalFuzzyOS is a [FuzzyOS] that corresponds to the local system.
	LocalFuzzyOS = FuzzyOS(runtime.GOOS)
)

var (
	_amdNames    = [...]string{"amd64", "x64", "x86_64"}
	_armNames    = [...]string{"aarch64", "arm64"}
	_darwinNames = [...]string{"darwin", "macos"}
	_debianNames = [...]string{"debian", "linux", "ubuntu"}
	_archIndex   = make(map[string][]string)
	_osIndex     = make(map[string][]string)
)

func init() {
	for _, name := range _amdNames {
		_archIndex[name] = _amdNames[:]
	}
	for _, name := range _armNames {
		_archIndex[name] = _armNames[:]
	}
	for _, name := range _darwinNames {
		_osIndex[name] = _darwinNames[:]
	}
	for _, name := range _debianNames {
		_osIndex[name] = _debianNames[:]
	}
}

// FuzzyArch is a fuzzy runtime.GOARCH-like value that allows alternate name
// matching. See [FuzzyArch.Names] and [FuzzyArch.Matches] for more details.
type FuzzyArch string

// Names returns any names by which this [FuzzyArch] might be known. Always
// includes the value of this FuzzyArch verbatim.
func (a FuzzyArch) Names() []string {
	return appendName(_archIndex, a.String())
}

// Matches returns whether this [FuzzyArch] is a substring of str.
func (a FuzzyArch) Matches(str string) bool {
	return matches(a.Names(), str)
}

// String returns this [FuzzyArch] as a string.
func (a FuzzyArch) String() string {
	return string(a)
}

// FuzzyOS is a fuzzy runtime.GOOS-like value that allows alternate name
// matching. See [FuzzyOS.Names] and [FuzzyOS.Matches] for more details.
type FuzzyOS string

// Names returns any names by which this [FuzzyOS] might be known. Always
// includes the value of this FuzzyOS verbatim.
func (o FuzzyOS) Names() []string {
	return appendName(_osIndex, o.String())
}

// Matches returns whether this [FuzzyOS] is a substring of str.
func (o FuzzyOS) Matches(str string) bool {
	return matches(o.Names(), str)
}

// String returns this [FuzzyOS] as a string.
func (o FuzzyOS) String() string {
	return string(o)
}

func appendName(index map[string][]string, name string) []string {
	if name = strings.ToLower(strings.TrimSpace(name)); len(name) == 0 {
		return nil
	}

	names, ok := index[name]
	if len(names) == 0 || !ok {
		return []string{name}
	}

	tmp := set.New(names...)
	tmp.Add(name)

	names = tmp.ToSlice()
	sort.Strings(names)
	return names
}

func matches(names []string, target string) bool {
	target = strings.ToLower(target)
	for _, name := range names {
		if name == "*" || strings.Contains(target, name) {
			return true
		}
	}
	return false
}
