package migratory

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/korfairo/migratory/internal/migrator"
	"github.com/korfairo/migratory/internal/sqlmigration"
)

var ErrUnsupportedMigrationType = errors.New("migration type is unsupported")

// Up applies all available database migrations in order,
// using the given database connection and optional configurations.
func Up(db *sql.DB, opts ...OptionsFunc) (n int, err error) {
	ctx := context.Background()
	return UpContext(ctx, db, opts...)
}

// UpContext applies any pending database migrations using the provided context, database connection, and options.
// It returns the number of migrations applied and any error encountered during the process.
func UpContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc) (n int, err error) {
	option := applyOptions(opts)
	m, err := migrator.New(ctx, db, option.dialect, option.table)
	if err != nil {
		return 0, err
	}

	migrations, err := getMigrations(option.migrationType, option.directory)
	if err != nil {
		return 0, err
	}

	appliedCount, err := m.Up(ctx, migrations, db, option.forceUp)
	if err != nil {
		return appliedCount, err
	}

	return appliedCount, nil
}

// Down rolls back the most recently applied migration in the database.
// Accepts optional configuration via OptionsFunc.
func Down(db *sql.DB, opts ...OptionsFunc) error {
	ctx := context.Background()
	return DownContext(ctx, db, opts...)
}

// DownContext rolls back database migrations using the provided context, database connection,
// and optional configuration.
func DownContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc) error {
	return rollback(ctx, db, false, opts)
}

// Redo rolls back and re-applies the last migration in the database using the provided options.
func Redo(db *sql.DB, opts ...OptionsFunc) error {
	ctx := context.Background()
	return RedoContext(ctx, db, opts...)
}

// RedoContext re-applies the most recently rolled back migration within the provided context and database connection.
func RedoContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc) error {
	return rollback(ctx, db, true, opts)
}

// MigrationResult represents the result of a migration,
// including its ID, name, application status, and applied timestamp.
type MigrationResult struct {
	ID        int64
	Name      string
	IsApplied bool
	AppliedAt time.Time
}

// GetStatus retrieves the migration status from the database,
// including applied status and application time for each migration.
func GetStatus(db *sql.DB, opts ...OptionsFunc) ([]MigrationResult, error) {
	ctx := context.Background()
	return GetStatusContext(ctx, db, opts...)
}

// GetStatusContext retrieves the migration status from the database
// and returns a list of MigrationResult with their details.
func GetStatusContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc) ([]MigrationResult, error) {
	option := applyOptions(opts)
	m, err := migrator.New(ctx, db, option.dialect, option.table)
	if err != nil {
		return nil, err
	}

	migrations, err := getMigrations(option.migrationType, option.directory)
	if err != nil {
		return nil, err
	}

	results, err := m.GetStatus(ctx, migrations, db)
	if err != nil {
		return nil, err
	}

	migrationResults := make([]MigrationResult, 0, len(results))
	for _, r := range results {
		migrationResults = append(migrationResults, MigrationResult{
			ID:        r.ID,
			Name:      r.Name,
			IsApplied: !r.AppliedAt.IsZero(),
			AppliedAt: r.AppliedAt,
		})
	}

	return migrationResults, nil
}

// GetDBVersion retrieves the current database schema version based on the migrations table.
// The database version is represented by the ID of the last applied migration.
// Takes an *sql.DB instance and optional configuration through OptionsFunc.
// Returns the schema version as an int64 or an error if retrieval fails.
func GetDBVersion(db *sql.DB, opts ...OptionsFunc) (int64, error) {
	ctx := context.Background()
	return GetDBVersionContext(ctx, db, opts...)
}

// GetDBVersionContext retrieves the current database version by querying the migrations metadata table.
// The database version is represented by the ID of the last applied migration.
func GetDBVersionContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc) (int64, error) {
	option := applyOptions(opts)
	m, err := migrator.New(ctx, db, option.dialect, option.table)
	if err != nil {
		return -1, err
	}

	version, err := m.GetDBVersion(ctx, db)
	if err != nil {
		return -1, err
	}

	return version, nil
}

func getMigrations(migrationType, directory string) (m migrator.Migrations, err error) {
	switch migrationType {
	case migrationTypeGo:
		m, err = registerGoMigrations(globalGoMigrations)
	case migrationTypeSQL:
		m, err = sqlmigration.SeekMigrations(directory, nil)
	default:
		return nil, ErrUnsupportedMigrationType
	}
	return m, err
}

func rollback(ctx context.Context, db *sql.DB, redo bool, opts []OptionsFunc) error {
	option := applyOptions(opts)
	m, err := migrator.New(ctx, db, option.dialect, option.table)
	if err != nil {
		return err
	}

	migrations, err := getMigrations(option.migrationType, option.directory)
	if err != nil {
		return err
	}

	if err = m.Down(ctx, migrations, db, redo); err != nil {
		return err
	}

	return nil
}
