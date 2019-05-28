package out

import "io"

type Outputter interface {
	io.Writer
	Flusher
	io.Closer
}

type Flusher interface {
	Flush()
}
