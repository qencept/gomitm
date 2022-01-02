package sni

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
)

func Detect(tcpClientConn net.TCPConn) (*tls.ClientHelloInfo, net.Conn, bool) {
	var clientHello *tls.ClientHelloInfo
	savedBytes := &bytes.Buffer{}

	_ = tls.Server(ConnReader{io.TeeReader(&tcpClientConn, savedBytes)}, &tls.Config{
		GetConfigForClient: func(info *tls.ClientHelloInfo) (*tls.Config, error) {
			clientHello = info
			return nil, nil
		},
	}).Handshake()

	tcpClientReader := io.MultiReader(savedBytes, &tcpClientConn)
	return clientHello, &ConnWrapper{r: tcpClientReader, TCPConn: tcpClientConn}, clientHello != nil
}
