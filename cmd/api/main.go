package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/artrctx/noliteo-core/internal/server"
)

func gracefulShutdown(srv *http.Server, done chan bool) {
	// listen to interrupt signal from os
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal
	<-ctx.Done()

	log.Println("shutting down. press Ctrl+C again to force quit")
	stop() // Allow Ctrl+C to force shutdown

	// 5 second max request handling time
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func main() {
	srv := server.NewServer()

	// app close channel
	done := make(chan bool, 1)

	// Handle graceful shutdown in a seperate goroutine
	go gracefulShutdown(srv, done)

	log.Printf("quoin server starting at %v", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Gracefully shutdown server.")
}
