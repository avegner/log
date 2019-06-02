package out

import (
	"fmt"
	"os"
)

func NewStderrOut() (Outputter, error) {
	return &stderrOut{}, nil
}

type stderrOut struct {
}

func (o *stderrOut) Write(p []byte) (n int, err error) {
	return fmt.Fprint(os.Stderr, string(p))
}

func (o *stderrOut) Flush() error {
	return nil
}

func (o *stderrOut) Close() error {
	return nil
}
