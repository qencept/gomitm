package session

import "io"

type WriteCloseWriter interface {
	io.Writer
	CloseWrite() error
}

type Mutator interface {
	MutateForward(w io.Writer, r io.Reader)
	MutateBackward(w io.Writer, r io.Reader)
}
