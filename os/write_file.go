package os

import (
	"io"
	"os"

	"go.mway.dev/errors"
)

const (
	_fileFlags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	_fileMode  = 0o644
)

var (
	// ErrNilWriteReader indicates that a given reader value is nil or produced
	// a nil reader.
	ErrNilWriteReader = errors.New("nil write reader provided")
	// ErrUnsupportedWriteReader indicates that a given value is not a
	// supported write reader type.
	ErrUnsupportedWriteReader = errors.New("unsupported write reader")
)

// ReaderFunc is a function that returns an [io.Reader] or an error.
type ReaderFunc = func() (io.Reader, error)

// ReadCloserFunc is a function that returns an [io.ReadCloser] or an error.
type ReadCloserFunc = func() (io.ReadCloser, error)

// WriteReaderToFileWithFlags creates path with flags and mode, and copies the
// given reader to it.
func WriteReaderToFileWithFlags(
	path string,
	reader any,
	flags int,
	mode os.FileMode,
) (written int, err error) {
	var src io.Reader
	switch t := reader.(type) {
	case io.ReadCloser:
		defer func() {
			err = errors.Join(err, t.Close())
		}()
		src = t
	case io.Reader:
		src = t
	case ReadCloserFunc:
		var tmp io.ReadCloser
		if tmp, err = t(); err != nil {
			return 0, err
		}
		defer func() {
			err = errors.Join(err, tmp.Close())
		}()
		src = tmp
	case ReaderFunc:
		if src, err = t(); err != nil {
			return 0, err
		}
	default:
		return 0, errors.Wrapf(ErrUnsupportedWriteReader, "%T", src)
	}

	var dst io.WriteCloser
	if dst, err = os.OpenFile(path, flags, mode); err != nil {
		return 0, err
	}
	defer func() {
		err = errors.Join(err, dst.Close())
	}()

	var n int64
	n, err = io.Copy(dst, src)
	return int(n), err
}

// WriteReaderToFile creates path and copies the given reader to it.
func WriteReaderToFile(path string, reader any) (written int, err error) {
	return WriteReaderToFileWithFlags(path, reader, _fileFlags, _fileMode)
}
