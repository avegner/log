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
	Flusher
	io.Closer
}

type Flusher interface {
	Flush() error
}
