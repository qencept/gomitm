package proxy

import (
	"fmt"
	"github.com/qencept/gomitm/internal/destination"
	"github.com/qencept/gomitm/internal/forging"
	"github.com/qencept/gomitm/internal/mitm"
	"github.com/qencept/gomitm/internal/shuttle"
	"github.com/qencept/gomitm/internal/trusted"
	"github.com/sirupsen/logrus"
	"net"
)

type Proxy struct {
	trusted  *trusted.Trusted
	forging  *forging.Forging
	shuttler shuttle.Shuttle
	addr     string
}

func New(addr string, trusted *trusted.Trusted, forging *forging.Forging, shuttler shuttle.Shuttle) *Proxy {
	return &Proxy{trusted: trusted, forging: forging, addr: addr, shuttler: shuttler}
}

func (p *Proxy) Run() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", p.addr)
	if err != nil {
		return fmt.Errorf("resolving addr: %w", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return fmt.Errorf("listening: %w", err)
	}

	for {
		tcpClientConn, err := listener.AcceptTCP()
		if err != nil {
			return fmt.Errorf("accepting: %w", err)
		}

		go p.Handle(tcpClientConn)
	}
}

func (p *Proxy) Handle(tcpClientConn *net.TCPConn) {
	defer tcpClientConn.Close()

	serverAddr, err := destination.Detect(tcpClientConn)
	if err != nil {
		logrus.Warnln("Destination detect:", err)
		return
	}

	tcpServerConn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		logrus.Warnln("Dial Server:", err)
		return
	}
	defer tcpServerConn.Close()

	mitm.New(tcpClientConn, tcpServerConn, p.trusted, p.forging, p.shuttler).Handle()
}
