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

	return &fileOut{z: z}, nil
}

type fileOut struct {
	mu   sync.Mutex
	z    *zlib.Writer
}

func (o *fileOut) Write(p []byte) (n int, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.z.Write(p)
}

func (o *fileOut) Flush() {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.z.Flush()
}

func (o *fileOut) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.z.Close()
}
