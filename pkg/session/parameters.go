package session

import (
	"crypto/tls"
	"github.com/qencept/gomitm/pkg/shuttler"
	"strconv"
)

type Direction int

const (
	Forward Direction = iota + 1
	Backward
)

type Parameters struct {
	Client    string
	Server    string
	Sni       string
	Timestamp string
}

func NewParameters(client, server shuttler.Connection, ts int64) *Parameters {
	sni := ""
	if conn, ok := server.(*tls.Conn); ok {
		sni = conn.ConnectionState().ServerName
	}
	return &Parameters{
		Client:    client.RemoteAddr().String(),
		Server:    server.RemoteAddr().String(),
		Sni:       sni,
		Timestamp: strconv.FormatInt(ts, 10),
	}
}
