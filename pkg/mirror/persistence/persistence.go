package persistence

import (
	"fmt"
	"github.com/qencept/gomitm/pkg/mirror/session"
	"io"
	"os"
	"strconv"
	"time"
)

const (
	CliSer int = iota + 1
	SerCli
)

func fileName(direction int, path, ts, sni, client, server string) string {
	template := path + "/" + ts + "[" + sni + "]" + client + "%s" + server
	switch direction {
	case CliSer:
		return fmt.Sprintf(template, "->")
	case SerCli:
		return fmt.Sprintf(template, "<-")
	default:
		return ""
	}
}

func CreateFile(direction int, path string, params *session.Parameters) (io.WriteCloser, error) {
	ts := strconv.Itoa(int(time.Now().Unix()))
	name := fileName(direction, path, ts, params.Sni, params.ClientAddr.String(), params.ServerAddr.String())
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}
