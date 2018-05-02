package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cooldarkdryplace/debugserver"
)

// NewRedirectServer to handle redirects to HTTPS.
func NewRedirectServer(handler http.Handler) *http.Server {
	s := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	return s
}

func main() {
	http.Handle("/", debugserver.API())

	var (
		errChan    = make(chan error)
		signalChan = make(chan os.Signal, 1)
	)

	signal.Notify(signalChan, os.Interrupt)

	s := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	go func() {
		errChan <- s.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		log.Println(err)
	case <-signalChan:
		log.Println("Interrupt recieved. Graceful shutdown.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %s", err)
	}
}
