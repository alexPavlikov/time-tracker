package service

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/time-tracker/internal/domain"
	postgres "github.com/alexPavlikov/time-tracker/internal/server/db"
)

type Services struct {
	ctx      context.Context
	postgres *postgres.Postgres
}

func New(ctx context.Context, postgres *postgres.Postgres) *Services {
	return &Services{
		ctx:      ctx,
		postgres: postgres,
	}
}

func (s *Services) Add(user domain.User) (id int, err error) {

	id, err = s.postgres.Insert(s.ctx, user)
	if err != nil {
		return id, fmt.Errorf("insert user error: %w", err)
	}

	return id, nil
}
