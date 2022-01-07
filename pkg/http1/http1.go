package http1

import (
	"bufio"
	"github.com/qencept/gomitm/pkg/backup"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/session"
	"io"
	"net/http"
	"net/http/httputil"
)

type Http1 struct {
	logger   logger.Logger
	copy     *session.Copy
	mutators []Mutator
}

func New(logger logger.Logger, mutators ...Mutator) *Http1 {
	return &Http1{logger: logger, mutators: mutators, copy: session.NewCopy(logger)}
}

func (h *Http1) MutateForward(w io.Writer, r io.Reader, sp session.Parameters) {
	br := backup.NewReader(r)
	for {
		req, err := http.ReadRequest(bufio.NewReader(br))
		if err == io.EOF {
			break
		} else if err != nil {
			br.Reset()
			h.copy.MutateForward(w, br, sp)
			return
		}
		defer func(req *http.Request) {
			_ = req.Body.Close()
		}(req)
		// this is HTTP/2?
		if req.Method == "PRI" {
			br.Reset()
			h.copy.MutateForward(w, br, sp)
			return
		}
		for _, mutator := range h.mutators {
			req = mutator.MutateRequest(req, sp)
		}
		request, err := httputil.DumpRequest(req, true)
		if err != nil {
			h.logger.Warnln("Http1 req serialization: ", err)
			return
		}
		if _, err = w.Write(request); err != nil {
			h.logger.Warnln("Http1 req writing: ", err)
			return
		}
	}
}

func (h *Http1) MutateBackward(w io.Writer, r io.Reader, sp session.Parameters) {
	br := backup.NewReader(r)
	for {
		resp, err := http.ReadResponse(bufio.NewReader(br), nil)
		if err == io.EOF {
			break
		} else if err != nil {
			if err == io.ErrUnexpectedEOF {
				h.logger.Warnln("MutateBackward workarounds io.ErrUnexpectedEOF")
			} else {
				br.Reset()
				h.copy.MutateBackward(w, br, sp)
			}
			return
		}
		defer func(resp *http.Response) {
			_ = resp.Body.Close()
		}(resp)
		for _, mutator := range h.mutators {
			resp = mutator.MutateResponse(resp, sp)
		}
		response, err := httputil.DumpResponse(resp, true)
		if err != nil {
			h.logger.Warnln("Http1 resp serialization: ", err)
			return
		}
		if _, err = w.Write(response); err != nil {
			h.logger.Warnln("Http1 resp writing: ", err)
			return
		}
	}
}

var _ session.Mutator = (*Http1)(nil)
