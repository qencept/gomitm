package mirror

import (
	"github.com/qencept/gomitm/internal/shuttle"
	"io"
)

type Inspector interface {
	Session(client, server shuttle.Stream, sni string) (io.WriteCloser, io.WriteCloser)
}
