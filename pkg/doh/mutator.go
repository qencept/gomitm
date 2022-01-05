package doh

import "golang.org/x/net/dns/dnsmessage"

type Mutator interface {
	MutateQuestion(questions []dnsmessage.Question) []dnsmessage.Question
	MutateAnswer(answers []dnsmessage.Resource) []dnsmessage.Resource
}
