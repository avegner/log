package net

import (
	"testing"
	"time"

	"github.com/avegner/log/out"
)

var (
	address = "localhost:50300"
	network = "tcp"
)

func TestMultiWrite(t *testing.T) {
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

func createNetOutput(t *testing.T, network, address string) out.Outputter {
	o, err := New(network, address)
	if err != nil {
		t.Fatalf("New(): got '%v' error, want no error", err)
	}
	return o
}
