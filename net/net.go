package net

import (
	"net"
	"time"

	"github.com/nanachi-sh/go-mc/constants"
)

type Conn struct {
	net.Conn
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

func (c *Conn) Read(b []byte) (int, error) {
	return c.Conn.Read(b)
}

func (c *Conn) WriteRaw(b []byte) (int, error) {
	return c.Conn.Write(b)
}

func (c *Conn) Write(packet_id constants.ClientID, data []byte) (int, error) {
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
