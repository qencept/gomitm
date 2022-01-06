package clienthello

import (
	"crypto/tls"
	"github.com/qencept/gomitm/pkg/backup"
	"github.com/qencept/gomitm/pkg/shuttler"
)

func Detect(tcp shuttler.Connection) (shuttler.Connection, tls.ClientHelloInfo, bool) {
	var clientHelloInfo tls.ClientHelloInfo
	var ok bool
	backupReader := backup.NewReader(tcp)
	_ = tls.Server(conn{Reader: backupReader}, &tls.Config{
		GetConfigForClient: func(info *tls.ClientHelloInfo) (*tls.Config, error) {
			clientHelloInfo = *info
			ok = true
			return nil, nil
		},
	}).Handshake()
	backupReader.Reset()
	return &connection{Connection: tcp, reader: backupReader}, clientHelloInfo, ok
}
