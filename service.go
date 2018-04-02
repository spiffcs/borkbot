// Package borkbot provides a small server that will send borks to the Samantha and Christopher Slack
package borkbot

import (
	"errors"
	"math/rand"
	"time"
)

// Service is the interface that provides the borkbot methods.
type Service interface {
	// FetchBork receives a request from the slack app and responds with a dog meme
	FetchBork(req fetchBorkRequest) (string, error)
}

type service struct {
	verficationToken string
}

func (s *service) FetchBork(req fetchBorkRequest) (string, error) {
	if s.verficationToken != req.Token {
		return "", errnotFromSlack
	}
	bork, err := s.borkGenerator()
	if err != nil {
		return "", err
	}
	return bork, nil
}

func (s *service) borkGenerator() (string, error) {
	rand.Seed(time.Now().Unix())
	borkURL := dogMemes[rand.Intn(len(dogMemes))]
	return borkURL, nil
}

func (s *service) borkFetcher() {

}

// NewService creates a borkbot service
func NewService(token string) Service {
	return &service{
		verficationToken: token,
	}
}

var dogMemes = []string{
	"https://i.imgur.com/utGgHkP_d.jpg?maxwidth=640&shape=thumb&fidelity=medium",
	"https://i.redd.it/z6idgq152gp01.png",
	"https://i.imgur.com/SzNHUWv.jpg",
	"https://i.imgur.com/VjCocZY.gifv",
	"https://i.redd.it/gru1mw246mn01.jpg",
	"https://i.imgur.com/eJTKvOa.jpg",
	"https://i.imgur.com/un5eaeD.gifv",
	"https://i.imgur.com/XQQh7v0.jpg",
	"https://i.redd.it/7tobptj97bk01.jpg",
	"https://i.redd.it/9cvfwf5vwik01.jpg",
	"https://i.redd.it/keszh72c3io01.jpg",
}

var errnotFromSlack = errors.New("why you try to be a slack?")
