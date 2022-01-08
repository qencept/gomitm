package session

import (
	"github.com/qencept/gomitm/pkg/storage"
	"io"
)

type WriteCloseWriter interface {
	io.Writer
	CloseWrite() error
}

type Mutator interface {
	MutateForward(w io.Writer, r io.Reader, sp storage.Parameters)
	MutateBackward(w io.Writer, r io.Reader, sp storage.Parameters)
}
