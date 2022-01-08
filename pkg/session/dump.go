package session

import (
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/storage"
	"io"
)

type dump struct {
	logger logger.Logger
	path   string
}

func NewDump(logger logger.Logger, path string) *dump {
	return &dump{logger: logger, path: path}
}

func (d *dump) MutateForward(w io.Writer, r io.Reader, sp storage.Parameters) {
	f, err := storage.New(storage.Forward, d.path, sp)
	if err != nil {
		d.logger.Errorln("session new dump: ", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err = io.Copy(w, io.TeeReader(r, f)); err != nil {
		d.logger.Warnln("Copy.MutateForward: ", err)
	}
}

func (d *dump) MutateBackward(w io.Writer, r io.Reader, sp storage.Parameters) {
	f, err := storage.New(storage.Backward, d.path, sp)
	if err != nil {
		d.logger.Errorln("Http1 new dump: ", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err = io.Copy(w, io.TeeReader(r, f)); err != nil {
		d.logger.Warnln("Copy.MutateBackward: ", err)
	}
}

var _ Mutator = (*dump)(nil)
