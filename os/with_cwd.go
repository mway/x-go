package os

import (
	"os"
)

// n.b. See with_cwd_internals_test.go.
var (
	_getwd = os.Getwd
	_chdir = os.Chdir
)

// WithCwd attempts to change the working directory to dir and, if successful,
// calls f. WithCwd expects dir to exist already; if dir does not exist, or is
// removed or renamed during f's execution, an error will be returned.
//
// WithCwd will attempt to restore the original working directory, even if the
// given function panics.
func WithCwd(dir string, f func()) (err error) {
	var orig string
	if orig, err = _getwd(); err != nil {
		return
	}

	// If we're already in the target directory, just execute f.
	if dir == orig {
		f()
		return
	}

	if err = _chdir(dir); err != nil {
		return
	}

	defer func() {
		err = os.Chdir(orig)
	}()

	f()
	return
}