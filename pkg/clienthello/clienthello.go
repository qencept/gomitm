package clienthello

import (
	"bytes"
	"crypto/tls"
	"github.com/qencept/gomitm/pkg/shuttler"
	"io"
)

func Detect(tcp shuttler.Connection) (shuttler.Connection, tls.ClientHelloInfo, bool) {
	var clientHelloInfo tls.ClientHelloInfo
	var ok bool
	savedBytes := &bytes.Buffer{}

	readOnly := connReadOnly{io.TeeReader(tcp, savedBytes)}
	_ = tls.Server(readOnly, &tls.Config{
		GetConfigForClient: func(info *tls.ClientHelloInfo) (*tls.Config, error) {
			clientHelloInfo = *info
			ok = true
			return nil, nil
		},
	}).Handshake()

	tcpMultiRead := &connectionMultiRead{multiReader: io.MultiReader(savedBytes, tcp), Connection: tcp}
	return tcpMultiRead, clientHelloInfo, ok
}
