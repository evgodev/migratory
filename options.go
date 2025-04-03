package migratory

import "github.com/korfairo/migratory/internal/migrator"

const (
	Postgres   Dialect = migrator.DialectPostgres
	ClickHouse Dialect = migrator.DialectClickHouse

	migrationTypeGo  = "go"
	migrationTypeSQL = "sql"
)

// Dialect determines how the migrations table is managed based on the database system.
type Dialect = string

type options struct {
	migrationType string
	directory     string
	dialect       string
	table         string

	forceUp bool
}

var defaultOpts = options{
	migrationType: migrationTypeGo,
	dialect:       Postgres,
	directory:     ".",
	table:         "migrations",
	forceUp:       false,
}

type OptionsFunc func(o *options)

// WithGoMigration sets the migration type to "go" for the specified options configuration.
func WithGoMigration() OptionsFunc {
	return func(o *options) { o.migrationType = migrationTypeGo }
}

// WithSQLMigrationDir sets the migration type to SQL and specifies the directory containing migration files.
func WithSQLMigrationDir(d string) OptionsFunc {
	return func(o *options) { o.migrationType = migrationTypeSQL; o.directory = d }
}

// WithTable configures a custom table name for tracking migrations within the database.
func WithTable(n string) OptionsFunc {
	return func(o *options) { o.table = n }
}

// WithForce allows the migrator to apply any unapplied migrations out of their original sequence,
// potentially resulting in migrations being applied in a non-linear order.
// For example, if there are migrations with numbers 1, 2, 3, and migration 2 was not applied before,
// this option allows you to apply it. Otherwise, the migrator will return an error migrator.ErrDirtyMigrations.
func WithForce() OptionsFunc {
	return func(o *options) { o.forceUp = true }
}

// WithDialect sets the database dialect.
func WithDialect(d Dialect) OptionsFunc {
	return func(o *options) { o.dialect = d }
}

// SetTable updates the default table name used for managing migrations.
func SetTable(s string) { defaultOpts.table = s }

// SetSQLDirectory configures the default options to use the specified directory for SQL migration files.
func SetSQLDirectory(path string) {
	defaultOpts.migrationType = migrationTypeSQL
	defaultOpts.directory = path
}

// SetDialect sets the default SQL dialect for database migrations.
// The chosen dialect determines how the migrations table is managed
// based on the database system (e.g., Postgres, ClickHouse).
func SetDialect(d Dialect) {
	defaultOpts.dialect = d
}

func applyOptions(optionsFns []OptionsFunc) options {
	opts := defaultOpts
	for _, apply := range optionsFns {
		apply(&opts)
	}
	return opts
}
