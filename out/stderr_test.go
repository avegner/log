package out

import (
	"testing"
	"time"
)

func TestMultiWriteStderr(t *testing.T) {
	o := createStderrOut(t)
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

func createStderrOut(t *testing.T) Outputter {
	o, err := NewStderrOut()
	if err != nil {
		t.Fatalf("NewStderrOut(): got '%v' error, want no error", err)
	}
	return o
}
