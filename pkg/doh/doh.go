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

type creator struct {
	logger   logger.Logger
	creators []Creator
}

func New(logger logger.Logger, creators ...Creator) http1.Creator {
	return &creator{logger: logger, creators: creators}
}

func (c *creator) Create() http1.Mutator {
	var mutators []Mutator
	for _, http1Creator := range c.creators {
		mutators = append(mutators, http1Creator.Create())
	}
	return &doh{
		logger:   c.logger,
		mutators: mutators,
	}
}

type doh struct {
	logger   logger.Logger
	mutators []Mutator
}

func (d *doh) MutateRequest(req *http.Request, sp session.Parameters) *http.Request {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		d.logger.Warnln("doh body reading: ", err)
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
		d.logger.Warnln("doh req pack: ", err)
		return req
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(pack))
	req.ContentLength = int64(len(pack))
	return req
}

func (d *doh) MutateResponse(resp *http.Response, sp session.Parameters) *http.Response {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		d.logger.Warnln("doh body reading: ", err)
		return resp
	}
	msg := dnsmessage.Message{}
	if err = msg.Unpack(body); err != nil {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		return resp
	}
	for i := len(d.mutators) - 1; i >= 0; i-- {
		msg.Answers = d.mutators[i].MutateAnswer(msg.Answers, sp)
	}
	pack, err := msg.Pack()
	if err != nil {
		d.logger.Warnln("doh resp pack: ", err)
		return resp
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(pack))
	resp.ContentLength = int64(len(pack))
	return resp
}

var _ http1.Creator = (*creator)(nil)
var _ http1.Mutator = (*doh)(nil)
