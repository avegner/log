package stderr

import (
	"fmt"
	"os"

	"github.com/avegner/log/out"
)

func New() (out.Outputter, error) {
	return &output{}, nil
}

type output struct {
}

func (o *output) Write(p []byte) (n int, err error) {
	return fmt.Fprint(os.Stderr, string(p))
}

func (o *output) Flush() {
}

func (o *output) Close() error {
	return nil
}
