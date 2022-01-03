package http1

import (
	"github.com/qencept/gomitm/pkg/mirror/session"
	"net/http"
)

type Parameters struct {
	Req           *http.Request
	Resp          *http.Response
	ReqBytes      []byte
	RespBytes     []byte
	SessionParams *session.Parameters
}
