package mirror

import (
	"github.com/qencept/gomitm/internal/shuttle"
	"io"
	"os"
	"strconv"
	"time"
)

type Dumper struct {
	dir string
}

func (d *Dumper) Session(client, server shuttle.Stream, sni string) (io.WriteCloser, io.WriteCloser) {
	ts := strconv.Itoa(int(time.Now().Unix()))
	c2s, _ := os.Create(d.dir + "/" + ts + "[" + sni + "]" + client.RemoteAddr().String() + "->" + server.RemoteAddr().String())
	s2c, _ := os.Create(d.dir + "/" + ts + "[" + sni + "]" + client.RemoteAddr().String() + "<-" + server.RemoteAddr().String())
	return c2s, s2c
}

func NewDumper() SessionInspector {
	return &Dumper{"dump"}
}
