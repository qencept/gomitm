package doh

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/qencept/gomitm/pkg/mirror/http"
	"github.com/qencept/gomitm/pkg/mitm/shuttle"
	"golang.org/x/net/dns/dnsmessage"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Doh struct {
	dir string
}

func NewDoh(dir string) httpi.Inspector {
	return &Doh{dir: dir}
}

func (d *Doh) Inspect(req *http.Request, resp *http.Response, session *shuttle.Session) {
	sni := ""
	if conn, ok := session.Server.(*tls.Conn); ok {
		sni = conn.ConnectionState().ServerName
	}
	ts := strconv.Itoa(int(time.Now().Unix()))

	buf, msg := bytes.Buffer{}, dnsmessage.Message{}
	buf.ReadFrom(req.Body)
	err := msg.Unpack(buf.Bytes())
	if err == nil {
		c2s, _ := os.Create(d.dir + "/" + ts + "[" + sni + "]" + session.Client.RemoteAddr().String() + "->" + session.Server.RemoteAddr().String())
		defer c2s.Close()
		for _, q := range msg.Questions {
			fmt.Fprintln(c2s, q)
		}
	}

	buf.Reset()
	buf.ReadFrom(resp.Body)
	err = msg.Unpack(buf.Bytes())
	if err == nil {
		s2c, _ := os.Create(d.dir + "/" + ts + "[" + sni + "]" + session.Client.RemoteAddr().String() + "<-" + session.Server.RemoteAddr().String())
		defer s2c.Close()
		for _, a := range msg.Answers {
			fmt.Fprintln(s2c, a.Header.Name, a.Header.Type, a.Body)
		}
	}
}
