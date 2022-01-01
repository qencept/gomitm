package mirror

import (
	"bytes"
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"net/http"
)

type Doh struct {
}

func NewDoh() HttpInspector {
	return &Doh{}
}

func (d *Doh) Http(req *http.Request, resp *http.Response) {
	buf, msg := bytes.Buffer{}, dnsmessage.Message{}
	buf.ReadFrom(resp.Body)
	err := msg.Unpack(buf.Bytes())
	if err == nil {
		for _, a := range msg.Answers {
			fmt.Println(a.Header.Name, a.Header.Type, a.Body)
		}
	}
}
