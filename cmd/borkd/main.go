package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"crypto/tls"

	"github.com/go-kit/kit/log"
	"github.com/sparklycb/borkbot"
)

func main() {
	var (
		httpAddr          = flag.String("listen", ":8080", "HTTP listen and serve address for service")
		verificationToken = flag.String("verification_token", "", "Slack token used to verify requests come from slack")
	)
	flag.Parse()

	// Configure logger
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	httpLogger := log.With(logger, "component", "http")

	// Configure borkServiceHandlers
	var bs borkbot.Service
	bs = borkbot.NewService(*verificationToken)
	bs = borkbot.NewLoggingService(log.With(logger, "component", "borkbot"), bs)
	mux := http.NewServeMux()
	mux.Handle("/borkbot/v1/", borkbot.MakeHandler(bs, httpLogger))
	mux.Handle("/borkbot/", accessControl(mux))

	// Configure TLS
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	// Configure Server
	srv := &http.Server{
		Addr:         *httpAddr,
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	// Configure Listen and Serve
	errs := make(chan error, 2)
	errChan := make(chan error)
	// Server goroutine
	go func() {
		logger.Log("transport", "https", "address", *httpAddr, "msg", "listening")
		errs <- srv.ListenAndServeTLS("/secure/cert.pem", "/secure/key.pem")
	}()
	// operator cancel go routines
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
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
