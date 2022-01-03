package doh

import (
	"fmt"
	"github.com/qencept/gomitm/pkg/mirror/http1"
	"github.com/qencept/gomitm/pkg/mirror/persistence"
	"golang.org/x/net/dns/dnsmessage"
	"net"
)

type Doh struct {
	path string
}

func New(path string) http1.Inspector {
	return &Doh{path: path}
}

func (d *Doh) Inspect(params *http1.Parameters) error {
	msg := dnsmessage.Message{}

	if msg.Unpack(params.ReqBytes) == nil {
		c2s, err := persistence.CreateFile(persistence.CliSer, d.path, params.SessionParams)
		if err != nil {
			return err
		}
		defer func() { _ = c2s.Close() }()
		for _, q := range msg.Questions {
			_, err = fmt.Fprintln(c2s, q.Name, q.Type)
			if err != nil {
				return err
			}
		}
	}

	if msg.Unpack(params.RespBytes) == nil {
		s2c, err := persistence.CreateFile(persistence.SerCli, d.path, params.SessionParams)
		if err != nil {
			return err
		}
		defer func() { _ = s2c.Close() }()
		for _, a := range msg.Answers {
			var str string
			switch b := a.Body.(type) {
			case *dnsmessage.AResource:
				str = net.IPv4(b.A[0], b.A[1], b.A[2], b.A[3]).String()
			case *dnsmessage.CNAMEResource:
				str = b.CNAME.String()
			default:
				str = b.GoString()
			}
			_, err = fmt.Fprintln(s2c, a.Header.Name, a.Header.Type, str)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
