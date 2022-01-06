package http1

import (
	"bytes"
	"fmt"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/session"
	"github.com/qencept/gomitm/pkg/storage"
	"io/ioutil"
	"net/http"
)

type dump struct {
	logger logger.Logger
	path   string
}

func NewDump(logger logger.Logger, path string) *dump {
	return &dump{logger: logger, path: path}
}

func (d *dump) MutateRequest(req *http.Request, sp session.Parameters) *http.Request {
	f, err := storage.New(session.Forward, d.path, sp)
	if err != nil {
	}
	defer func() { _ = f.Close() }()
	if _, err = fmt.Fprintf(f, "%s %s %s%s\n\n", req.Proto, req.Method, req.Host, req.RequestURI); err != nil {
	}
	for k, v := range req.Header {
		if _, err = fmt.Fprintf(f, "%s:%s\n", k, v); err != nil {
		}
	}
	if _, err = fmt.Fprintf(f, "\n"); err != nil {
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
	}
	if _, err = f.Write(body); err != nil {
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return req
}

func (d *dump) MutateResponse(resp *http.Response, sp session.Parameters) *http.Response {
	f, err := storage.New(session.Backward, d.path, sp)
	if err != nil {
	}
	defer func() { _ = f.Close() }()
	if _, err = fmt.Fprintf(f, "%s %s\n\n", resp.Proto, resp.Status); err != nil {
	}
	for k, v := range resp.Header {
		if _, err = fmt.Fprintf(f, "%s:%s\n", k, v); err != nil {
		}
	}
	if _, err = fmt.Fprintf(f, "\n"); err != nil {
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	if _, err = f.Write(body); err != nil {
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return resp
}

var _ Mutator = (*dump)(nil)
