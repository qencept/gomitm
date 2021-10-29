package sni

import (
	"io"
	"net"
	"time"
)

type ConnWrapper struct {
	r io.Reader
	c net.Conn
}

func (c ConnWrapper) Read(p []byte) (int, error)         { return c.r.Read(p) }
func (c ConnWrapper) Write(p []byte) (int, error)        { return c.c.Write(p) }
func (c ConnWrapper) Close() error                       { return c.c.Close() }
func (c ConnWrapper) LocalAddr() net.Addr                { return c.c.LocalAddr() }
func (c ConnWrapper) RemoteAddr() net.Addr               { return c.c.RemoteAddr() }
func (c ConnWrapper) SetDeadline(t time.Time) error      { return c.c.SetDeadline(t) }
func (c ConnWrapper) SetReadDeadline(t time.Time) error  { return c.c.SetReadDeadline(t) }
func (c ConnWrapper) SetWriteDeadline(t time.Time) error { return c.c.SetDeadline(t) }
