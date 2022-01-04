package http1

import "net/http"

type Modifier interface {
	ModifyRequest(req *http.Request) bool
	ModifyResponse(resp *http.Response) bool
}
