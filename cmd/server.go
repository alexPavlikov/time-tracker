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
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	_ "github.com/alexPavlikov/time-tracker/internal/migrations"
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

	UUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("create user UUID error: %w", err)
	}

	conn, err := db.Connect(cfg)
	if err != nil {
		return fmt.Errorf("failed connect to database: %w", err)
	}

	defer conn.Close()

	//----------------------------------------------
	db := stdlib.OpenDBFromPool(conn)

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to start migrations: %w", err)
	}
	//----------------------------------------------

	ctx := context.Background()

	ctx = context.WithValue(ctx, "UUID", UUID.String())

	slog.Info("generate user UUID", "uuid", UUID.String())

	postgres := postgres.New(context.TODO(), conn)
	service := service.New(context.TODO(), postgres)
	handler := locations.New(ctx, *service)
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
