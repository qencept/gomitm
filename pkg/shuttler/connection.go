package shuttler

import "net"

type Connection interface {
	net.Conn
	CloseWrite() error
}
