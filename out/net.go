package out

import (
	"net"
)

func NewNetOut(network, address string) (Outputter, error) {
	c, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &netOut{c: c}, nil
}

type netOut struct {
	c net.Conn
}

func (o *netOut) Write(b []byte) (n int, err error) {
	return o.c.Write(b)
}

func (o *netOut) Flush() {
}

func (o *netOut) Close() error {
	return o.c.Close()
}
