package session

import "io"

type WriteCloseWriter interface {
	io.Writer
	CloseWrite() error
}

type Modifier interface {
	Modify(cr, sr io.Reader, cw, sw WriteCloseWriter) bool
}
