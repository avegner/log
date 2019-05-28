package file

import (
	"compress/zlib"
	"os"

	"github.com/avegner/log/out"
)

func New(name string, perm os.FileMode, append bool, comprLevel int) (out.Outputter, error) {
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

	return &output{z: z}, nil
}

type output struct {
	z *zlib.Writer
}

func (o *output) Write(p []byte) (n int, err error) {
	return o.z.Write(p)
}

func (o *output) Flush() {
	o.z.Flush()
}

func (o *output) Close() error {
	return o.z.Close()
}
