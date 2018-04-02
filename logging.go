package borkbot

import (
	"time"

	kitlog "github.com/go-kit/kit/log"
)

type loggingService struct {
	logger kitlog.Logger
	Service
}

// NewLoggingService returns a new instance of a loggin Service.
func NewLoggingService(logger kitlog.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) FetchBork(req fetchBorkRequest) (text string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "bork",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.FetchBork(req)
}
