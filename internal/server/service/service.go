package service

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/time-tracker/internal/domain"
	postgres "github.com/alexPavlikov/time-tracker/internal/server/db"
	"github.com/gofrs/uuid/v5"
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

func (s *Services) Add(user domain.User) (id int64, err error) {
	id, err = s.postgres.Insert(s.ctx, user)
	if err != nil {
		return id, fmt.Errorf("insert user error: %w", err)
	}
	return id, nil
}

func (s *Services) Update(user domain.User) error {
	err := s.postgres.Update(s.ctx, user)
	if err != nil {
		return fmt.Errorf("update user error: %w", err)
	}
	return nil
}

func (s *Services) Delete(id int) error {
	err := s.postgres.Delete(s.ctx, id)
	if err != nil {
		return fmt.Errorf("delete user error: %w", err)
	}
	return nil
}

func (s *Services) GetOne(id int64) (domain.User, error) {
	user, err := s.postgres.SelectOne(s.ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("select one user error with id: %d: %w", id, err)
	}
	return user, nil
}

func (s *Services) Get() ([]domain.User, error) {

	var test = domain.UserSortParameters{
		PassportSeries: 312,
		PassportNumber: 534,
		Surname:        "312",
		Name:           "31321",
		Patronymic:     "312",
		Address:        "31",
	}

	users, err := s.postgres.Select(s.ctx, test)
	if err != nil {
		return nil, fmt.Errorf("select users error: %w", err)
	}
	return users, nil
}

func (s *Services) GetPag(limit int, offset int) ([]domain.User, error) {
	users, err := s.postgres.PagSelect(s.ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select pagination users error: %w", err)
	}

	return users, nil
}

func (s *Services) AddMetrics(m domain.Metrics) error {
	err := s.postgres.InsertMetrics(s.ctx, m)
	if err != nil {
		return fmt.Errorf("failed to insert metrics: %w", err)
	}
	return nil
}

func (s *Services) GetMetrics(uuid uuid.UUID) ([]domain.Metrics, error) {
	metrics, err := s.postgres.SelectMetrics(s.ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return metrics, nil
}
