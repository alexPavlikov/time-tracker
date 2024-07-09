package postgres

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/time-tracker/internal/domain"
	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	ctx context.Context
	DB  *pgx.Conn
}

func New(ctx context.Context, DB *pgx.Conn) *Postgres {
	return &Postgres{
		ctx: ctx,
		DB:  DB,
	}
}

func (p *Postgres) Insert(ctx context.Context, user domain.User) (int, error) {
	query := `
	INSERT INTO public."users" (passport_series, passport_number, surname, name, patronymic, address) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`

	row := p.DB.QueryRow(ctx, query, fmt.Sprint(user.PassportSeries), fmt.Sprint(user.PassportNumber), user.Surname, user.Name, user.Patronymic, user.Address)

	err := row.Scan(&user.ID)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
