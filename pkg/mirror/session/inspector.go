package session

import (
	"io"
)

type Inspector interface {
	InitWriteClosers(params *Parameters) (c2s, s2c io.WriteCloser, err error)
}
