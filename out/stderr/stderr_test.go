package stderr

import (
	"testing"
	"time"

	"github.com/avegner/log/out"
)

func TestMultiWrite(t *testing.T) {
	o := createStderrOutput(t)
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

func createStderrOutput(t *testing.T) out.Outputter {
	o, err := New()
	if err != nil {
		t.Fatalf("New(): got '%v' error, want no error", err)
	}
	return o
}
