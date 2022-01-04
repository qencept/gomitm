package proxy

import (
	"fmt"
	"github.com/qencept/gomitm/pkg/forgery"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/mitm"
	"github.com/qencept/gomitm/pkg/shuttler"
	"github.com/qencept/gomitm/pkg/trusted"
	"net"
)

type Proxy struct {
	addr   string
	mitm   *mitm.Mitm
	logger logger.Logger
}

func New(addr string, t *trusted.Trusted, f *forgery.Forgery, s shuttler.Shuttler, l logger.Logger) *Proxy {
	m := mitm.New(t, f, s, l)
	return &Proxy{logger: l, mitm: m, addr: addr}
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
	p.logger.Infoln("Listening", tcpAddr)
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
