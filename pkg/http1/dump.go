package http1

import (
	"github.com/qencept/gomitm/pkg/logger"
	"net/http"
)

type dump struct {
	logger logger.Logger
}

func NewDump(logger logger.Logger) *dump {
	return &dump{logger: logger}
}

func (d *dump) MutateRequest(req *http.Request) *http.Request {
	return req
}

func (d *dump) MutateResponse(resp *http.Response) *http.Response {
	return resp
}

var _ Mutator = (*dump)(nil)
