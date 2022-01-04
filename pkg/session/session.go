package session

import (
	"github.com/qencept/gomitm/pkg/backup"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/shuttler"
)

type Session struct {
	logger    logger.Logger
	modifiers []Modifier
}

func New(l logger.Logger, modifiers ...Modifier) shuttler.Shuttler {
	return &Session{
		logger:    l,
		modifiers: append(modifiers, NewDefault(l)),
	}
}

func (s *Session) Shuttle(client, server shuttler.Connection) {
	clientReader := backup.NewReader(client)
	serverReader := backup.NewReader(server)
	for _, modifier := range s.modifiers {
		if modifier.Modify(clientReader, serverReader, client, server) {
			break
		}
		clientReader.Reset()
		serverReader.Reset()
	}
}
