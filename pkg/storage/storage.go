package storage

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func fileName(direction int, path, ts string, sp Parameters) string {
	template := path + "/" + ts + "#" + sp.Client + "%s" + sp.Server + "#" + sp.Sni
	switch direction {
	case Forward:
		return fmt.Sprintf(template, "->")
	case Backward:
		return fmt.Sprintf(template, "<-")
	default:
		return ""
	}
}

func New(direction int, path string, sp Parameters) (io.WriteCloser, error) {
	ts := strconv.Itoa(int(time.Now().Unix()))
	name := fileName(direction, path, ts, sp)
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}
