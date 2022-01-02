package shuttle

import (
	"net"
)

type Session struct {
	Client net.Conn
	Server net.Conn
}

type Shuttle interface {
	Shuttle(session *Session)
}
