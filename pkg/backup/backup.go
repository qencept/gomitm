package backup

import (
	"bytes"
	"io"
)

type Backup struct {
	r io.Reader
	b []*bytes.Buffer
}

func NewReader(r io.Reader) *Backup {
	return &Backup{r: r, b: []*bytes.Buffer{{}, {}}}
}

func (b *Backup) Read(p []byte) (n int, err error) {
	return io.TeeReader(io.MultiReader(b.b[0], b.r), b.b[1]).Read(p)
}

func (b *Backup) Reset() {
	_, _ = b.b[0].WriteTo(b.b[1])
	b.b = b.b[1:]
	b.b = append(b.b, &bytes.Buffer{})
}
