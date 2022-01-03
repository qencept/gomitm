package sessiondump

import (
	"github.com/qencept/gomitm/pkg/mirror/persistence"
	"github.com/qencept/gomitm/pkg/mirror/session"
	"io"
)

type Dump struct {
	path string
}

func (d *Dump) InitWriteClosers(params *session.Parameters) (c2s, s2c io.WriteCloser, err error) {
	c2s, err = persistence.CreateFile(persistence.CliSer, d.path, params)
	if err != nil {
		return
	}
	s2c, err = persistence.CreateFile(persistence.SerCli, d.path, params)
	if err != nil {
		return
	}
	return
}

func New(path string) session.Inspector {
	return &Dump{path: path}
}
