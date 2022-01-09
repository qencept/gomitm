package storage

import (
	"fmt"
	"github.com/qencept/gomitm/pkg/session"
	"io"
	"os"
)

func fileName(dir session.Direction, path, ts string, sp session.Parameters) string {
	template := path + "/" + ts + "#" + sp.Client + "%s" + sp.Server + "#" + sp.Sni
	switch dir {
	case session.Forward:
		return fmt.Sprintf(template, "->")
	case session.Backward:
		return fmt.Sprintf(template, "<-")
	default:
		return ""
	}
}

func New(dir session.Direction, path string, sp session.Parameters) (io.WriteCloser, error) {
	name := fileName(dir, path, sp.Timestamp, sp)
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}
