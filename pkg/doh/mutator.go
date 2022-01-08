package doh

import (
	"github.com/qencept/gomitm/pkg/storage"
	"golang.org/x/net/dns/dnsmessage"
)

type Mutator interface {
	MutateQuestion(questions []dnsmessage.Question, sp storage.Parameters) []dnsmessage.Question
	MutateAnswer(answers []dnsmessage.Resource, sp storage.Parameters) []dnsmessage.Resource
}
