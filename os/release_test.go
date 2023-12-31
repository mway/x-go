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
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/stub"
)

var (
	_runtimeGOOSDarwin      = func() string { return "darwin" }
	_runtimeGOOSLinux       = func() string { return "linux" }
	_runtimeGOOSUnsupported = func() string { return "unsupported" }
	_catCommand             = func(name string) func() *exec.Cmd {
		return func() *exec.Cmd {
			return exec.Command("cat", filepath.Join("testdata", name))
		}
	}
	_badCommand = func() *exec.Cmd {
		return _catCommand("this-command-will-fail")()
	}
)

func TestWrappers(t *testing.T) {
	t.Run("runtimeGOOS", func(t *testing.T) {
		require.Equal(t, runtime.GOOS, _runtimeGOOS())
	})
	t.Run("swVersCommand", func(t *testing.T) {
		have := _swVersCommand()
		require.Equal(t, []string{"sw_vers"}, have.Args)
	})
	t.Run("unameCommand", func(t *testing.T) {
		have := _unameCommand()
		require.Equal(t, []string{"uname", "-o"}, have.Args)
	})
	t.Run("lsbReleaseCommand", func(t *testing.T) {
		have := _lsbReleaseCommand()
		require.Equal(t, []string{"lsb_release", "-a"}, have.Args)
	})
	t.Run("readOSRelease", func(t *testing.T) {
		_, err := _readOSRelease()
		if err != nil {
			require.ErrorIs(t, err, os.ErrNotExist)
		}
	})
}

func TestRelease(t *testing.T) {
	first := Release()
	for i := 0; i < 1000; i++ {
		require.Equal(t, first, Release())
	}
}

func TestRelease_MacOS(t *testing.T) {
	cases := map[string]struct {
		cmd  func() *exec.Cmd
		want ReleaseInfo
	}{
		"sonoma": {
			cmd: _catCommand("sw_vers-sonoma"),
			want: ReleaseInfo{
				OS:       "darwin",
				Family:   "macos",
				Version:  "14.0.0",
				Codename: "sonoma",
			},
		},
		"ventura": {
			cmd: _catCommand("sw_vers-ventura"),
			want: ReleaseInfo{
				OS:       "darwin",
				Family:   "macos",
				Version:  "13.0.0",
				Codename: "ventura",
			},
		},
		"monterey": {
			cmd: _catCommand("sw_vers-monterey"),
			want: ReleaseInfo{
				OS:       "darwin",
				Family:   "macos",
				Version:  "12.0.0",
				Codename: "monterey",
			},
		},
		"legacy": {
			cmd: _catCommand("sw_vers-legacy"),
			want: ReleaseInfo{
				OS:       "darwin",
				Family:   "macos",
				Version:  "11.0.0",
				Codename: "legacy",
			},
		},
		"error": {
			cmd: _badCommand,
			want: ReleaseInfo{
				OS: "darwin",
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			stub.With(&_runtimeGOOS, _runtimeGOOSDarwin, func() {
				stub.With(&_swVersCommand, tt.cmd, func() {
					have := loadRelease()
					require.Equal(t, tt.want, have)
				})
			})
		})
	}
}

func TestRelease_Debian_OSRelease(t *testing.T) {
	cases := map[string]struct {
		cmd  func() ([]byte, error)
		want ReleaseInfo
	}{
		"nominal": {
			cmd: func() ([]byte, error) {
				return os.ReadFile("testdata/os-release")
			},
			want: ReleaseInfo{
				OS:       "linux",
				Family:   "debian",
				Version:  "12",
				Codename: "bookworm",
			},
		},
		"error": {
			cmd: func() ([]byte, error) {
				return nil, errors.New("error")
			},
			want: ReleaseInfo{
				OS: "linux",
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			stub.With(&_runtimeGOOS, _runtimeGOOSLinux, func() {
				stub.With(&_readOSRelease, tt.cmd, func() {
					stub.With(&_lsbReleaseCommand, _badCommand, func() {
						have := loadRelease()
						require.Equal(t, tt.want, have)
					})
				})
			})
		})
	}
}

func TestRelease_Debian_LSBRelease(t *testing.T) {
	var (
		badOSRelease = func() ([]byte, error) {
			return nil, errors.New("error")
		}
		lsbRelease = func() *exec.Cmd {
			return exec.Command("cat", "testdata/lsb_release")
		}

		cases = map[string]struct {
			cmd  func() *exec.Cmd
			want ReleaseInfo
		}{
			"nominal": {
				cmd: lsbRelease,
				want: ReleaseInfo{
					OS:       "linux",
					Family:   "debian",
					Version:  "12",
					Codename: "bookworm",
				},
			},
			"error": {
				cmd: _badCommand,
				want: ReleaseInfo{
					OS: "linux",
				},
			},
		}
	)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			stub.With(&_runtimeGOOS, _runtimeGOOSLinux, func() {
				stub.With(&_readOSRelease, badOSRelease, func() {
					stub.With(&_lsbReleaseCommand, tt.cmd, func() {
						have := loadRelease()
						require.Equal(t, tt.want, have)
					})
				})
			})
		})
	}
}

func TestRelease_Unknown(t *testing.T) {
	stub.With(&_runtimeGOOS, _runtimeGOOSUnsupported, func() {
		want := ReleaseInfo{
			OS:       _runtimeGOOSUnsupported(),
			Family:   "unknown",
			Version:  "unknown",
			Codename: "unknown",
		}
		require.Equal(t, want, loadRelease())
	})
}
