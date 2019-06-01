package out

import (
	"testing"
	"time"
)

var (
	address = "localhost:50300"
	network = "tcp"
)

func TestMultiWriteNet(t *testing.T) {
	o := createNetOutput(t, network, address)
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

func createNetOutput(t *testing.T, network, address string) Outputter {
	o, err := NewNetOut(network, address)
	if err != nil {
		t.Fatalf("NewNetOut(): got '%v' error, want no error", err)
	}
	return o
}
