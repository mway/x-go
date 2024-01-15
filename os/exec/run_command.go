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

// Package exec provides command execution related helpers and utilities.
package exec

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"slices"

	xio "go.mway.dev/x/io"
)

// ErrNotFound re-exports [exec.ErrNotFound] for convenience.
var ErrNotFound = exec.ErrNotFound

type (
	// Error re-exports [exec.Error] for convenience.
	Error = exec.Error
	// ExitError re-exports [exec.ExitError] for convenience.
	ExitError = exec.ExitError
)

// RunCommand runs the given program or executable with the given options.
func RunCommand(name string, opts ...CommandOption) error {
	cmd := exec.Command(name)
	for _, opt := range opts {
		opt(cmd)
	}

	// If any pipes are nil, use nop to avoid opening os.DevNull.
	useFirstOf(&cmd.Stdout, xio.NopWriter)
	useFirstOf(&cmd.Stderr, xio.NopWriter)
	useFirstOf(&cmd.Stdin, xio.NopReader)

	return cmd.Run()
}

// RunAttachedCommand is a convenience alias for [RunCommand] where
// [WithAttachedPipes] is appended to the provided options.
func RunAttachedCommand(name string, opts ...CommandOption) error {
	opts = append(slices.Clone(opts), WithAttachedPipes())
	return RunCommand(name, opts...)
}

// RunSilentCommand is a convenience alias for [RunCommand] where
// [WithNopPipes] is appended to the provided options.
func RunSilentCommand(name string, opts ...CommandOption) error {
	opts = append(slices.Clone(opts), WithStdout(xio.Nop), WithStderr(xio.Nop))
	return RunCommand(name, opts...)
}

// RunCommandOutput calls [RunCommand] with the given name and options, and
// returns the command's standard output.
func RunCommandOutput(name string, opts ...CommandOption) (string, error) {
	buf := bytes.NewBuffer(nil)

	opts = append(slices.Clone(opts), WithStdout(buf))
	if err := RunCommand(name, opts...); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// RunCommandSplitOutput calls [RunCommand] with the given name and options,
// and returns the command's standard output and standard error individually.
func RunCommandSplitOutput(
	name string,
	opts ...CommandOption,
) (string, string, error) {
	var (
		outbuf = bytes.NewBuffer(nil)
		errbuf = bytes.NewBuffer(nil)
	)

	opts = append(slices.Clone(opts), WithStdout(outbuf), WithStderr(errbuf))
	if err := RunCommand(name, opts...); err != nil {
		return "", "", err
	}

	return outbuf.String(), errbuf.String(), nil
}

// RunCommandCombinedOutput calls [RunCommand] with the given name and options,
// and returns the command's standard output and standard error combined.
func RunCommandCombinedOutput(
	name string,
	opts ...CommandOption,
) (string, error) {
	buf := bytes.NewBuffer(nil)

	opts = append(opts, WithStdout(buf), WithStderr(buf))
	if err := RunCommand(name, opts...); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// A CommandOption specializes the behavior of [RunCommand] and its variants.
type CommandOption = func(*exec.Cmd)

// CommandOptions combines multiple [CommandOption] to specialize [RunCommand]
// and its variants.
func CommandOptions(opts ...CommandOption) CommandOption {
	return func(dst *exec.Cmd) {
		for _, opt := range opts {
			opt(dst)
		}
	}
}

// WithArgs returns a new [CommandOption] that uses the given args as command
// arguments.
func WithArgs(args ...string) CommandOption {
	return func(dst *exec.Cmd) {
		dst.Args = append([]string{dst.Path}, args...)
	}
}

// WithDir returns a new [CommandOption] that uses the given dir as as the
// working directory.
func WithDir(dir string) CommandOption {
	return func(dst *exec.Cmd) {
		dst.Dir = dir
	}
}

// WithEnv returns a new [CommandOption] that uses the given key=val pairs as
// the environment.
func WithEnv(env []string) CommandOption {
	return func(dst *exec.Cmd) {
		dst.Env = slices.Clone(env)
	}
}

// WithStdout returns a new [CommandOption] that uses w for standard output.
func WithStdout(w io.Writer) CommandOption {
	return func(dst *exec.Cmd) {
		dst.Stdout = w
	}
}

// WithStderr returns a new [CommandOption] that uses w for standard error.
func WithStderr(w io.Writer) CommandOption {
	return func(dst *exec.Cmd) {
		dst.Stderr = w
	}
}

// WithStdin returns a new [CommandOption] that uses w for standard input.
func WithStdin(r io.Reader) CommandOption {
	return func(dst *exec.Cmd) {
		dst.Stdin = r
	}
}

// WithNopPipes returns a new [CommandOption] that uses nop types for standard
// output, standard error, and standard input.
func WithNopPipes() CommandOption {
	return func(dst *exec.Cmd) {
		dst.Stdout = xio.Nop
		dst.Stdin = xio.Nop
		dst.Stderr = xio.Nop
	}
}

// WithAttachedPipes returns a new [CommandOption] that uses this program's
// standard output, standard error, and standard input.
func WithAttachedPipes() CommandOption {
	return func(dst *exec.Cmd) {
		dst.Stdout = os.Stdout
		dst.Stdin = os.Stdin
		dst.Stderr = os.Stderr
	}
}

func useFirstOf[T comparable](val *T, vals ...T) {
	var zero T
	if *val != zero {
		return
	}

	for _, v := range vals {
		if v != zero {
			*val = v
		}
	}
}
