package session

import (
	"github.com/qencept/gomitm/pkg/mitm/shuttle"
	"io"
	"net"
)

type Mirror struct {
	inspectors []Inspector
}

func mirror(dst, src net.Conn, mirror io.Writer) {
	io.Copy(dst, io.TeeReader(src, mirror))
	dst.(interface{ CloseWrite() error }).CloseWrite()
}

func (m *Mirror) Shuttle(session *shuttle.Session) {
	var c2s, s2c []io.Writer
	for _, inspector := range m.inspectors {
		c, s := inspector.Inspect(session)
		defer c.Close()
		defer s.Close()
		c2s, s2c = append(c2s, c), append(s2c, s)
	}

	done := make(chan struct{})
	go func() {
		mirror(session.Server, session.Client, io.MultiWriter(c2s...))
		done <- struct{}{}
	}()
	mirror(session.Client, session.Server, io.MultiWriter(s2c...))
	<-done
}

func New(inspectors ...Inspector) shuttle.Shuttle {
	return &Mirror{inspectors: inspectors}
}
