package session

import (
	"github.com/qencept/gomitm/pkg/logger"
	"io"
)

type defaultMutator struct {
	logger logger.Logger
}

func NewDefault(logger logger.Logger) *defaultMutator {
	return &defaultMutator{logger: logger}
}

func (d *defaultMutator) MutateForward(w io.Writer, r io.Reader, _ *Parameters) {
	if _, err := io.Copy(w, r); err != nil {
		d.logger.Warnln("defaultMutator.MutateForward: ", err)
	}
}

func (d *defaultMutator) MutateBackward(w io.Writer, r io.Reader, _ *Parameters) {
	if _, err := io.Copy(w, r); err != nil {
		d.logger.Warnln("defaultMutator.MutateBackward: ", err)
	}
}

var _ Mutator = (*defaultMutator)(nil)
