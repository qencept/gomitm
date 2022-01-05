package storage

import (
	"fmt"
	"github.com/qencept/gomitm/pkg/session"
	"io"
	"os"
	"strconv"
	"time"
)

func fileName(direction int, path, ts string, sp *session.Parameters) string {
	template := path + "/" + ts + "#" + sp.Client.String() + "%s" + sp.Server.String() + "#" + sp.Sni
	switch direction {
	case session.Forward:
		return fmt.Sprintf(template, "->")
	case session.Backward:
		return fmt.Sprintf(template, "<-")
	default:
		return ""
	}
}

func New(direction int, path string, sp *session.Parameters) (io.WriteCloser, error) {
	ts := strconv.Itoa(int(time.Now().Unix()))
	name := fileName(direction, path, ts, sp)
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}
