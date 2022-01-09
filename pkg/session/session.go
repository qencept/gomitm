package session

import (
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/shuttler"
	"io"
	"sync"
)

type Session struct {
	logger   logger.Logger
	creators []Creator
}

func New(logger logger.Logger, creators ...Creator) shuttler.Shuttler {
	return &Session{logger: logger, creators: creators}
}

func (s *Session) Shuttle(client, server shuttler.Connection) {
	sp := NewParameters(client, server)
	var fwg, bwg sync.WaitGroup
	var fcr io.Reader = client
	var bcw io.WriteCloser = client
	for i, creator := range s.creators {
		var fnr, bcr io.Reader
		var fcw, bnw io.WriteCloser
		if i == len(s.creators)-1 {
			fcw, bcr = server, server
		} else {
			fnr, fcw = io.Pipe()
			bcr, bnw = io.Pipe()
		}
		mutator := creator.Create()
		startShuttle(Forward, mutator, fcw, fcr, sp, fcw == server, &fwg)
		startShuttle(Backward, mutator, bcw, bcr, sp, bcw == client, &bwg)
		fcr, bcw = fnr, bnw
	}
	fwg.Wait()
	bwg.Wait()
}

func startShuttle(dir Direction, m Mutator, w io.WriteCloser, r io.Reader, sp *Parameters, ep bool, wg *sync.WaitGroup) {
	wg.Add(1)
	go func(m Mutator, w io.WriteCloser, r io.Reader, sp *Parameters) {
		defer wg.Done()
		defer func() {
			if ep {
				_ = w.(interface{ CloseWrite() error }).CloseWrite()
			} else {
				_ = w.Close()
			}
		}()
		switch dir {
		case Forward:
			m.MutateForward(w, r, *sp)
		case Backward:
			m.MutateBackward(w, r, *sp)
		}
	}(m, w, r, sp)
}

var _ shuttler.Shuttler = (*Session)(nil)
