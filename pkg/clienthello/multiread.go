package clienthello

import (
	"github.com/qencept/gomitm/pkg/shuttler"
	"io"
)

type connectionMultiRead struct {
	shuttler.Connection
	multiReader io.Reader
}

func (c *connectionMultiRead) Read(p []byte) (int, error) { return c.multiReader.Read(p) }

var _ shuttler.Connection = (*connectionMultiRead)(nil)
