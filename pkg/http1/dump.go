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
		d.logger.Errorln("Http1 new dump: ", err)
		return req
	}
	defer func() { _ = f.Close() }()
	if _, err = fmt.Fprintf(f, "%s %s %s%s\n", req.Proto, req.Method, req.Host, req.RequestURI); err != nil {
		d.logger.Errorln("Http1 dumping: ", err)
		return req
	}
	for k, v := range req.Header {
		if _, err = fmt.Fprintf(f, "%s:%s\n", k, v); err != nil {
			d.logger.Errorln("Http1 dumping: ", err)
			return req
		}
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		d.logger.Errorln("Http1 dump body restoring: ", err)
		return req
	}
	if _, err = f.Write(body); err != nil {
		d.logger.Errorln("Http1 dump body restoring(2): ", err)
		return req
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return req
}

func (d *dump) MutateResponse(resp *http.Response, sp session.Parameters) *http.Response {
	f, err := storage.New(session.Backward, d.path, sp)
	if err != nil {
		d.logger.Errorln("Http1 new dump: ", err)
		return resp
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err = fmt.Fprintf(f, "%s %s\n", resp.Proto, resp.Status); err != nil {
		d.logger.Errorln("Http1 dumping: ", err)
		return resp
	}
	for k, v := range resp.Header {
		if _, err = fmt.Fprintf(f, "%s:%s\n", k, v); err != nil {
			d.logger.Errorln("Http1 dumping: ", err)
			return resp
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		d.logger.Errorln("Http1 dump body restoring: ", err)
		return resp
	}
	if _, err = f.Write(body); err != nil {
		d.logger.Errorln("Http1 dump body restoring(2): ", err)
		return resp
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return resp
}

var _ Mutator = (*dump)(nil)
