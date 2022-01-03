package http1

import (
	"bufio"
	"github.com/qencept/gomitm/pkg/mirror/session"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type Http struct {
	logger     *logrus.Logger
	inspectors []Inspector
}

func New(logger *logrus.Logger, inspectors ...Inspector) session.Inspector {
	return &Http{logger: logger, inspectors: inspectors}
}

func (h *Http) InitWriteClosers(params *session.Parameters) (io.WriteCloser, io.WriteCloser, error) {
	cpr, cpw := io.Pipe()
	spr, spw := io.Pipe()
	go func() {
		for {
			req, err := http.ReadRequest(bufio.NewReader(cpr))
			if err != nil {
				if err != io.EOF {
					h.logger.Warnf("http.ReadRequest %v -> %v: %v", params.ClientAddr, params.ServerAddr, err)
				}
				return
			}
			defer func() { _ = req.Body.Close() }()

			resp, err := http.ReadResponse(bufio.NewReader(spr), req)
			if err != nil {
				if err != io.EOF {
					h.logger.Warnf("http.ReadResponse %v <- %v: %v", params.ClientAddr, params.ServerAddr, err)
				}
				return
			}
			defer func() { _ = resp.Body.Close() }()

			reqBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				h.logger.Warnf("ioutil.ReadAll(req.Body) %v -> %v: %v", params.ClientAddr, params.ServerAddr, err)
				return
			}
			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				h.logger.Warnf("ioutil.ReadAll(resp.Body) %v <- %v: %v", params.ClientAddr, params.ServerAddr, err)
				return
			}
			for _, inspector := range h.inspectors {
				p := &Parameters{Req: req, Resp: resp, ReqBytes: reqBytes, RespBytes: respBytes, SessionParams: params}
				if err = inspector.Inspect(p); err != nil {
					h.logger.Warnf("inspector.Inspect %v <-> %v: %v", params.ClientAddr, params.ServerAddr, err)
					return
				}
			}
		}
	}()
	return cpw, spw, nil
}
