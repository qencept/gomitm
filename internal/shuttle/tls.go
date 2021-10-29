package shuttle

import (
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
)

func OverTLS(tlsClientConn, tlsServerConn *tls.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		if _, err := io.Copy(tlsServerConn, tlsClientConn); err != nil {
			logrus.Warnf("Copy TLS %v -> %v: %v\n", tlsClientConn.RemoteAddr(), tlsServerConn.RemoteAddr(), err)
		}
		if err := tlsServerConn.CloseWrite(); err != nil {
			logrus.Warnf("CloseWrite TLS %v -> %v: %v\n", tlsClientConn.RemoteAddr(), tlsServerConn.RemoteAddr(), err)
		}
		wg.Done()
	}()
	go func() {
		if _, err := io.Copy(tlsClientConn, tlsServerConn); err != nil {
			logrus.Warnf("Copy TLS %v <- %v: %v\n", tlsClientConn.RemoteAddr(), tlsServerConn.RemoteAddr(), err)
		}
		if err := tlsClientConn.CloseWrite(); err != nil {
			logrus.Warnf("CloseWrite TLS %v <- %v: %v\n", tlsClientConn.RemoteAddr(), tlsServerConn.RemoteAddr(), err)
		}
		wg.Done()
	}()
	wg.Wait()
}
