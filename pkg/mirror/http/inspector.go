package httpi

import (
	"github.com/qencept/gomitm/pkg/mitm/shuttle"
	"net/http"
)

type Inspector interface {
	Inspect(req *http.Request, resp *http.Response, session *shuttle.Session)
}
