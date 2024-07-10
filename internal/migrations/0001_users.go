package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upUsers, downUsers)
}

func upUsers(tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS users (
		"id" serial PRIMARY KEY NOT NULL,
		"surname" character varying NOT NULL,
		"name" character varying NOT NULL,
		"patronymic" character varying NOT NULL,
		"address" character varying NOT NULL,
		"passport_series" bigint NOT NULL,
		"passport_number" bigint NOT NULL);`

	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("migrate failed to create table users: %w", err)
	}

	return nil
}

func downUsers(tx *sql.Tx) error {

	query := `DROP TABLE users;`

	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("migrate failed to drop table users: %w", err)
	}

	return nil
}
