package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/alexPavlikov/time-tracker/internal/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	ctx context.Context
	DB  *pgxpool.Pool
}

func New(ctx context.Context, DB *pgxpool.Pool) *Postgres {
	return &Postgres{
		ctx: ctx,
		DB:  DB,
	}
}

// Функция выборки одного user по id
func (p *Postgres) SelectOne(ctx context.Context, id int64) (user domain.User, err error) {
	query := `SELECT id, passport_series, passport_number, surname, name, patronymic, address FROM public."users"
	WHERE id = $1`

	row := p.DB.QueryRow(ctx, query, id)

	err = row.Scan(&user.ID, &user.PassportSeries, &user.PassportNumber, &user.Surname, &user.Name, &user.Patronymic, &user.Address)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// Функция выборки всех user с возможностью сортировки
func (p *Postgres) Select(ctx context.Context, options ...domain.UserSortParameters) (users []domain.User, err error) {
	query := `SELECT id, passport_series, passport_number, surname, name, patronymic, address FROM public."users"`

	if options != nil {
		opt := validateSortParameters(options[0])
		query += " WHERE " + opt
	}

	rows, err := p.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.PassportSeries, &user.PassportNumber, &user.Surname, &user.Name, &user.Patronymic, &user.Address)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Функция добавления пользователя с возвращением его id
func (p *Postgres) Insert(ctx context.Context, user domain.User) (int64, error) {
	query := `
	INSERT INTO public."users" (passport_series, passport_number, surname, name, patronymic, address) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`

	row := p.DB.QueryRow(ctx, query, user.PassportSeries, user.PassportNumber, user.Surname, user.Name, user.Patronymic, user.Address)

	err := row.Scan(&user.ID)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// Функция удаления пользователя по id
func (p *Postgres) Delete(ctx context.Context, id int) error {
	query := `DELETE public."users" WHERE id = $1`

	_, err := p.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// Функция полного обновления пользователя
func (p *Postgres) Update(ctx context.Context, us domain.User) error {
	query := `
	UPDATE public."users" SET passport_series = $1, passport_number= $2, surname= $3, name= $4, patronymic= $5, address= $6
	WHERE id = $7
	`

	_ = p.DB.QueryRow(ctx, query, us.PassportSeries, us.PassportNumber, us.Surname, us.Name, us.Patronymic, us.Address, us.ID)

	return nil
}

// Функция выборки всех пользователей с пагинацией и сортировкой
func (p *Postgres) PagSelect(ctx context.Context, limit int, offset int, options ...domain.UserSortParameters) (users []domain.User, err error) {
	query := `SELECT id, passport_series, passport_number, surname, name, patronymic, address FROM public."users"`

	if options != nil {
		opt := validateSortParameters(options[0])
		if opt != "" {
			query += " WHERE " + opt
		}
	}

	query += " ORDER BY id LIMIT $1 OFFSET $2"

	rows, err := p.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.PassportSeries, &user.PassportNumber, &user.Surname, &user.Name, &user.Patronymic, &user.Address)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Функция парсинга параметров сортировки
func validateSortParameters(options domain.UserSortParameters) string {
	var void domain.UserSortParameters
	if options == void {
		return ""
	}

	var res string
	var count int

	if options.Address != "" {
		if count != 0 {
			res += "AND "
		}
		res += fmt.Sprintf(`address ILIKE '%%%s%%' `, options.Address)
		count++
	}

	if options.Name != "" {
		if count != 0 {
			res += "AND "
		}
		res += fmt.Sprintf(`name ILIKE '%%%s%%' `, options.Name)
		count++
	}

	if options.Patronymic != "" {
		if count != 0 {
			res += "AND "
		}
		res += fmt.Sprintf(`patronymic ILIKE '%%%s%%' `, options.Patronymic)
		count++
	}

	if options.Surname != "" {
		if count != 0 {
			res += "AND "
		}
		res += fmt.Sprintf(`surname ILIKE '%%%s%%' `, options.Surname)
		count++
	}

	if options.PassportNumber != 0 {
		if count != 0 {
			res += "AND "
		}
		res += fmt.Sprintf("passport_number = %d ", options.PassportNumber)
		count++
	}

	if options.PassportSeries != 0 {
		if count != 0 {
			res += "AND "
		}
		res += fmt.Sprintf("passport_series = %d ", options.PassportSeries)
		count++
	}

	return res
}

// Функция добавления metrics
func (p *Postgres) InsertMetrics(ctx context.Context, m domain.Metrics) error {
	query := `INSERT INTO public."metrics" (user_id, func_name, time_micro) VALUES ($1, $2, $3)`

	time := m.Time.Microseconds() / 1000

	_ = p.DB.QueryRow(ctx, query, m.User_ID.String(), m.FuncName, time)
	return nil
}

// Функция выборки всех метрив по uuid
func (p *Postgres) SelectMetrics(ctx context.Context, user uuid.UUID) ([]domain.Metrics, error) {
	query := `SELECT user_id, func_name, time_micro FROM public."metrics" WHERE user_id = $1 ORDER BY time_micro DESC`

	rows, err := p.DB.Query(ctx, query, user.String())
	if err != nil {
		return nil, err
	}

	var metrics = make([]domain.Metrics, 0)

	for rows.Next() {
		var m domain.Metrics
		var t int
		err := rows.Scan(&m.User_ID, &m.FuncName, &t)
		if err != nil {
			fmt.Println("error", err)
			return nil, err
		}

		m.Time = time.Duration(t)

		metrics = append(metrics, m)
	}

	return metrics, nil
}
