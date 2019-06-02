package out

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

func NewNetOut(network, address string, queueSize int) (Outputter, error) {
	o := &netOut{
		network: network,
		address: address,
		queue:   make(chan []byte, queueSize),
		done:    make(chan struct{}),
	}

	o.wg.Add(1)
	go o.output()

	return o, nil
}

type netOut struct {
	conn    net.Conn
	network string
	address string
	queue   chan []byte
	done    chan struct{}
	wg      sync.WaitGroup
}

func (o *netOut) Write(bs []byte) (n int, err error) {
	select {
	case <-o.done:
		return 0, ErrClosed
	default:
	}

	select {
	case o.queue <- bs:
		return len(bs), nil
	default:
		return 0, errors.New("queue overflow")
	}
}

func (o *netOut) output() {
	defer o.wg.Done()

reconnect:
	if err := o.connect(); err != nil {
		if err == ErrClosed {
			return
		}
		select {
		case <-time.After(100 * time.Millisecond):
			goto reconnect
		case <-o.done:
			return
		}
	}

	for {
		select {
		case bs := <-o.queue:
			// TODO: use write timeout
			if _, err := o.conn.Write(bs); err != nil {
				_ = o.conn.Close()
				o.conn = nil
				goto reconnect
			}
		case <-o.done:
			return
		}
	}
}

func (o *netOut) Flush() error {
	select {
	case <-o.done:
		return ErrClosed
	default:
	}

	return nil
}

func (o *netOut) Close() error {
	select {
	case <-o.done:
		return ErrClosed
	default:
	}

	close(o.done)
	o.wg.Wait()

	if o.conn == nil {
		return nil
	}
	return o.conn.Close()
}

func (o *netOut) connect() error {
	errc := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	go func() {
		errc <- o.dialContext(ctx)
	}()

	select {
	case <-o.done:
		return ErrClosed
	case err := <-errc:
		return err
	}
}

func (o *netOut) dialContext(ctx context.Context) error {
	d := net.Dialer{}
	var err error

	o.conn, err = d.DialContext(ctx, o.network, o.address)
	return err
}
