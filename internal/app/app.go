package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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

type Application struct {
	router *chi.Mux
	cfg    *config.Config
	logger *slog.Logger
	server *http.Server
}

func New() (*Application, error) {
	app := &Application{}

	if err := app.initConfig(); err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}

	if err := app.initLogger(); err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	if err := app.initRouter(); err != nil {
		return nil, fmt.Errorf("failed to init router: %w", err)
	}

	return app, nil
}

func (a *Application) initConfig() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	a.cfg = cfg
	return nil
}

func (a *Application) initLogger() error {
	log := logger.New(a.cfg.Env)
	if log == nil {
		return fmt.Errorf("failed to create logger")
	}
	a.logger = log
	return nil
}

func (a *Application) initRouter() error {
	a.router = chi.NewRouter()

	// Инициализация зависимостей
	ctx := context.Background()

	// База данных
	client, err := postgres.NewClient(ctx, &a.cfg.Postgres)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Пользовательский сервис
	repo := repository.NewRepository(client, a.logger)
	srv := service.NewService(repo, a.logger)
	userHandler := handler.NewUserHandler(srv, a.logger)

	// WebSocket
	hub := websocket.NewHub(a.logger)
	wsHandler := websocket.NewHandlerWS(hub, srv, a.logger)
	go hub.Run()

	// Регистрация маршрутов
	route := routes.NewRoute(userHandler, wsHandler)
	route.RegisterRoutes(a.router)

	// Настройка HTTP сервера
	a.server = &http.Server{
		Addr:    a.cfg.HTTPServer.Addr,
		Handler: a.router,
	}

	return nil
}

func (a *Application) Run() error {
	a.logger.Info("Starting server", "addr", a.cfg.HTTPServer.Addr)

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down server")
	return a.server.Shutdown(ctx)
}
