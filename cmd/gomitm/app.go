package main

import (
	"github.com/qencept/gomitm/pkg/storage"
	"golang.org/x/net/dns/dnsmessage"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (d *App) MutateQuestion(questions []dnsmessage.Question, _ storage.Parameters) []dnsmessage.Question {
	return questions
}

func (d *App) MutateAnswer(answers []dnsmessage.Resource, _ storage.Parameters) []dnsmessage.Resource {
	for _, a := range answers {
		if a.Header.Name.String() == "www.example.com." {
			if ar, ok := a.Body.(*dnsmessage.AResource); ok {
				ar.A = [4]byte{1, 1, 1, 1}
			}
		}
	}
	return answers
}
