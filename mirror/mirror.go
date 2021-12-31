package mirror

import (
	"github.com/qencept/gomitm/internal/shuttle"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

type Mirror struct {
}

func (m *Mirror) Shuttle(client, server shuttle.Stream) error {
	ts := strconv.Itoa(int(time.Now().Unix()))
	c2s, _ := os.Create(ts + "#" + client.RemoteAddr().String() + "->" + server.RemoteAddr().String())
	defer c2s.Close()
	s2c, _ := os.Create(ts + "#" + client.RemoteAddr().String() + "<-" + server.RemoteAddr().String())
	defer s2c.Close()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		io.Copy(server, io.TeeReader(client, c2s))
		server.CloseWrite()
		wg.Done()
	}()
	go func() {
		io.Copy(client, io.TeeReader(server, s2c))
		client.CloseWrite()
		wg.Done()
	}()
	wg.Wait()

	return nil
}

func New() shuttle.Shuttle {
	return &Mirror{}
}
