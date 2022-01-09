package doh

import (
	"bytes"
	"github.com/qencept/gomitm/pkg/http1"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/session"
	"golang.org/x/net/dns/dnsmessage"
	"io/ioutil"
	"net/http"
)

type Doh struct {
	logger   logger.Logger
	mutators []Mutator
}

func New(logger logger.Logger, mutators ...Mutator) *Doh {
	return &Doh{logger: logger, mutators: mutators}
}

func (d *Doh) MutateRequest(req *http.Request, sp session.Parameters) *http.Request {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		d.logger.Errorln("Doh body reading: ", err)
	}
	msg := dnsmessage.Message{}
	if err = msg.Unpack(body); err != nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		return req
	}
	for _, mutator := range d.mutators {
		msg.Questions = mutator.MutateQuestion(msg.Questions, sp)
	}
	pack, err := msg.Pack()
	if err != nil {
		d.logger.Errorln("Doh Request Pack: ", err)
		return req
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(pack))
	req.ContentLength = int64(len(pack))
	return req
}

func (d *Doh) MutateResponse(resp *http.Response, sp session.Parameters) *http.Response {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		d.logger.Errorln("Doh body reading: ", err)
		return resp
	}
	msg := dnsmessage.Message{}
	if err = msg.Unpack(body); err != nil {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		return resp
	}
	for _, mutator := range d.mutators {
		msg.Answers = mutator.MutateAnswer(msg.Answers, sp)
	}
	pack, err := msg.Pack()
	if err != nil {
		d.logger.Errorln("Doh Response Pack: ", err)
		return resp
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(pack))
	resp.ContentLength = int64(len(pack))
	return resp
}

var _ http1.Mutator = (*Doh)(nil)
