package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/httpapi"
	"github.com/cvsz/zneondrive/services/game-api/internal/store"
)

func main() {
	ctx := context.Background()
	databaseURL := env("DATABASE_URL", "postgres://zneondrive:zneondrive@127.0.0.1:55432/zneondrive?sslmode=disable")
	listenAddr := env("LISTEN_ADDR", ":8080")

	db, err := store.OpenPostgres(ctx, databaseURL)
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}
	defer db.Close()

	if err = db.EnsureSchema(ctx); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}

	server := &http.Server{
		Addr:              listenAddr,
		Handler:           httpapi.New(db),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("zNeonDrive service plane listening on %s", listenAddr)
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			log.Fatalf("http server: %v", serveErr)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
