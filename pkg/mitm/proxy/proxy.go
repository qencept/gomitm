package proxy

import (
	"fmt"
	"github.com/qencept/gomitm/pkg/mitm/destination"
	"github.com/qencept/gomitm/pkg/mitm/forgery"
	"github.com/qencept/gomitm/pkg/mitm/shuttle"
	"github.com/qencept/gomitm/pkg/mitm/trusted"
	"github.com/sirupsen/logrus"
	"net"
)

type Proxy struct {
	trusted *trusted.Trusted
	forgery *forgery.Forgery
	shuttle shuttle.Shuttle
	addr    string
}

func New(addr string, trusted *trusted.Trusted, forgery *forgery.Forgery, shuttle shuttle.Shuttle) *Proxy {
	return &Proxy{trusted: trusted, forgery: forgery, addr: addr, shuttle: shuttle}
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
		go func() {
			defer tcpClientConn.Close()
			p.Handle(tcpClientConn)
		}()
	}
}

func (p *Proxy) Handle(tcpClientConn *net.TCPConn) {
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

	Mitm(tcpClientConn, tcpServerConn, p.trusted, p.forgery, p.shuttle)
}
