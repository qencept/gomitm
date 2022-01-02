package httpi

import (
	"bufio"
	"github.com/qencept/gomitm/pkg/mirror/session"
	"github.com/qencept/gomitm/pkg/mitm/shuttle"
	"io"
	"net/http"
)

type Http struct {
	dir        string
	inspectors []Inspector
}

func NewHttp(inspectors ...Inspector) session.Inspector {
	return &Http{"http", inspectors}
}

func (p *Http) Inspect(session *shuttle.Session) (io.WriteCloser, io.WriteCloser) {
	cpr, cpw := io.Pipe()
	spr, spw := io.Pipe()
	go func() {
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

			for _, inspector := range p.inspectors {
				inspector.Inspect(req, resp, session)
			}
		}
	}()
	return cpw, spw
}
