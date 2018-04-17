package borkbot

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/form"
	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	fetchBorkHandler := kithttp.NewServer(
		makeFetchBorkEndpoint(s),
		decodeFetchBorkRequest,
		encodeResponse,
		opts...,
	)

	healthHandler := kithttp.NewServer(
		makeHealthEndpoint(s),
		decodeHealthRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/borkbot/v1/health", healthHandler).Methods("GET")
	r.Handle("/borkbot/v1/bork", fetchBorkHandler).Methods("POST")
	return r
}

var errBadRoute = errors.New("bad route")

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	// based on the type of error we can set different status codes in this block
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch e := err.Error(); e {
	case "why you try to be a slack?":
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func decodeFetchBorkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var br fetchBorkRequest
	r.ParseForm()
	decoder := form.NewDecoder()
	err := decoder.Decode(&br, r.Form)
	if err != nil {
		return nil, err
	}
	return br, nil
}

func decodeHealthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return healthRequest{}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// SlackForm contains all of the information that comes in from a slack POST request
// Requests from slack come as application/x-www-form-urlencoded so we need a way to
// decode these values when receving the request
type SlackForm struct {
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
