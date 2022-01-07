package doh

import (
	"fmt"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/session"
	"github.com/qencept/gomitm/pkg/storage"
	"golang.org/x/net/dns/dnsmessage"
	"net"
)

type dump struct {
	logger logger.Logger
	path   string
}

func NewDump(logger logger.Logger, path string) *dump {
	return &dump{logger: logger, path: path}
}

func (d *dump) MutateQuestion(questions []dnsmessage.Question, sp session.Parameters) []dnsmessage.Question {
	f, err := storage.New(session.Forward, d.path, sp)
	if err != nil {
		d.logger.Errorln("Doh new dump: ", err)
		return questions
	}
	defer func() {
		_ = f.Close()
	}()
	for _, q := range questions {
		if _, err = fmt.Fprintln(f, q.Name, q.Type); err != nil {
			d.logger.Errorln("Doh dumping: ", err)
			return questions
		}
	}
	return questions
}

func (d *dump) MutateAnswer(answers []dnsmessage.Resource, sp session.Parameters) []dnsmessage.Resource {
	f, err := storage.New(session.Backward, d.path, sp)
	if err != nil {
		d.logger.Errorln("Doh new dump: ", err)
		return answers
	}
	defer func() { _ = f.Close() }()
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
		if _, err = fmt.Fprintln(f, a.Header.Name, a.Header.Type, str); err != nil {
			d.logger.Errorln("Doh dumping: ", err)
			return answers
		}
	}
	return answers
}

var _ Mutator = (*dump)(nil)
