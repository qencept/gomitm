package http1

import (
	"github.com/qencept/gomitm/pkg/logger"
	"net/http"
)

type defaultModifier struct {
	logger logger.Logger
}

func NewDefault(l logger.Logger) *defaultModifier {
	return &defaultModifier{logger: l}
}

func (d *defaultModifier) ModifyRequest(req *http.Request) bool {
	return true
}

func (d *defaultModifier) ModifyResponse(resp *http.Response) bool {
	return true
}
