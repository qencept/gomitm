package clienthello

import (
	"io"
	"net"
	"time"
)

type conn struct {
	io.Reader
}

func (c conn) Read(p []byte) (int, error)       { return c.Reader.Read(p) }
func (c conn) Write([]byte) (int, error)        { return 0, io.ErrClosedPipe }
func (c conn) Close() error                     { return nil }
func (c conn) LocalAddr() net.Addr              { return nil }
func (c conn) RemoteAddr() net.Addr             { return nil }
func (c conn) SetDeadline(time.Time) error      { return nil }
func (c conn) SetReadDeadline(time.Time) error  { return nil }
func (c conn) SetWriteDeadline(time.Time) error { return nil }

var _ net.Conn = (*conn)(nil)
