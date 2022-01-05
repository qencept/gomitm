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
	mutators []Mutator
}

func New(logger logger.Logger, mutators ...Mutator) *Http1 {
	return &Http1{logger: logger, mutators: mutators}
}

func (h *Http1) MutateForward(w io.Writer, r io.Reader, sp *session.Parameters) {
	br := backup.NewReader(r)
	req, err := http.ReadRequest(bufio.NewReader(br))
	if err != nil && err != io.EOF {
		h.logger.Infoln("Http1 parsing failed, fallback to session.storage: ", err)
		br.Reset()
		session.NewDefault(h.logger).MutateForward(w, br, sp)
	}
	defer func() { _ = req.Body.Close() }()

	for _, mutator := range h.mutators {
		mutator.MutateRequest(req, sp)
	}

	request, err := httputil.DumpRequest(req, true)
	if err != nil {
		h.logger.Errorln("Http1 Dump Request: ", err)
		return
	}
	if _, err = w.Write(request); err != nil {
		h.logger.Errorln("Http1 Request Write: ", err)
	}
}

func (h *Http1) MutateBackward(w io.Writer, r io.Reader, sp *session.Parameters) {
	br := backup.NewReader(r)
	resp, err := http.ReadResponse(bufio.NewReader(r), nil)
	if err != nil && err != io.EOF {
		h.logger.Infoln("Http1 parsing failed, fallback to session.storage: ", err)
		br.Reset()
		session.NewDefault(h.logger).MutateForward(w, br, sp)
	}
	defer func() { _ = resp.Body.Close() }()

	for _, mutator := range h.mutators {
		mutator.MutateResponse(resp, sp)
	}

	response, err := httputil.DumpResponse(resp, true)
	if err != nil {
		h.logger.Errorln("Http1 Dump Response: ", err)
		return
	}
	if _, err = w.Write(response); err != nil {
		h.logger.Errorln("Http1 Response Write: ", err)
	}
}

var _ session.Mutator = (*Http1)(nil)
