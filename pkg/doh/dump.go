package doh

import (
	"fmt"
	"github.com/qencept/gomitm/pkg/logger"
	"golang.org/x/net/dns/dnsmessage"
	"net"
	"os"
)

type dump struct {
	logger logger.Logger
}

func NewDump(logger logger.Logger) *dump {
	return &dump{logger: logger}
}

func (d *dump) MutateQuestion(questions []dnsmessage.Question) []dnsmessage.Question {
	for _, q := range questions {
		_, _ = fmt.Fprintln(os.Stdout, q.Name, q.Type)
	}
	return questions
}

func (d *dump) MutateAnswer(answers []dnsmessage.Resource) []dnsmessage.Resource {
	for _, a := range answers {
		var str string
		switch b := a.Body.(type) {
		case *dnsmessage.AResource:
			str = net.IPv4(b.A[0], b.A[1], b.A[2], b.A[3]).String()
		case *dnsmessage.CNAMEResource:
			str = b.CNAME.String()
		default:
			str = b.GoString()
		}
		_, _ = fmt.Fprintln(os.Stdout, a.Header.Name, a.Header.Type, str)
	}
	return answers
}

var _ Mutator = (*dump)(nil)
