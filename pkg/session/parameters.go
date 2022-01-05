package session

import (
	"crypto/tls"
	"github.com/qencept/gomitm/pkg/shuttler"
	"net"
)

const (
	Forward int = iota + 1
	Backward
)

type Parameters struct {
	Client net.Addr
	Server net.Addr
	Sni    string
}

func NewParameters(client, server shuttler.Connection) *Parameters {
	sni := ""
	if conn, ok := server.(*tls.Conn); ok {
		sni = conn.ConnectionState().ServerName
	}
	return &Parameters{
		Client: client.RemoteAddr(),
		Server: server.RemoteAddr(),
		Sni:    sni,
	}
}
