package migrations

import (
	"context"
	"database/sql"

	"github.com/korfairo/migratory"
)

func init() {
	migratory.AddMigration(
		// Up migration - creates users table
		func(ctx context.Context, tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, `
				CREATE TABLE users (
					id SERIAL PRIMARY KEY,
					username VARCHAR(100) NOT NULL UNIQUE,
					email VARCHAR(255) NOT NULL UNIQUE,
					password_hash VARCHAR(255) NOT NULL,
					created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMP WITH TIME ZONE
				)
			`)
			return err
		},
		// Down migration - drops users table
		func(ctx context.Context, tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS users`)
			return err
		},
	)
}
