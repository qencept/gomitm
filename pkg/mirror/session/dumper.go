package session

import (
	"crypto/tls"
	"github.com/qencept/gomitm/pkg/mitm/shuttle"
	"io"
	"os"
	"strconv"
	"time"
)

type Dumper struct {
	dir string
}

func (d *Dumper) Inspect(session *shuttle.Session) (io.WriteCloser, io.WriteCloser) {
	sni := ""
	if conn, ok := session.Server.(*tls.Conn); ok {
		sni = conn.ConnectionState().ServerName
	}
	ts := strconv.Itoa(int(time.Now().Unix()))
	c2s, _ := os.Create(d.dir + "/" + ts + "[" + sni + "]" + session.Client.RemoteAddr().String() + "->" + session.Server.RemoteAddr().String())
	s2c, _ := os.Create(d.dir + "/" + ts + "[" + sni + "]" + session.Client.RemoteAddr().String() + "<-" + session.Server.RemoteAddr().String())
	return c2s, s2c
}

func NewDumper(dir string) Inspector {
	return &Dumper{dir: dir}
}
