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

package os

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	_runtimeGOARCH = func() string {
		return runtime.GOARCH
	}
	_runtimeGOOS = func() string {
		return runtime.GOOS
	}
	_releaseInfo   = sync.OnceValue(loadRelease)
	_swVersCommand = func() *exec.Cmd {
		return exec.Command("sw_vers")
	}
	_unameCommand = func() *exec.Cmd {
		return exec.Command("uname", "-o")
	}
	_lsbReleaseCommand = func() *exec.Cmd {
		return exec.Command("lsb_release", "-a")
	}
	_readOSRelease = func() ([]byte, error) {
		return os.ReadFile("/etc/os-release")
	}
	_newlineBytes = []byte{'\n'}
)

// ReleaseInfo provides information about the current OS release.
type ReleaseInfo struct {
	Arch     string
	OS       string
	Family   string
	Version  string
	Codename string
}

// Release returns the current system's [ReleaseInfo].
func Release() ReleaseInfo {
	return _releaseInfo()
}

func loadRelease() (info ReleaseInfo) {
	info.Arch = _runtimeGOARCH()
	info.OS = _runtimeGOOS()
	switch info.OS {
	case "darwin":
		loadDarwinRelease(&info)
	case "linux":
		loadLinuxRelease(&info)
	default:
		const unknown = "unknown"
		info.Family = unknown
		info.Codename = unknown
		info.Version = unknown
	}

	return
}

func parseOSRelease(dst *ReleaseInfo, raw []byte) {
	lines := bytes.Split(raw, _newlineBytes)
	for _, line := range lines {
		parts := strings.Split(string(line), "=")
		if len(parts) < 2 {
			continue
		}
		for i := range parts {
			parts[i] = strings.ToLower(strings.TrimSpace(parts[i]))
		}

		switch parts[0] {
		case "id":
			dst.Family = parts[1]
		case "version_id":
			dst.Version = strings.Trim(parts[1], `"`)
		case "version_codename":
			dst.Codename = parts[1]
		default:
			// ignore other fields
		}
	}
}

func parseLSBRelease(dst *ReleaseInfo, raw []byte) {
	lines := bytes.Split(raw, _newlineBytes)
	for _, line := range lines {
		parts := strings.Split(string(line), ":")
		if len(parts) < 2 {
			continue
		}
		for i := range parts {
			parts[i] = strings.ToLower(strings.TrimSpace(parts[i]))
		}

		switch parts[0] {
		case "distributor id":
			dst.Family = parts[1]
		case "release":
			dst.Version = parts[1]
		case "codename":
			dst.Codename = parts[1]
		default:
			// ignore other fields
		}
	}
}

func loadLinuxRelease(dst *ReleaseInfo) {
	if raw, err := _readOSRelease(); err != nil {
		if raw, err = _lsbReleaseCommand().Output(); err != nil {
			// nothing to do
			return
		}
		parseLSBRelease(dst, raw)
	} else {
		parseOSRelease(dst, raw)
	}
}

func loadDarwinRelease(dst *ReleaseInfo) {
	raw, err := _swVersCommand().Output()
	if err != nil {
		return
	}

	lines := bytes.Split(raw, _newlineBytes)
	for _, line := range lines {
		parts := strings.Split(string(line), ":")
		if len(parts) < 2 {
			continue
		}
		for i := range parts {
			parts[i] = strings.ToLower(strings.TrimSpace(parts[i]))
		}

		switch parts[0] {
		case "productname":
			dst.Family = parts[1]
		case "productversion":
			dst.Version = parts[1]
			setDarwinCodename(dst)
		default:
			// ignore other fields
		}
	}
}

func setDarwinCodename(dst *ReleaseInfo) {
	switch dst.Family {
	case "macos":
		switch {
		case strings.HasPrefix(dst.Version, "14."):
			dst.Codename = "sonoma"
		case strings.HasPrefix(dst.Version, "13."):
			dst.Codename = "ventura"
		case strings.HasPrefix(dst.Version, "12."):
			dst.Codename = "monterey"
		case strings.HasPrefix(dst.Version, "11."):
			dst.Codename = "bigsur"
		default:
			dst.Codename = "legacy"
		}
	default:
		// only macos is supported for darwin
	}
}
