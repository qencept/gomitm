package mitm

import (
	"crypto/tls"
	"github.com/qencept/gomitm/internal/forging"
	"github.com/qencept/gomitm/internal/shuttle"
	"github.com/qencept/gomitm/internal/sni"
	"github.com/qencept/gomitm/internal/trusted"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type Mitm struct {
	tcpClientConn *net.TCPConn
	tcpServerConn *net.TCPConn

	trusted *trusted.Trusted
	forging *forging.Forging
}

func New(tcpClientConn, tcpServerConn *net.TCPConn, trusted *trusted.Trusted, forging *forging.Forging) *Mitm {
	return &Mitm{tcpClientConn: tcpClientConn, tcpServerConn: tcpServerConn, trusted: trusted, forging: forging}
}

func (m *Mitm) Handle() {
	clientAddr := m.tcpClientConn.RemoteAddr()
	serverAddr := m.tcpServerConn.RemoteAddr()

	clientHello, tcpClientWrapper, err := sni.Detect(m.tcpClientConn)
	if err != nil {
		logrus.Infoln("Assumed raw TCP")
		shuttle.OverTCP(tcpClientWrapper, m.tcpServerConn)
	} else {
		tlsServerConn := tls.Client(m.tcpServerConn, &tls.Config{
			ServerName: clientHello.ServerName,
			NextProtos: clientHello.SupportedProtos,
			RootCAs:    m.trusted.CertPool(),
		})
		defer tlsServerConn.Close()
		if err := tlsServerConn.SetDeadline(time.Now().Add(time.Minute)); err != nil {
			logrus.Warnln("Server SetDeadline:", err)
			return
		}
		if err := tlsServerConn.Handshake(); err != nil {
			logrus.Warnln("Server Handshake:", err)
			return
		}

		tlsClientConn := tls.Server(tcpClientWrapper, &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				return m.forging.Forge(tlsServerConn.ConnectionState().PeerCertificates[0])
			},
			NextProtos: []string{tlsServerConn.ConnectionState().NegotiatedProtocol},
		})
		defer tlsClientConn.Close()
		if err := tlsClientConn.SetDeadline(time.Now().Add(time.Minute)); err != nil {
			logrus.Warnln("Client SetDeadline:", err)
			return
		}
		if err := tlsClientConn.Handshake(); err != nil {
			logrus.Warnln("Client Handshake:", err)
			return
		}

		logrus.Infof("MITM %v (v=0x0%x cs=0x%x alpn=%v) <-> %v (v=0x0%x cs=0x%x alpn=%v) '%v'\n",
			clientAddr,
			tlsClientConn.ConnectionState().Version,
			tlsClientConn.ConnectionState().CipherSuite,
			tlsClientConn.ConnectionState().NegotiatedProtocol,
			serverAddr,
			tlsServerConn.ConnectionState().Version,
			tlsServerConn.ConnectionState().CipherSuite,
			tlsServerConn.ConnectionState().NegotiatedProtocol,
			clientHello.ServerName)
		shuttle.OverTLS(tlsClientConn, tlsServerConn)
	}
}
