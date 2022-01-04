package clienthello

import (
	"io"
	"net"
	"time"
)

type connReadOnly struct {
	io.Reader
}

func (c connReadOnly) Read(p []byte) (int, error)       { return c.Reader.Read(p) }
func (c connReadOnly) Write([]byte) (int, error)        { return 0, io.ErrClosedPipe }
func (c connReadOnly) Close() error                     { return nil }
func (c connReadOnly) LocalAddr() net.Addr              { return nil }
func (c connReadOnly) RemoteAddr() net.Addr             { return nil }
func (c connReadOnly) SetDeadline(time.Time) error      { return nil }
func (c connReadOnly) SetReadDeadline(time.Time) error  { return nil }
func (c connReadOnly) SetWriteDeadline(time.Time) error { return nil }

var _ net.Conn = (*connReadOnly)(nil)
