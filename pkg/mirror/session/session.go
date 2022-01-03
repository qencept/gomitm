package session

import (
	"github.com/qencept/gomitm/pkg/shuttle"
	"github.com/sirupsen/logrus"
	"io"
	"net"
)

type Session struct {
	logger     *logrus.Logger
	inspectors []Inspector
}

func shuttleAndMirror(dst, src net.Conn, mirror io.Writer) error {
	if _, err := io.Copy(dst, io.TeeReader(src, mirror)); err != nil {
		return err
	}
	if err := dst.(interface{ CloseWrite() error }).CloseWrite(); err != nil {
		return err
	}
	return nil
}

func (s *Session) Shuttle(client, server shuttle.Connection) {
	var cli2ser []io.Writer
	var ser2cli []io.Writer
	for _, inspector := range s.inspectors {
		cs, sc, err := inspector.InitWriteClosers(NewParameters(client, server))
		if err != nil {
			s.logger.Warnf("inspector.InitWriteClosers %v <-> %v: %v", client.RemoteAddr(), server.RemoteAddr(), err)
		} else {
			defer func() { _ = cs.Close() }()
			defer func() { _ = sc.Close() }()
			cli2ser = append(cli2ser, cs)
			ser2cli = append(ser2cli, sc)
		}
	}

	done := make(chan struct{})
	go func() {
		if err := shuttleAndMirror(server, client, io.MultiWriter(cli2ser...)); err != nil {
			s.logger.Warnf("shuttleAndMirror %v -> %v: %v", client.RemoteAddr(), server.RemoteAddr(), err)
		}
		done <- struct{}{}
	}()
	if err := shuttleAndMirror(client, server, io.MultiWriter(ser2cli...)); err != nil {
		s.logger.Warnf("shuttleAndMirror %v <- %v: %v", client.RemoteAddr(), server.RemoteAddr(), err)
	}
	<-done
}

func New(logger *logrus.Logger, inspectors ...Inspector) shuttle.Shuttle {
	return &Session{logger: logger, inspectors: inspectors}
}
