package http1

import "net/http"

type Mutator interface {
	MutateRequest(req *http.Request) *http.Request
	MutateResponse(resp *http.Response) *http.Response
}
