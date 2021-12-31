package mirror

import (
	"github.com/qencept/gomitm/internal/shuttle"
	"io"
)

type Mirror struct {
	inspectors []Inspector
}

func mirror(dst, src shuttle.Stream, mirror io.Writer) {
	io.Copy(dst, io.TeeReader(src, mirror))
	dst.CloseWrite()
}

func (m *Mirror) Shuttle(client, server shuttle.Stream, sni string) error {
	c2s, s2c := make([]io.Writer, 0), make([]io.Writer, 0)
	for _, inspector := range m.inspectors {
		c, s := inspector.Session(client, server, sni)
		defer c.Close()
		defer s.Close()
		c2s, s2c = append(c2s, c), append(s2c, s)
	}

	done := make(chan struct{})
	go func() {
		mirror(server, client, io.MultiWriter(c2s...))
		done <- struct{}{}
	}()
	mirror(client, server, io.MultiWriter(s2c...))
	<-done

	return nil
}

func New(inspectors ...Inspector) shuttle.Shuttle {
	return &Mirror{inspectors: inspectors}
}
