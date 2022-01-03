package proxy

import (
	"fmt"
	"github.com/qencept/gomitm/pkg/mitm/destination"
	"github.com/qencept/gomitm/pkg/mitm/forgery"
	"github.com/qencept/gomitm/pkg/mitm/trusted"
	"github.com/qencept/gomitm/pkg/shuttle"
	"github.com/sirupsen/logrus"
	"net"
)

type Proxy struct {
	logger *logrus.Logger
	attack *Attack
	addr   string
}

func New(addr string, trusted *trusted.Trusted, forgery *forgery.Forgery, logger *logrus.Logger, shuttle shuttle.Shuttle) *Proxy {
	attack := &Attack{logger: logger, trusted: trusted, forgery: forgery, shuttle: shuttle}
	return &Proxy{logger: logger, attack: attack, addr: addr}
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
	p.logger.Infof("Listening %v \n", tcpAddr)

	for {
		tcpClientConn, err := listener.AcceptTCP()
		if err != nil {
			return fmt.Errorf("accepting: %w", err)
		}
		go func() {
			defer func() { _ = tcpClientConn.Close() }()
			p.Handle(tcpClientConn)
		}()
	}
}

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

	p.attack.Run(tcpClientConn, tcpServerConn)
}
