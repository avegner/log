package out

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

func NewNetOut(network, address string) (Outputter, error) {
	o := &netOut{
		network: network,
		address: address,
		done:    make(chan struct{}),
	}

	o.wg.Add(1)
	go o.connect()

	return o, nil
}

type netOut struct {
	conn    net.Conn
	mu      sync.Mutex
	network string
	address string
	done    chan struct{}
	wg      sync.WaitGroup
}

func (o *netOut) Write(b []byte) (n int, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.conn == nil {
		return 0, errors.New("not connected")
	}

	n, err = o.conn.Write(b)
	if err != nil {
		_ = o.conn.Close()
		o.conn = nil
		o.wg.Add(1)
		go o.connect()
	}

	return
}

func (o *netOut) Flush() {
}

func (o *netOut) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	select {
	case <-o.done:
		return nil
	default:
	}

	close(o.done)
	o.wg.Wait()

	if o.conn == nil {
		return nil
	}

	defer func() {
		o.conn = nil
	}()
	return o.conn.Close()
}

func (o *netOut) connect() {
	defer o.wg.Done()
reconnect:
	if err := o.dial(); err != nil {
		select {
		case <-time.After(1 * time.Second):
			goto reconnect
		case <-o.done:
			return
		}
	}
}

func (o *netOut) dial() error {
	errc := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	o.wg.Add(1)
	go func() {
		defer o.wg.Done()

		d := net.Dialer{}
		var err error

		conn, err := d.DialContext(ctx, o.network, o.address)

		o.mu.Lock()
		o.conn = conn
		o.mu.Unlock()
		errc <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-o.done:
		return errors.New("done")
	case err := <-errc:
		return err
	}
}
