package out

import (
	"compress/zlib"
	"testing"
	"time"
)

var (
	fileName      = "test.log"
	compressLevel = zlib.NoCompression
)

func TestMultiWriteFile(t *testing.T) {
	o := createFileOut(t, fileName, compressLevel)
	defer closeOut(t, o)
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

func createFileOut(t *testing.T, name string, level int) Outputter {
	o, err := NewFileOut(name, 0644, false, level)
	if err != nil {
		t.Fatalf("NewFileOut(): got '%v' error, want no error", err)
	}
	return o
}
