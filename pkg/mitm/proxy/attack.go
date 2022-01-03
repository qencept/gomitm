package proxy

import (
	"crypto/tls"
	"github.com/qencept/gomitm/pkg/mitm/clienthello"
	"github.com/qencept/gomitm/pkg/mitm/forgery"
	"github.com/qencept/gomitm/pkg/mitm/trusted"
	"github.com/qencept/gomitm/pkg/shuttle"
	"github.com/sirupsen/logrus"
	"time"
)

type Attack struct {
	logger  *logrus.Logger
	trusted *trusted.Trusted
	forgery *forgery.Forgery
	shuttle shuttle.Shuttle
}

func hack(original []string) []string {
	for _, a := range original {
		if a == "h2" {
			return []string{"http/1.1"}
		}
	}
	return original
}

func (a *Attack) Run(tcpOrigClient, tcpOrigServer shuttle.Connection) {
	tcpOrigClient, clientHelloInfo, ok := clienthello.Detect(tcpOrigClient)
	if !ok {
		a.logger.Infof("TCP %v <-> %v \n", tcpOrigClient.RemoteAddr(), tcpOrigServer.RemoteAddr())
		a.shuttle.Shuttle(tcpOrigClient, tcpOrigServer)
		return
	}

	tlsOrigServer := tls.Client(tcpOrigServer, &tls.Config{
		ServerName: clientHelloInfo.ServerName,
		NextProtos: hack(clientHelloInfo.SupportedProtos),
		RootCAs:    a.trusted.CertPool(),
	})
	defer func() { _ = tlsOrigServer.Close() }()
	if err := tlsOrigServer.SetDeadline(time.Now().Add(time.Minute)); err != nil {
		a.logger.Warnln("Server SetDeadline:", err)
		return
	}
	if err := tlsOrigServer.Handshake(); err != nil {
		a.logger.Warnln("Server Handshake:", err)
		return
	}

	tlsOrigClient := tls.Server(tcpOrigClient, &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return a.forgery.Forge(tlsOrigServer.ConnectionState().PeerCertificates[0])
		},
		NextProtos: []string{tlsOrigServer.ConnectionState().NegotiatedProtocol},
	})
	defer func() { _ = tlsOrigClient.Close() }()
	if err := tlsOrigClient.SetDeadline(time.Now().Add(time.Minute)); err != nil {
		a.logger.Warnln("Client SetDeadline:", err)
		return
	}
	if err := tlsOrigClient.Handshake(); err != nil {
		a.logger.Warnln("Client Handshake:", err)
		return
	}

	a.logger.Infof("TLS %v (alpn=%v) <-> %v (alpn=%v) '%v'\n",
		tlsOrigClient.RemoteAddr(),
		tlsOrigClient.ConnectionState().NegotiatedProtocol,
		tlsOrigServer.RemoteAddr(),
		tlsOrigServer.ConnectionState().NegotiatedProtocol,
		clientHelloInfo.ServerName)
	a.shuttle.Shuttle(tlsOrigClient, tlsOrigServer)
}
