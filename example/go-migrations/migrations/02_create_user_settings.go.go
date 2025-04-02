package migrations

import (
	"context"
	"database/sql"

	"github.com/korfairo/migratory"
)

func init() {
	migratory.AddMigrationNoTx(
		// Up migration - creates user_settings table and an index
		func(ctx context.Context, db *sql.DB) error {
			// Some operations like creating indexes or certain DDL operations
			// might need to be executed outside of a transaction in some databases
			_, err := db.ExecContext(ctx, `
				CREATE TABLE user_settings (
					user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
					theme VARCHAR(50) DEFAULT 'default',
					notifications_enabled BOOLEAN DEFAULT TRUE,
					last_login TIMESTAMP WITH TIME ZONE
				)
			`)
			if err != nil {
				return err
			}

			// Create an index (may require being outside a transaction in some DB systems)
			_, err = db.ExecContext(ctx, `
				CREATE INDEX idx_user_settings_theme ON user_settings(theme)
			`)
			return err
		},
		// Down migration - drops the table and its dependencies
		func(ctx context.Context, db *sql.DB) error {
			_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS user_settings`)
			return err
		},
	)
}
