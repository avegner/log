package out

import (
	"compress/zlib"
	"os"
	"sync"
)

func NewFileOut(name string, perm os.FileMode, append bool, comprLevel int) (Outputter, error) {
	flags := os.O_RDWR | os.O_CREATE
	if append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	f, err := os.OpenFile(name, flags, perm)
	if err != nil {
		return nil, err
	}

	z, err := zlib.NewWriterLevel(f, comprLevel)
	if err != nil {
		return nil, err
	}

	return &fileOut{
		z:    z,
		done: make(chan struct{}),
	}, nil
}

type fileOut struct {
	mu   sync.Mutex
	z    *zlib.Writer
	done chan struct{}
}

func (o *fileOut) Write(bs []byte) (n int, err error) {
	select {
	case <-o.done:
		return 0, ErrClosed
	default:
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	return o.z.Write(bs)
}

func (o *fileOut) Flush() error {
	select {
	case <-o.done:
		return ErrClosed
	default:
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	return o.z.Flush()
}

func (o *fileOut) Close() error {
	select {
	case <-o.done:
		return ErrClosed
	default:
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	return o.z.Close()
}
