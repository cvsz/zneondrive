package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/httpapi"
	"github.com/cvsz/zneondrive/services/game-api/internal/store"
)

func main() {
	ctx := context.Background()
	databaseURL := env("DATABASE_URL", "postgres://zneondrive:zneondrive@127.0.0.1:55432/zneondrive?sslmode=disable")
	listenAddr := env("LISTEN_ADDR", ":8080")
	redisAddr := strings.TrimSpace(os.Getenv("REDIS_ADDR"))
	trustedProxyCIDRs := strings.TrimSpace(os.Getenv("TRUSTED_PROXY_CIDRS"))
	gameServerKey := strings.TrimSpace(os.Getenv("GAME_SERVER_SHARED_KEY"))
	if len(gameServerKey) < 32 {
		log.Fatal("GAME_SERVER_SHARED_KEY must be configured with at least 32 characters")
	}

	proxyPolicy, err := httpapi.ParseTrustedProxyCIDRs(trustedProxyCIDRs)
	if err != nil {
		log.Fatalf("TRUSTED_PROXY_CIDRS: %v", err)
	}
	if trustedProxyCIDRs == "" {
		log.Printf("trusted proxy identity disabled; forwarded headers will be ignored")
	} else {
		log.Printf("trusted proxy identity enabled from configured CIDR allowlist")
	}

	db, err := store.OpenPostgres(ctx, databaseURL)
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}
	defer db.Close()

	if err = db.EnsureSchema(ctx); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}
	if err = db.EnsureRaceSchema(ctx); err != nil {
		log.Fatalf("ensure race schema: %v", err)
	}

	apiHandler := httpapi.New(db, gameServerKey)
	rateLimitedHandler := httpapi.NewRateLimitedHandlerWithTrustedProxies(apiHandler, proxyPolicy)
	if redisAddr != "" {
		rateLimitedHandler = httpapi.NewDistributedRateLimitedHandlerWithTrustedProxies(apiHandler, redisAddr, proxyPolicy)
		log.Printf("distributed rate limiting enabled via Redis")
	} else {
		log.Printf("WARN: REDIS_ADDR is not configured; rate limiting is process-local only")
	}
	server := &http.Server{
		Addr:              listenAddr,
		Handler:           rateLimitedHandler,
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
