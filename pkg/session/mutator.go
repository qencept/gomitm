package session

import (
	"io"
)

type Creator interface {
	Create() Mutator
}

type Mutator interface {
	MutateForward(w io.Writer, r io.Reader, sp Parameters)
	MutateBackward(w io.Writer, r io.Reader, sp Parameters)
}
