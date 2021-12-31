package mirror

import (
	"bufio"
	"fmt"
	"github.com/qencept/gomitm/internal/shuttle"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type HttpParser struct {
	dir string
}

func NewHttpParser() Inspector {
	return &HttpParser{"http"}
}

func (p *HttpParser) Session(client, server shuttle.Stream, sni string) (io.WriteCloser, io.WriteCloser) {
	cpr, cpw := io.Pipe()
	spr, spw := io.Pipe()
	go func() {
		ts := strconv.Itoa(int(time.Now().Unix()))
		c2s, _ := os.Create(p.dir + "/" + ts + "[" + sni + "]" + client.RemoteAddr().String() + "->" + server.RemoteAddr().String())
		s2c, _ := os.Create(p.dir + "/" + ts + "[" + sni + "]" + client.RemoteAddr().String() + "<-" + server.RemoteAddr().String())

		req, _ := http.ReadRequest(bufio.NewReader(cpr))
		resp, _ := http.ReadResponse(bufio.NewReader(spr), req)
		fmt.Fprint(c2s, req.Host, req.RequestURI, "\n", req.Header)
		fmt.Fprint(s2c, resp.Status, "\n", resp.Header)
		resp.Body.Close()
	}()
	return cpw, spw
}
