package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/QuUteO/video-communication/internal/config"
	"github.com/QuUteO/video-communication/internal/logger"
	"github.com/QuUteO/video-communication/internal/routes"
	"github.com/QuUteO/video-communication/internal/user/handler"
	"github.com/QuUteO/video-communication/internal/user/repository"
	"github.com/QuUteO/video-communication/internal/user/service"
	"github.com/QuUteO/video-communication/internal/websocket"
	"github.com/QuUteO/video-communication/pkg/db"
	"github.com/go-chi/chi/v5"
)

// go run ./cmd/main.go to run the application
func main() {
	// ctx
	ctx := context.Background()

	// router
	r := chi.NewRouter()

	// config initialization
	cfg, err := config.New()
	if err != nil {
		fmt.Println("failed to load config:", err)
	}

	// logger initialization
	log := logger.New(cfg.Env)

	// data base initialization
	client, err := postgres.NewClient(ctx, &cfg.Postgres)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
	}

	// service
	repo := repository.NewRepository(client, log)
	srv := service.NewService(repo, log)
	handle := handler.NewUserHandler(srv, log)

	// websocket
	hub := websocket.NewHub()
	wsHandle := websocket.NewWebSocket(hub, log, &srv)
	go hub.Run()

	// routes
	route := routes.NewRoute(handle, wsHandle)
	route.RegisterRoutes(r)

	log.Info("Starting server")
	if err := http.ListenAndServe(cfg.HTTPServer.Addr, r); err != nil {
		log.Error("failed to start http server", "error", err)
	}
	// graceful shutdown
}
