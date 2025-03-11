package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/linyows/ogenexample/api"
)

func main() {
	srv, err, closer := api.Server()
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: srv,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		if err := closer(context.TODO()); err != nil {
			log.Fatalf("Failed to connection close: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to server shutdown %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
