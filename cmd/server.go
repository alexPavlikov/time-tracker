package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexPavlikov/time-tracker/internal/config"
	"github.com/alexPavlikov/time-tracker/internal/db"
	router "github.com/alexPavlikov/time-tracker/internal/server"
	postgres "github.com/alexPavlikov/time-tracker/internal/server/db"
	"github.com/alexPavlikov/time-tracker/internal/server/locations"
	"github.com/alexPavlikov/time-tracker/internal/server/service"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config load error: %w", err)
	}

	setupLogger(cfg.LogLevel)

	slog.Info("starting application listen on", "server", cfg.ServerToString())

	conn, err := db.Connect(cfg)
	if err != nil {
		return fmt.Errorf("failed connect to database: %w", err)
	}

	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			panic("failed to close db connection")
		}
	}()

	postgres := postgres.New(context.TODO(), conn)
	service := service.New(context.TODO(), postgres)
	handler := locations.New(*service)
	repo := router.New(*handler)

	repo.Build()

	if err := http.ListenAndServe(cfg.ServerToString(), nil); err != nil {
		return fmt.Errorf("listen and serve error: %w", err)
	}

	return nil
}

func setupLogger(logLevel string) {
	var log *slog.Logger
	switch logLevel {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	slog.SetDefault(log)
}
