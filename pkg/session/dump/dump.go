package dump

import (
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/session"
	"github.com/qencept/gomitm/pkg/storage"
	"io"
)

type creator struct {
	logger logger.Logger
	path   string
}

func New(logger logger.Logger, path string) session.Creator {
	return &creator{logger: logger, path: path}
}

func (c *creator) Create() session.Mutator {
	return &dump{logger: c.logger, path: c.path}
}

type dump struct {
	logger logger.Logger
	path   string
}

func (d *dump) MutateForward(w io.Writer, r io.Reader, sp session.Parameters) {
	f, err := storage.New(session.Forward, d.path, sp)
	if err != nil {
		d.logger.Errorln("session.Forward new dump: ", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err = io.Copy(w, io.TeeReader(r, f)); err != nil {
		d.logger.Warnln("session.Forward dump copy: ", err)
	}
}

func (d *dump) MutateBackward(w io.Writer, r io.Reader, sp session.Parameters) {
	f, err := storage.New(session.Backward, d.path, sp)
	if err != nil {
		d.logger.Errorln("session.Backward new dump: ", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err = io.Copy(w, io.TeeReader(r, f)); err != nil {
		d.logger.Warnln("session.Backward dump copy: ", err)
	}
}

var _ session.Mutator = (*dump)(nil)
