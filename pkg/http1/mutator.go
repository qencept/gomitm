package http1

import (
	"github.com/qencept/gomitm/pkg/storage"
	"net/http"
)

type Mutator interface {
	MutateRequest(req *http.Request, sp storage.Parameters) *http.Request
	MutateResponse(resp *http.Response, sp storage.Parameters) *http.Response
}
