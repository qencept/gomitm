package proxy

import (
	"crypto/tls"
	"github.com/qencept/gomitm/pkg/mitm/forgery"
	"github.com/qencept/gomitm/pkg/mitm/shuttle"
	"github.com/qencept/gomitm/pkg/mitm/sni"
	"github.com/qencept/gomitm/pkg/mitm/trusted"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

func Mitm(tcpClientConn, tcpServerConn *net.TCPConn, t *trusted.Trusted, f *forgery.Forgery, s shuttle.Shuttle) {
	clientHello, tcpClientWrapper, ok := sni.Detect(*tcpClientConn)
	if !ok {
		session := &shuttle.Session{Client: tcpClientWrapper, Server: tcpServerConn}
		logrus.Infof("TCP %v <-> %v \n", session.Client.RemoteAddr(), session.Server.RemoteAddr())
		s.Shuttle(session)
		return
	}

	tlsServerConn := tls.Client(tcpServerConn, &tls.Config{
		ServerName: clientHello.ServerName,
		NextProtos: clientHello.SupportedProtos,
		RootCAs:    t.CertPool(),
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
			return f.Forge(tlsServerConn.ConnectionState().PeerCertificates[0])
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

	session := &shuttle.Session{Client: tlsClientConn, Server: tlsServerConn}
	logrus.Infof("TLS %v (alpn=%v) <-> %v (alpn=%v) '%v'\n",
		session.Client.RemoteAddr(),
		tlsClientConn.ConnectionState().NegotiatedProtocol,
		session.Server.RemoteAddr(),
		tlsServerConn.ConnectionState().NegotiatedProtocol,
		clientHello.ServerName)
	s.Shuttle(session)

}
