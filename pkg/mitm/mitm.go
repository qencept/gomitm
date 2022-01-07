package mitm

import (
	"crypto/tls"
	"fmt"
	"github.com/qencept/gomitm/pkg/clienthello"
	"github.com/qencept/gomitm/pkg/forgery"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/shuttler"
	"github.com/qencept/gomitm/pkg/trusted"
	"time"
)

type Mitm struct {
	trusted  *trusted.Trusted
	forgery  *forgery.Forgery
	shuttler shuttler.Shuttler
	logger   logger.Logger
}

func New(t *trusted.Trusted, f *forgery.Forgery, s shuttler.Shuttler, l logger.Logger) *Mitm {
	return &Mitm{trusted: t, forgery: f, shuttler: s, logger: l}
}

func (m *Mitm) Run(tcpOrigClient, tcpOrigServer shuttler.Connection) {
	tcpOrigClient, clientHelloInfo, ok := clienthello.Detect(tcpOrigClient)
	if !ok {
		m.logger.Infoln("TCP", tcpOrigClient.RemoteAddr(), "<->", tcpOrigServer.RemoteAddr())
		m.shuttler.Shuttle(tcpOrigClient, tcpOrigServer)
		return
	}

	tlsOrigServer := tls.Client(tcpOrigServer, &tls.Config{
		ServerName: clientHelloInfo.ServerName,
		NextProtos: clientHelloInfo.SupportedProtos,
		RootCAs:    m.trusted.CertPool(),
	})
	defer func() {
		_ = tlsOrigServer.Close()
	}()
	if err := tlsOrigServer.SetDeadline(time.Now().Add(time.Minute)); err != nil {
		m.logger.Warnln("Server SetDeadline:", err)
		return
	}
	if err := tlsOrigServer.Handshake(); err != nil {
		m.logger.Warnln("Server Handshake:", err)
		return
	}

	tlsOrigClient := tls.Server(tcpOrigClient, &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return m.forgery.Forge(tlsOrigServer.ConnectionState().PeerCertificates[0])
		},
		NextProtos: []string{tlsOrigServer.ConnectionState().NegotiatedProtocol},
	})
	defer func() {
		_ = tlsOrigClient.Close()
	}()
	if err := tlsOrigClient.SetDeadline(time.Now().Add(time.Minute)); err != nil {
		m.logger.Warnln("Client SetDeadline:", err)
		return
	}
	if err := tlsOrigClient.Handshake(); err != nil {
		m.logger.Warnln("Client Handshake:", err)
		return
	}

	m.logger.Infoln(fmt.Sprintf("TLS %v (alpn=%v) <-> %v (alpn=%v) '%v'\n",
		tlsOrigClient.RemoteAddr(),
		tlsOrigClient.ConnectionState().NegotiatedProtocol,
		tlsOrigServer.RemoteAddr(),
		tlsOrigServer.ConnectionState().NegotiatedProtocol,
		clientHelloInfo.ServerName))
	m.shuttler.Shuttle(tlsOrigClient, tlsOrigServer)
}
