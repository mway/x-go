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

package exec_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	goexec "os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	xos "go.mway.dev/x/os"
	"go.mway.dev/x/os/exec"
	"go.mway.dev/x/os/tempdir"
	"go.mway.dev/x/stub"
)

var (
	_nopOption  = func(*goexec.Cmd) {}
	_echoScript = "testdata/echo.sh"
)

func TestRunCommand(t *testing.T) {
	cases := map[string]struct { //nolint:govet
		giveName    string
		giveOptions exec.CommandOption
		checkError  func(*testing.T, error)
	}{
		"nominal without args": {
			giveName:    "echo",
			giveOptions: _nopOption,
			checkError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"nominal with args": {
			giveName: "echo",
			giveOptions: exec.CommandOptions(
				exec.WithArgs(t.Name()),
			),
			checkError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"no command": {
			giveName:    "",
			giveOptions: _nopOption,
			checkError: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "exec: no command")
			},
		},
		"command not exist": {
			giveName:    t.Name(),
			giveOptions: _nopOption,
			checkError: func(t *testing.T, err error) {
				require.ErrorIs(t, err, goexec.ErrNotFound)
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			tt.checkError(t, exec.RunCommand(
				tt.giveName,
				tt.giveOptions,
			))
		})
	}
}

func TestRunCommandOutput(t *testing.T) {
	cases := map[string]bool{
		"stdout": true,
		"stderr": false,
	}

	for dst, wantOutput := range cases {
		t.Run(dst, func(t *testing.T) {
			output, err := exec.RunCommandOutput(
				_echoScript,
				exec.WithArgs(dst, t.Name()),
			)
			require.NoError(t, err)
			if wantOutput {
				require.Equal(t, t.Name(), strings.TrimSpace(output))
			} else {
				require.Empty(t, output)
			}
		})
	}
}

func TestRunCommandSplitOutput(t *testing.T) {
	cases := map[string]struct {
		wantStdout bool
		wantStderr bool
	}{
		"stdout": {
			wantStdout: true,
			wantStderr: false,
		},
		"stderr": {
			wantStdout: false,
			wantStderr: true,
		},
	}

	for dst, tt := range cases {
		t.Run(dst, func(t *testing.T) {
			stdout, stderr, err := exec.RunCommandSplitOutput(
				_echoScript,
				exec.WithArgs(dst, t.Name()),
			)
			require.NoError(t, err)

			if tt.wantStdout {
				require.Equal(t, t.Name(), strings.TrimSpace(stdout))
			} else {
				require.Empty(t, stdout)
			}

			if tt.wantStderr {
				require.Equal(t, t.Name(), strings.TrimSpace(stderr))
			} else {
				require.Empty(t, stderr)
			}
		})
	}
}

func TestRunCommandCombinedOutput(t *testing.T) {
	output, err := exec.RunCommandCombinedOutput("sh", exec.WithArgs(
		"-c",
		strings.TrimSpace(`
			testdata/echo.sh stdout 1:stdout
			testdata/echo.sh stdout 2:stderr
			testdata/echo.sh stdout 3:stdout
			testdata/echo.sh stdout 4:stderr
		`),
	))
	require.NoError(t, err)
	output = strings.TrimSpace(output)

	lines := strings.Split(output, "\n")
	require.Len(t, lines, 4)

	require.Equal(t, "1:stdout", lines[0])
	require.Equal(t, "2:stderr", lines[1])
	require.Equal(t, "3:stdout", lines[2])
	require.Equal(t, "4:stderr", lines[3])
}

func TestRunAttachedCommand(t *testing.T) {
	testdir, err := os.Getwd()
	require.NoError(t, err)

	testfile := filepath.Join(testdir, "testdata", "echo.sh")
	err = tempdir.With(func(tmp string) {
		var (
			stdoutPath = filepath.Join(tmp, "stdout")
			stderrPath = filepath.Join(tmp, "stderr")
		)
		e := xos.WithFileWriter(stdoutPath, func(outw io.Writer) error {
			stdout, ok := outw.(*os.File)
			require.True(t, ok)

			return xos.WithFileWriter(stderrPath, func(errw io.Writer) error {
				stderr, ok := errw.(*os.File)
				require.True(t, ok)

				return stub.WithError(&os.Stdout, stdout, func() error {
					return stub.WithError(&os.Stderr, stderr, func() error {
						return exec.RunAttachedCommand("sh", exec.WithArgs(
							"-c",
							strings.TrimSpace(fmt.Sprintf(`
								%[1]s stdout 1:stdout
								%[1]s stderr 2:stderr
								%[1]s stdout 3:stdout
								%[1]s stderr 4:stderr
							`, testfile)),
						))
					})
				})
			})
		})
		require.NoError(t, e)

		requireFileEquals(t, stdoutPath, "1:stdout\n3:stdout", true)
		requireFileEquals(t, stderrPath, "2:stderr\n4:stderr", true)
	})
	require.NoError(t, err)
}

func TestRunSilentCommand(t *testing.T) {
	var (
		stdout = bytes.NewBuffer(nil)
		stderr = bytes.NewBuffer(nil)
		stdin  = bytes.NewBufferString(t.Name())
	)

	err := exec.RunSilentCommand(
		"sh",
		exec.WithArgs(
			"-c",
			strings.TrimSpace(fmt.Sprintf(`
				cat -
				%[1]s stdout 1:stdout
				%[1]s stderr 2:stderr
				%[1]s stdout 3:stdout
				%[1]s stderr 4:stderr
			`, _echoScript)),
		),
		exec.WithStdout(stdout),
		exec.WithStderr(stderr),
		exec.WithStdin(stdin),
	)
	require.NoError(t, err)

	require.Equal(t, 0, stdout.Len())
	require.Equal(t, 0, stderr.Len())
	require.Equal(t, 0, stdin.Len())
}

func TestRunCommandOutput_CommandError(t *testing.T) {
	err := tempdir.With(func(_ string) {
		output, err := exec.RunCommandOutput("does-not-exist")
		require.ErrorIs(t, err, goexec.ErrNotFound)
		require.Empty(t, output)

		output, err = exec.RunCommandOutput("does/not/exist")
		require.ErrorIs(t, err, os.ErrNotExist)
		require.Empty(t, output)
	})
	require.NoError(t, err)
}

func TestRunCommandSplitOutput_CommandError(t *testing.T) {
	err := tempdir.With(func(_ string) {
		stdout, stderr, err := exec.RunCommandSplitOutput("does-not-exist")
		require.ErrorIs(t, err, goexec.ErrNotFound)
		require.Empty(t, stdout)
		require.Empty(t, stderr)

		stdout, stderr, err = exec.RunCommandSplitOutput("does/not/exist")
		require.ErrorIs(t, err, os.ErrNotExist)
		require.Empty(t, stdout)
		require.Empty(t, stderr)
	})
	require.NoError(t, err)
}

func TestRunCommandCombinedOutput_CommandError(t *testing.T) {
	err := tempdir.With(func(_ string) {
		output, err := exec.RunCommandCombinedOutput("does-not-exist")
		require.ErrorIs(t, err, goexec.ErrNotFound)
		require.Empty(t, output)

		output, err = exec.RunCommandCombinedOutput("does/not/exist")
		require.ErrorIs(t, err, os.ErrNotExist)
		require.Empty(t, output)
	})
	require.NoError(t, err)
}

func TestRunCommand_WithDir(t *testing.T) {
	err := tempdir.With(func(wantPwd string) {
		havePwd, err := exec.RunCommandOutput("pwd", exec.WithDir(wantPwd))
		require.NoError(t, err)

		wantAbs, err := filepath.Abs(wantPwd)
		require.NoError(t, err)

		haveAbs, err := filepath.Abs(strings.TrimSpace(havePwd))
		require.NoError(t, err)

		require.Equal(t, wantAbs, haveAbs)
	})
	require.NoError(t, err)
}

func TestRunCommand_WithEnv(t *testing.T) {
	wantEnv := []string{
		fmt.Sprintf("TEST_VAR_1_%s=%s", t.Name(), t.Name()),
		fmt.Sprintf("TEST_VAR_2_%s=%s", t.Name(), t.Name()),
	}

	raw, err := exec.RunCommandOutput("env", exec.WithEnv(wantEnv))
	require.NoError(t, err)

	haveEnv := strings.Split(strings.TrimSpace(raw), "\n")
	require.Equal(t, wantEnv, haveEnv)
}

func TestRunCommand_WithStdin(t *testing.T) {
	stdin := bytes.NewBufferString(t.Name())

	output, err := exec.RunCommandOutput(
		"cat",
		exec.WithArgs("-"),
		exec.WithStdin(stdin),
	)
	require.NoError(t, err)
	require.Equal(t, 0, stdin.Len())
	require.Equal(t, t.Name(), strings.TrimSpace(output))
}

func TestRunCommand_WithNopPipes(t *testing.T) {
	var (
		stdout = bytes.NewBuffer(nil)
		stderr = bytes.NewBuffer(nil)
		stdin  = bytes.NewBufferString(t.Name())
	)

	err := exec.RunSilentCommand(
		"sh",
		exec.WithArgs(
			"-c",
			strings.TrimSpace(fmt.Sprintf(`
				cat -
				%[1]s stdout 1:stdout
				%[1]s stderr 2:stderr
				%[1]s stdout 3:stdout
				%[1]s stderr 4:stderr
			`, _echoScript)),
		),
		exec.WithStdout(stdout),
		exec.WithStderr(stderr),
		exec.WithNopPipes(), // specify nop pipes before stdin
		exec.WithStdin(stdin),
	)
	require.NoError(t, err)

	require.Equal(t, 0, stdout.Len())
	require.Equal(t, 0, stderr.Len())
	require.Equal(t, 0, stdin.Len())
}

func requireFileEquals(
	t *testing.T,
	path string,
	want string,
	wantExist bool,
) {
	raw, err := os.ReadFile(path)
	if !wantExist {
		require.ErrorIs(t, err, os.ErrNotExist)
		return
	}
	require.NoError(t, err)
	require.Equal(t, want, string(bytes.TrimSpace(raw)))
}
