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

	trusted  *trusted.Trusted
	forging  *forging.Forging
	shuttler shuttle.Shuttle
}

func New(tcpClientConn, tcpServerConn *net.TCPConn, trusted *trusted.Trusted, forging *forging.Forging, shuttler shuttle.Shuttle) *Mitm {
	return &Mitm{tcpClientConn: tcpClientConn, tcpServerConn: tcpServerConn, trusted: trusted, forging: forging, shuttler: shuttler}
}

func (m *Mitm) Handle() {
	clientAddr := m.tcpClientConn.RemoteAddr()
	serverAddr := m.tcpServerConn.RemoteAddr()

	clientHello, tcpClientWrapper, err := sni.Detect(*m.tcpClientConn)
	if err != nil {
		logrus.Infof("TCP %v <-> %v \n", clientAddr, serverAddr)
		m.shuttler.Shuttle(tcpClientWrapper, m.tcpServerConn, "TCP")
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

		logrus.Infof("TLS %v (alpn=%v) <-> %v (alpn=%v) '%v'\n",
			clientAddr,
			tlsClientConn.ConnectionState().NegotiatedProtocol,
			serverAddr,
			tlsServerConn.ConnectionState().NegotiatedProtocol,
			clientHello.ServerName)
		m.shuttler.Shuttle(tlsClientConn, tlsServerConn, clientHello.ServerName)
	}
}
