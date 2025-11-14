package net

import (
	"bytes"
	"io"
	"net"
	"sync"
	"time"

	"github.com/nanachi-sh/go-mc/constants"
)

type Conn struct {
	net.Conn

	lock sync.Mutex
}

func Dial(network string, address string) (*Conn, error) {
	c, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &Conn{
		Conn: c,
	}, nil
}

func (c *Conn) Read() ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var buffer bytes.Buffer
	for {
		tmp := make([]byte, 1024)
		n, err := c.Conn.Read(tmp)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		buffer.Write(tmp[:n])
	}
	return buffer.Bytes(), nil
}

func (c *Conn) Write(packet_id constants.ClientID, data []byte) (int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	b := []byte{}
	// length
	b = append(b, byte(len(data)+1))
	// packet_id
	b = append(b, byte(packet_id))
	// data
	b = append(b, data...)
	return c.Conn.Write(b)
}

func (c *Conn) Close() error {
	return c.Conn.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.Conn.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}
