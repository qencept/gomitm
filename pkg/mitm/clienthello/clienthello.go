package clienthello

import (
	"bytes"
	"crypto/tls"
	"github.com/qencept/gomitm/pkg/shuttle"
	"io"
)

func Detect(tcp shuttle.Connection) (shuttle.Connection, tls.ClientHelloInfo, bool) {
	var clientHelloInfo tls.ClientHelloInfo
	var ok bool
	savedBytes := &bytes.Buffer{}

	readOnly := ConnReadOnly{io.TeeReader(tcp, savedBytes)}
	_ = tls.Server(readOnly, &tls.Config{
		GetConfigForClient: func(info *tls.ClientHelloInfo) (*tls.Config, error) {
			clientHelloInfo = *info
			ok = true
			return nil, nil
		},
	}).Handshake()

	tcpMultiRead := &ConnectionMultiRead{multiReader: io.MultiReader(savedBytes, tcp), Connection: tcp}
	return tcpMultiRead, clientHelloInfo, ok
}
