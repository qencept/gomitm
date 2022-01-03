package session

import (
	"crypto/tls"
	"github.com/qencept/gomitm/pkg/shuttle"
	"net"
)

type Parameters struct {
	ClientAddr net.Addr
	ServerAddr net.Addr
	Sni        string
}

func NewParameters(client, server shuttle.Connection) *Parameters {
	sni := ""
	if conn, ok := server.(*tls.Conn); ok {
		sni = conn.ConnectionState().ServerName
	}
	return &Parameters{
		ClientAddr: client.RemoteAddr(),
		ServerAddr: server.RemoteAddr(),
		Sni:        sni,
	}
}
