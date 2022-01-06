package clienthello

import (
	"github.com/qencept/gomitm/pkg/shuttler"
	"io"
)

type connection struct {
	shuttler.Connection
	reader io.Reader
}

func (c *connection) Read(p []byte) (int, error) { return c.reader.Read(p) }

var _ shuttler.Connection = (*connection)(nil)
