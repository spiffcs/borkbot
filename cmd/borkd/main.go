package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/sparklycb/borkbot"
)

func main() {
	var (
		httpAddr          = flag.String("listen", ":9000", "HTTP listen and serve address for service")
		verificationToken = flag.String("verification_token", "", "Slack token used to verify requests come from slack")
	)
	flag.Parse()

	// Configure logger
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	httpLogger := log.With(logger, "component", "http")

	// Configure borkService
	var bs borkbot.Service
	bs = borkbot.NewService(*verificationToken)
	bs = borkbot.NewLoggingService(log.With(logger, "component", "borkbot"), bs)
	mux := http.NewServeMux()
	mux.Handle("/borkbot/v1/", borkbot.MakeHandler(bs, httpLogger))
	http.Handle("/", accessControl(mux))

	// Configure Listen and Serve
	errs := make(chan error, 2)
	errChan := make(chan error)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
