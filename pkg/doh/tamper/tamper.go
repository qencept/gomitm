package tamper

import (
	"github.com/qencept/gomitm/pkg/doh"
	"github.com/qencept/gomitm/pkg/session"
	"golang.org/x/net/dns/dnsmessage"
)

type SubstitutionTypeA map[string][4]byte

type creator struct {
	typeA SubstitutionTypeA
}

func New(typeA SubstitutionTypeA) doh.Creator {
	return &creator{typeA: typeA}
}

func (c *creator) Create() doh.Mutator {
	return &tamper{typeA: c.typeA}
}

type tamper struct {
	typeA SubstitutionTypeA
}

func (t *tamper) MutateQuestion(questions []dnsmessage.Question, _ session.Parameters) []dnsmessage.Question {
	return questions
}

func (t *tamper) MutateAnswer(answers []dnsmessage.Resource, _ session.Parameters) []dnsmessage.Resource {
	for _, a := range answers {
		if ip, ok := t.typeA[a.Header.Name.String()]; ok {
			if ar, ok := a.Body.(*dnsmessage.AResource); ok {
				ar.A = ip
			}
		}
	}
	return answers
}
