package session

import (
	"github.com/qencept/gomitm/pkg/logger"
	"io"
)

type defaultModifier struct {
	logger logger.Logger
}

func NewDefault(l logger.Logger) *defaultModifier {
	return &defaultModifier{logger: l}
}

func (d *defaultModifier) Modify(cr, sr io.Reader, cw, sw WriteCloseWriter) bool {
	done := make(chan struct{})
	go func() {
		defer func() { done <- struct{}{} }()
		if _, err := io.Copy(sw, cr); err != nil {
			d.logger.Warnln("defaultModifier: ", err)
		}
		if err := sw.CloseWrite(); err != nil {
			d.logger.Warnln("defaultModifier: ", err)
		}
	}()
	if _, err := io.Copy(cw, sr); err != nil {
		d.logger.Warnln("defaultModifier: ", err)
	}
	if err := cw.CloseWrite(); err != nil {
		d.logger.Warnln("defaultModifier: ", err)
	}
	<-done
	return true
}
