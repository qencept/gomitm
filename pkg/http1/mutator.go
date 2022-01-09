package http1

import (
	"github.com/qencept/gomitm/pkg/session"
	"net/http"
)

type Creator interface {
	Create() Mutator
}

type Mutator interface {
	MutateRequest(req *http.Request, sp session.Parameters) *http.Request
	MutateResponse(resp *http.Response, sp session.Parameters) *http.Response
}
