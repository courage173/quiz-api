package main

import (
	"context"

	"flag"

	"os"

	"fmt"

	"time"

	"os/signal"

	"sync/atomic"

	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"

	"github.com/go-ozzo/ozzo-routing/v2/content"

	"github.com/go-ozzo/ozzo-routing/v2/cors"

	"github.com/courage173/quiz-api/pkg/accesslog"

	"github.com/courage173/quiz-api/pkg/log"

	"github.com/courage173/quiz-api/internal/errors"

	"github.com/courage173/quiz-api/internal/healthcheck"
)

var (
	Version    string = "1.0.0"
	listenAddr string
	healthy    int32
)

func main() {
	flag.StringVar(&listenAddr, "listen-addr", "localhost:4000", "server listen address")
	flag.Parse()

	// create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	fmt.Println(listenAddr)
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      buildHandler(logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Errorf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	fmt.Println("Server is ready to handle requests at", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	fmt.Println("Server stopped")
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger) http.Handler {
	router := routing.New()

	router.Use(
		accesslog.Handler(logger),
		errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	healthcheck.RegisterHandlers(router, Version)

	// rg := router.Group("/v1")

	// authHandler := auth.Handler(cfg.JWTSigningKey)

	// album.RegisterHandlers(rg.Group(""),
	// 	album.NewService(album.NewRepository(db, logger), logger),
	// 	authHandler, logger,
	// )

	// auth.RegisterHandlers(rg.Group(""),
	// 	auth.NewService(jwtSigningKey, expiry, logger, users.NewRepository(db, logger)),
	// 	logger,
	// )

	return router
}
