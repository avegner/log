package net

import (
	"net"

	"github.com/avegner/log/out"
)

func New(network, address string) (out.Outputter, error) {
	c, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &output{c: c}, nil
}

type output struct {
	c net.Conn
}

func (o *output) Write(b []byte) (n int, err error) {
	return o.c.Write(b)
}

func (o *output) Flush() {
}

func (o *output) Close() error {
	return o.c.Close()
}
