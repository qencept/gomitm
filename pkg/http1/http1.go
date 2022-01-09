package http1

import (
	"bufio"
	"github.com/qencept/gomitm/pkg/backup"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/session"
	"github.com/qencept/gomitm/pkg/session/copier"
	"io"
	"net/http"
	"net/http/httputil"
)

type creator struct {
	logger   logger.Logger
	creators []Creator
}

func New(logger logger.Logger, creators ...Creator) session.Creator {
	return &creator{logger: logger, creators: creators}
}

func (c *creator) Create() session.Mutator {
	var mutators []Mutator
	for _, http1Creator := range c.creators {
		mutators = append(mutators, http1Creator.Create())
	}
	return &http1{
		logger:       c.logger,
		mutators:     mutators,
		copier:       copier.New(c.logger).Create(),
		skipResponse: make(chan bool, 1),
	}
}

type http1 struct {
	logger       logger.Logger
	copier       session.Mutator
	mutators     []Mutator
	skipResponse chan bool
}

func (h *http1) justCopy(w io.Writer, br *backup.Backup, sp session.Parameters) {
	br.Reset()
	h.copier.MutateForward(w, br, sp)
}

func (h *http1) MutateForward(w io.Writer, r io.Reader, sp session.Parameters) {
	br := backup.NewReader(r)
	for {
		req, err := http.ReadRequest(bufio.NewReader(br))
		if err == io.EOF {
			break
		} else if err != nil {
			h.skipResponse <- true
			h.logger.Warnln("http1 req read: ", err)
			h.justCopy(w, br, sp)
			return
		}
		defer func(req *http.Request) {
			_ = req.Body.Close()
		}(req)
		if req.Proto == "HTTP/2.0" {
			h.skipResponse <- true
			h.logger.Debugln("http1 req read: ", "HTTP/2.0")
			h.justCopy(w, br, sp)
			return
		}
		h.skipResponse <- false
		for _, mutator := range h.mutators {
			req = mutator.MutateRequest(req, sp)
		}
		request, err := httputil.DumpRequest(req, true)
		if err != nil {
			h.logger.Warnln("http1 req serialization: ", err)
			return
		}
		if _, err = w.Write(request); err != nil {
			h.logger.Warnln("http1 req writing: ", err)
			return
		}
	}
}

func (h *http1) MutateBackward(w io.Writer, r io.Reader, sp session.Parameters) {
	br := backup.NewReader(r)
	for {
		if <-h.skipResponse {
			h.justCopy(w, br, sp)
			return
		}
		resp, err := http.ReadResponse(bufio.NewReader(br), nil)
		if err == io.EOF {
			break
		} else if err != nil {
			h.logger.Warnln("http1 resp read: ", err)
			h.justCopy(w, br, sp)
			return
		}
		defer func(resp *http.Response) {
			_ = resp.Body.Close()
		}(resp)
		for i := len(h.mutators) - 1; i >= 0; i-- {
			resp = h.mutators[i].MutateResponse(resp, sp)
		}
		response, err := httputil.DumpResponse(resp, true)
		if err != nil {
			h.logger.Warnln("http1 resp serialization: ", err)
			return
		}
		if _, err = w.Write(response); err != nil {
			h.logger.Warnln("http1 resp writing: ", err)
			return
		}
	}
}

var _ session.Creator = (*creator)(nil)
var _ session.Mutator = (*http1)(nil)
