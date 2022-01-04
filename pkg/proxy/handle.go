package proxy

import (
	"github.com/qencept/gomitm/pkg/destination"
	"net"
)

func (p *Proxy) Handle(tcpClientConn *net.TCPConn) {
	serverAddr, err := destination.Detect(tcpClientConn)
	if err != nil {
		p.logger.Warnln("Destination detect:", err)
		return
	}
	tcpServerConn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		p.logger.Warnln("Dial Server:", err)
		return
	}
	defer func() { _ = tcpServerConn.Close() }()
	p.mitm.Run(tcpClientConn, tcpServerConn)
}
