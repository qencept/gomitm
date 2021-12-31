package sni

import (
	"io"
	"net"
)

type ConnWrapper struct {
	net.TCPConn

	r io.Reader
}

func (c ConnWrapper) Read(p []byte) (int, error) { return c.r.Read(p) }
