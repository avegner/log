package file

import (
	"compress/zlib"
	"testing"
	"time"

	"github.com/avegner/log/out"
)

var (
	name  = "test.log"
	level = zlib.NoCompression
)

func TestMultiWrite(t *testing.T) {
	o := createFileOutput(t, name, level)
	done := make(chan struct{})

	r := func(rec string) {
		for {
			select {
			case <-done:
				return
			default:
				_, _ = o.Write([]byte(rec))
			}
		}
	}

	go r("i am routine 1\n")
	go r("i am routine 2\n")

	<-time.After(1 * time.Second)
	close(done)
}

func createFileOutput(t *testing.T, name string, level int) out.Outputter {
	o, err := New(name, 0644, false, level)
	if err != nil {
		t.Fatalf("New(): got '%v' error, want no error", err)
	}
	return o
}
