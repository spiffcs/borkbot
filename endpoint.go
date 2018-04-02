package borkbot

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints models a collection of all bork endpoints

type fetchBorkRequest struct {
	UserID      string `form:"user_id"`
	UserName    string `form:"user_name"`
	TeamDomain  string `form:"team_domain"`
	ResponseURL string `form:"response_url"`
	TriggerID   string `form:"trigger_id"`
	Text        string `form:"text"`
	Token       string `form:"token"`
	ChannelName string `form:"channel_name"`
	Command     string `form:"command"`
}

type fetchBorkResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
	Err          error  `json:"error,omitempty"`
}

func (r fetchBorkResponse) error() error { return r.Err }

func makeFetchBorkEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(fetchBorkRequest)
		bork, err := s.FetchBork(req)
		return fetchBorkResponse{
			ResponseType: "in_channel",
			Text:         bork,
			Err:          err,
		}, nil
	}
}
