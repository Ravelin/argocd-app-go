package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/health", healthHandler)
	router.HandleFunc("/ready", readinessHandler)
	router.HandleFunc("/ping", pingHandler)

	http.Handle("/", router)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":9010",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interuptChan := make(chan os.Signal, 1)
	signal.Notify(interuptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until receive signal
	<-interuptChan

	// Create a deadline and wait
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting Down")
	os.Exit(0)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
