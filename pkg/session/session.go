package session

import (
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/shuttler"
	"github.com/qencept/gomitm/pkg/storage"
	"io"
	"sync"
)

type Session struct {
	logger   logger.Logger
	mutators []Mutator
}

func New(logger logger.Logger, mutators ...Mutator) shuttler.Shuttler {
	if len(mutators) == 0 {
		mutators = append(mutators, NewCopy(logger))
	}
	return &Session{logger: logger, mutators: mutators}
}

func (s *Session) Shuttle(client, server shuttler.Connection) {
	sp := storage.NewParameters(client, server)
	var fwg, bwg sync.WaitGroup
	var fcr, fnr, bcr io.Reader
	var bcw, fcw, bnw io.WriteCloser
	fcr, bcw = client, client
	for i, mutator := range s.mutators {
		if i == len(s.mutators)-1 {
			fcw, bcr = server, server
		} else {
			fnr, fcw = io.Pipe()
			bcr, bnw = io.Pipe()
		}
		fwg.Add(1)
		bwg.Add(1)
		go func(m Mutator, w io.WriteCloser, r io.Reader, sp storage.Parameters) {
			defer fwg.Done()
			defer func() {
				if w == server {
					_ = server.CloseWrite()
				} else {
					_ = w.Close()
				}
			}()
			m.MutateForward(w, r, sp)
		}(mutator, fcw, fcr, *sp)
		go func(m Mutator, w io.WriteCloser, r io.Reader, sp storage.Parameters) {
			defer bwg.Done()
			defer func() {
				if w == client {
					_ = client.CloseWrite()
				} else {
					_ = w.Close()
				}
			}()
			m.MutateBackward(w, r, sp)
		}(mutator, bcw, bcr, *sp)
		fcr, bcw = fnr, bnw
	}
	fwg.Wait()
	bwg.Wait()
}

var _ shuttler.Shuttler = (*Session)(nil)
