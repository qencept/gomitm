package http1dump

import (
	"bytes"
	"fmt"
	"github.com/qencept/gomitm/pkg/mirror/http1"
	"github.com/qencept/gomitm/pkg/mirror/persistence"
	"io"
)

type Dump struct {
	path string
}

func (d *Dump) Inspect(p *http1.Parameters) error {
	c2s, err := persistence.CreateFile(persistence.CliSer, d.path, p.SessionParams)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintf(c2s, "%s %s %s%s\n\n", p.Req.Proto, p.Req.Method, p.Req.Host, p.Req.RequestURI); err != nil {
		return err
	}
	for k, v := range p.Req.Header {
		if _, err = fmt.Fprintf(c2s, "%s:%s\n", k, v); err != nil {
			return err
		}
	}
	if _, err = fmt.Fprintf(c2s, "\n"); err != nil {
		return err
	}
	if _, err = io.Copy(c2s, bytes.NewReader(p.ReqBytes)); err != nil {
		return err
	}
	_ = c2s.Close()

	s2c, err := persistence.CreateFile(persistence.SerCli, d.path, p.SessionParams)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintf(s2c, "%s %s\n\n", p.Resp.Proto, p.Resp.Status); err != nil {
		return err
	}
	for k, v := range p.Resp.Header {
		if _, err = fmt.Fprintf(s2c, "%s:%s\n", k, v); err != nil {
			return err
		}
	}
	if _, err = fmt.Fprintf(s2c, "\n"); err != nil {
		return err
	}
	if _, err = io.Copy(s2c, bytes.NewReader(p.RespBytes)); err != nil {
		return err
	}
	_ = s2c.Close()

	return nil
}

func New(path string) http1.Inspector {
	return &Dump{path: path}
}
