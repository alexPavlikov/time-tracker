package db

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/time-tracker/internal/config"
	"github.com/jackc/pgx/v5"
)

func Connect(cfg *config.Config) (conn *pgx.Conn, err error) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.PostgreUser, cfg.PostgrePassword, cfg.PostgresPath, cfg.PostgresPort, cfg.PostgreDatabaseName)
	conn, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
