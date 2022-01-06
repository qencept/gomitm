package session

import (
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/shuttler"
	"io"
	"sync"
)

type Session struct {
	logger   logger.Logger
	mutators []Mutator
}

func New(logger logger.Logger, mutators ...Mutator) shuttler.Shuttler {
	return &Session{logger: logger, mutators: mutators}
}

func (s *Session) Shuttle(client, server shuttler.Connection) {
	sp := NewParameters(client, server)
	var wg sync.WaitGroup
	var fcr, fnr, bcr io.Reader
	var bcw, fcw, bnw io.Writer
	fcr, bcw = client, client
	for i, mutator := range s.mutators {
		if i == len(s.mutators)-1 {
			fcw, bcr = server, server
		} else {
			fnr, fcw = io.Pipe()
			bcr, bnw = io.Pipe()
		}
		wg.Add(2)
		go func(m Mutator, w io.Writer, r io.Reader, sp Parameters) {
			defer wg.Done()
			m.MutateForward(w, r, sp)
		}(mutator, fcw, fcr, *sp)
		go func(m Mutator, w io.Writer, r io.Reader, sp Parameters) {
			defer wg.Done()
			m.MutateBackward(w, r, sp)
		}(mutator, bcw, bcr, *sp)
		fcr, bcw = fnr, bnw
	}
	wg.Wait()
	if err := server.CloseWrite(); err != nil {
		s.logger.Warnln("server.CloseWrite: ", err)
	}
	if err := client.CloseWrite(); err != nil {
		s.logger.Warnln("client.CloseWrite: ", err)
	}
}

var _ shuttler.Shuttler = (*Session)(nil)
