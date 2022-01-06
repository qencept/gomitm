package session

import (
	"github.com/qencept/gomitm/pkg/logger"
	"io"
)

type Copy struct {
	logger logger.Logger
}

func NewCopy(logger logger.Logger) *Copy {
	return &Copy{logger: logger}
}

func (d *Copy) MutateForward(w io.Writer, r io.Reader, _ Parameters) {
	if _, err := io.Copy(w, r); err != nil {
		d.logger.Warnln("Copy.MutateForward: ", err)
	}
}

func (d *Copy) MutateBackward(w io.Writer, r io.Reader, _ Parameters) {
	if _, err := io.Copy(w, r); err != nil {
		d.logger.Warnln("Copy.MutateBackward: ", err)
	}
}

var _ Mutator = (*Copy)(nil)
