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

type Http struct {
	dir        string
	inspectors []HttpInspector
}

func NewHttp(inspectors ...HttpInspector) SessionInspector {
	return &Http{"http", inspectors}
}

func (p *Http) Session(client, server shuttle.Stream, sni string) (io.WriteCloser, io.WriteCloser) {
	cpr, cpw := io.Pipe()
	spr, spw := io.Pipe()
	go func() {
		ts := strconv.Itoa(int(time.Now().Unix()))
		c2s, _ := os.Create(p.dir + "/" + ts + "[" + sni + "]" + client.RemoteAddr().String() + "->" + server.RemoteAddr().String())
		s2c, _ := os.Create(p.dir + "/" + ts + "[" + sni + "]" + client.RemoteAddr().String() + "<-" + server.RemoteAddr().String())

		for {
			req, err1 := http.ReadRequest(bufio.NewReader(cpr))
			if err1 == io.EOF {
				break
			}
			defer req.Body.Close()
			resp, err2 := http.ReadResponse(bufio.NewReader(spr), req)
			if err2 == io.EOF {
				break
			}
			defer resp.Body.Close()

			// should ne separate inspector
			fmt.Fprint(c2s, req.Host, req.RequestURI, "\n", req.Header)
			fmt.Fprint(s2c, resp.Status, "\n", resp.Header)

			for _, inspector := range p.inspectors {
				inspector.Http(req, resp)
			}
		}
	}()
	return cpw, spw
}
