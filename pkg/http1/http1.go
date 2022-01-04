package http1

import (
	"bufio"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/session"
	"io"
	"net/http"
	"net/http/httputil"
)

type Http1 struct {
	logger    logger.Logger
	modifiers []Modifier
}

func New(l logger.Logger, modifiers ...Modifier) *Http1 {
	return &Http1{
		logger:    l,
		modifiers: append(modifiers, NewDefault(l)),
	}
}

func (h *Http1) Modify(cr, sr io.Reader, cw, sw session.WriteCloseWriter) bool {
	req, err := http.ReadRequest(bufio.NewReader(cr))
	if err != nil {
		if err != io.EOF {
		}
	}
	defer func() { _ = req.Body.Close() }()

	for _, modifier := range h.modifiers {
		if modifier.ModifyRequest(req) {
			break
		}
	}

	request, err := httputil.DumpRequest(req, true)
	if err != nil {
	}

	_, err = sw.Write(request)
	if err != nil {
	}

	resp, err := http.ReadResponse(bufio.NewReader(sr), req)
	if err != nil {
		if err != io.EOF {
		}
	}
	defer func() { _ = resp.Body.Close() }()

	for _, modifier := range h.modifiers {
		if modifier.ModifyResponse(resp) {
			break
		}
	}

	response, err := httputil.DumpResponse(resp, true)
	if err != nil {
	}

	_, err = cw.Write(response)
	if err != nil {
	}

	return true
}
