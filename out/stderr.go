package out

import (
	"fmt"
	"os"
)

func NewStderrOut() (Outputter, error) {
	return &stderrOut{
		done: make(chan struct{}),
	}, nil
}

type stderrOut struct {
	done chan struct{}
}

func (o *stderrOut) Write(p []byte) (n int, err error) {
	select {
	case <-o.done:
		return 0, ErrClosed
	default:
	}

	return fmt.Fprint(os.Stderr, string(p))
}

func (o *stderrOut) Flush() error {
	select {
	case <-o.done:
		return ErrClosed
	default:
	}

	return nil
}

func (o *stderrOut) Close() error {
	select {
	case <-o.done:
		return ErrClosed
	default:
	}

	close(o.done)
	return nil
}
