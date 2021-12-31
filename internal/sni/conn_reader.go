package sni

import (
	"io"
	"net"
	"time"
)

type ConnReader struct {
	io.Reader
}

func (c ConnReader) Read(p []byte) (int, error)         { return c.Reader.Read(p) }
func (c ConnReader) Write([]byte) (int, error)          { return 0, io.ErrClosedPipe }
func (c ConnReader) Close() error                       { return nil }
func (c ConnReader) LocalAddr() net.Addr                { return nil }
func (c ConnReader) RemoteAddr() net.Addr               { return nil }
func (c ConnReader) SetDeadline(t time.Time) error      { return nil }
func (c ConnReader) SetReadDeadline(t time.Time) error  { return nil }
func (c ConnReader) SetWriteDeadline(t time.Time) error { return nil }
