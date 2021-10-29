package shuttle

import (
	"github.com/qencept/gomitm/internal/sni"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
)

func OverTCP(tcpClientWrapper *sni.ConnWrapper, tcpServerConn net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		if _, err := io.Copy(tcpClientWrapper, tcpServerConn); err != nil {
			logrus.Warnf("Copy TCP %v -> %v: %v\n", tcpClientWrapper.RemoteAddr(), tcpServerConn.RemoteAddr(), err)
		}
		if err := tcpClientWrapper.Close(); err != nil {
			logrus.Warnf("Close TCP %v -> %v: %v\n", tcpClientWrapper.RemoteAddr(), tcpServerConn.RemoteAddr(), err)
		}
		wg.Done()
	}()
	go func() {
		if _, err := io.Copy(tcpServerConn, tcpClientWrapper); err != nil {
			logrus.Warnf("Copy TCP %v <- %v: %v\n", tcpClientWrapper.RemoteAddr(), tcpServerConn.RemoteAddr(), err)
		}
		if err := tcpServerConn.(*net.TCPConn).CloseWrite(); err != nil {
			logrus.Warnf("CloseWrite TCP %v <- %v: %v\n", tcpClientWrapper.RemoteAddr(), tcpServerConn.RemoteAddr(), err)
		}
		wg.Done()
	}()
	wg.Wait()
}
