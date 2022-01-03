package shuttle

import "net"

type Connection interface {
	net.Conn
	CloseWrite() error
}

type Shuttle interface {
	Shuttle(client, server Connection)
}
