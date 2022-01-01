package mirror

import (
	"github.com/qencept/gomitm/internal/shuttle"
	"io"
	"net/http"
)

type SessionInspector interface {
	Session(client, server shuttle.Stream, sni string) (io.WriteCloser, io.WriteCloser)
}

type HttpInspector interface {
	Http(req *http.Request, resp *http.Response)
}
