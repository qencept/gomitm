package session

import (
	"github.com/qencept/gomitm/pkg/mitm/shuttle"
	"io"
)

type Inspector interface {
	Inspect(session *shuttle.Session) (io.WriteCloser, io.WriteCloser)
}
