package shuttle

import (
	"io"
	"net"
)

type Stream interface {
	io.Reader
	io.Writer

	RemoteAddr() net.Addr
	CloseWrite() error
}

type Shuttle interface {
	Shuttle(client Stream, server Stream, sni string) error
}
