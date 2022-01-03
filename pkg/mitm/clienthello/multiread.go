package clienthello

import (
	"github.com/qencept/gomitm/pkg/shuttle"
	"io"
)

type ConnectionMultiRead struct {
	shuttle.Connection
	multiReader io.Reader
}

func (c *ConnectionMultiRead) Read(p []byte) (int, error) { return c.multiReader.Read(p) }

var _ shuttle.Connection = (*ConnectionMultiRead)(nil)
