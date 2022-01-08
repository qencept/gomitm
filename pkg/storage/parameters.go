package storage

import (
	"crypto/tls"
	"github.com/qencept/gomitm/pkg/shuttler"
)

const (
	Forward int = iota + 1
	Backward
)

type Parameters struct {
	Client string
	Server string
	Sni    string
}

func NewParameters(client, server shuttler.Connection) *Parameters {
	sni := ""
	if conn, ok := server.(*tls.Conn); ok {
		sni = conn.ConnectionState().ServerName
	}
	return &Parameters{
		Client: client.RemoteAddr().String(),
		Server: server.RemoteAddr().String(),
		Sni:    sni,
	}
}
