package main

import (
	"context"
	"encoding/json"
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

func main() {
	var srv http.Server

	ll := log.New(err, "user-service: ", log.LstdFlags)

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

	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users := GetAllUsers()
		if jsonStr, err := json.Marshal(users); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(jsonStr)
			w.WriteHeader(http.StatusOK)
		}
	})

	srv.Handler = mux

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	ll.Printf("reached the waiting state")
	<-idleConnsClosed
}