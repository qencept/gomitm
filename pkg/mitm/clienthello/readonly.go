package clienthello

import (
	"io"
	"net"
	"time"
)

type ConnReadOnly struct {
	io.Reader
}

func (c ConnReadOnly) Read(p []byte) (int, error)       { return c.Reader.Read(p) }
func (c ConnReadOnly) Write([]byte) (int, error)        { return 0, io.ErrClosedPipe }
func (c ConnReadOnly) Close() error                     { return nil }
func (c ConnReadOnly) LocalAddr() net.Addr              { return nil }
func (c ConnReadOnly) RemoteAddr() net.Addr             { return nil }
func (c ConnReadOnly) SetDeadline(time.Time) error      { return nil }
func (c ConnReadOnly) SetReadDeadline(time.Time) error  { return nil }
func (c ConnReadOnly) SetWriteDeadline(time.Time) error { return nil }

var _ net.Conn = (*ConnReadOnly)(nil)
