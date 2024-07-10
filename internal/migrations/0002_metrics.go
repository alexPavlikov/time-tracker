package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upMetrics, downMetrics)
}

func upMetrics(tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS metrics (
		"id" serial PRIMARY KEY NOT NULL,
		"user_id" uuid NOT NULL,
		"func_name" character varying NOT NULL,
		"time_micro" bigint NOT NULL);`

	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("migrate failed to create table metrics: %w", err)
	}

	return nil
}

func downMetrics(tx *sql.Tx) error {

	query := `DROP TABLE metrics;`

	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("migrate failed to drop table metrics: %w", err)
	}

	return nil
}
