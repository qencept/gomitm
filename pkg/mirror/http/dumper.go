package httpi

import (
	"crypto/tls"
	"fmt"
	"github.com/qencept/gomitm/pkg/mitm/shuttle"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Dumper struct {
	dir string
}

func (d *Dumper) Inspect(req *http.Request, resp *http.Response, session *shuttle.Session) {
	sni := ""
	if conn, ok := session.Server.(*tls.Conn); ok {
		sni = conn.ConnectionState().ServerName
	}
	ts := strconv.Itoa(int(time.Now().Unix()))
	c2s, _ := os.Create(d.dir + "/" + ts + "[" + sni + "]" + session.Client.RemoteAddr().String() + "->" + session.Server.RemoteAddr().String())
	s2c, _ := os.Create(d.dir + "/" + ts + "[" + sni + "]" + session.Client.RemoteAddr().String() + "<-" + session.Server.RemoteAddr().String())

	fmt.Fprint(c2s, req.Host, req.RequestURI, "\n", req.Header)
	fmt.Fprint(s2c, resp.Status, "\n", resp.Header)

	c2s.Close()
	s2c.Close()
}

func NewDumper(dir string) Inspector {
	return &Dumper{dir: dir}
}
