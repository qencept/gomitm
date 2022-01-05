package doh

import (
	"github.com/qencept/gomitm/pkg/session"
	"golang.org/x/net/dns/dnsmessage"
)

type Mutator interface {
	MutateQuestion(questions []dnsmessage.Question, sp *session.Parameters) []dnsmessage.Question
	MutateAnswer(answers []dnsmessage.Resource, sp *session.Parameters) []dnsmessage.Resource
}
