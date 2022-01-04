package session

import (
	"github.com/qencept/gomitm/pkg/shuttler"
	"io"
)

type Session struct {
}

func (s *Session) Shuttle(client, server shuttler.Connection) {
	done := make(chan struct{})
	go func() {
		io.Copy(server, client)
		server.CloseWrite()
		done <- struct{}{}
	}()
	io.Copy(client, server)
	client.CloseWrite()
	<-done
}

func New() shuttler.Shuttler {
	return &Session{}
}
