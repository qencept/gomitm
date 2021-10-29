package sni

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
)

func Detect(tcpClientConn net.Conn) (*tls.ClientHelloInfo, *ConnWrapper, error) {
	var clientHello *tls.ClientHelloInfo
	savedBytes := new(bytes.Buffer)

	_ = tls.Server(ConnReader{r: io.TeeReader(tcpClientConn, savedBytes)}, &tls.Config{
		GetConfigForClient: func(info *tls.ClientHelloInfo) (*tls.Config, error) {
			clientHello = new(tls.ClientHelloInfo)
			clientHello = info
			return nil, nil
		},
	}).Handshake()

	var err error
	if clientHello == nil {
		err = fmt.Errorf("no ClientHello")
	}

	tcpClientConnReader := io.MultiReader(savedBytes, tcpClientConn)
	return clientHello, &ConnWrapper{r: tcpClientConnReader, c: tcpClientConn}, err
}
