package out

import (
	"errors"
	"io"
)

var (
	ErrClosed = errors.New("closed")
)

type Outputter interface {
	io.Writer
	io.Closer
	Flush() error
}
