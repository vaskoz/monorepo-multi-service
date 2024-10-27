package app

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	out    = os.Stdout
	err    = os.Stderr
	in     = os.Stdin
	args   = os.Args
	exit   = os.Exit
	sigint = make(chan os.Signal, 1)
)

func RealMain() {
	var srv http.Server

	ll := log.New(err, "public-api: ", log.LstdFlags)

	signal.Notify(sigint, os.Interrupt)

	idleConnsClosed := make(chan struct{})

	go func() {
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			ll.Printf("HTTP server Shutdown error: %v", err)
		}
		close(idleConnsClosed)
	}()

	usersClient := &http.Client{}
	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		req, _ := http.NewRequest("GET", "http://localhost:8081/users", nil)
		if resp, err := usersClient.Do(req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			defer resp.Body.Close()
			w.WriteHeader(resp.StatusCode)
			payload, _ := io.ReadAll(resp.Body)
			w.Write(payload)
		}
	})

	srv.Handler = mux
	srv.Addr = ":8080"

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	ll.Printf("reached the waiting state")
	<-idleConnsClosed
}
