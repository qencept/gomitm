package dump

import (
	"github.com/qencept/gomitm/pkg/http1"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/session"
	"github.com/qencept/gomitm/pkg/storage"
	"net/http"
	"net/http/httputil"
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
	defer func() {
		_ = f.Close()
	}()
	request, err := httputil.DumpRequest(req, true)
	if err != nil {
		d.logger.Warnln("Http1 dump req serialization: ", err)
		return req
	}
	if _, err = f.Write(request); err != nil {
		d.logger.Errorln("Http1 dump req writing: ", err)
		return req
	}
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
	response, err := httputil.DumpResponse(resp, true)
	if err != nil {
		d.logger.Warnln("Http1 dump resp serialization: ", err)
		return resp
	}
	if _, err = f.Write(response); err != nil {
		d.logger.Errorln("Http1 dump resp writing: ", err)
		return resp
	}
	return resp
}

var _ http1.Mutator = (*dump)(nil)
