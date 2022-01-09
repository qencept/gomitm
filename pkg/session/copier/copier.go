package copier

import (
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/session"
	"io"
)

type creator struct {
	logger logger.Logger
}

func New(logger logger.Logger) session.Creator {
	return &creator{logger: logger}
}

func (c *creator) Create() session.Mutator {
	return &copier{logger: c.logger}
}

type copier struct {
	logger logger.Logger
}

func (d *copier) MutateForward(w io.Writer, r io.Reader, _ session.Parameters) {
	if _, err := io.Copy(w, r); err != nil {
		d.logger.Warnln("copier.MutateForward: ", err)
	}
}

func (d *copier) MutateBackward(w io.Writer, r io.Reader, _ session.Parameters) {
	if _, err := io.Copy(w, r); err != nil {
		d.logger.Warnln("copier.MutateBackward: ", err)
	}
}

var _ session.Creator = (*creator)(nil)
var _ session.Mutator = (*copier)(nil)
