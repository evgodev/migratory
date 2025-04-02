package migratory

const (
	migrationTypeGo  = "go"
	migrationTypeSQL = "sql"
	dialectPostgres  = "postgres"
)

type options struct {
	migrationType string
	directory     string
	dialect       string
	table         string

	forceUp bool
}

var defaultOpts = options{
	migrationType: migrationTypeGo,
	dialect:       dialectPostgres,
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

// WithForce is an OptionsFunc that sets the forceUp flag to true in the options configuration.
func WithForce() OptionsFunc {
	return func(o *options) { o.forceUp = true }
}

// SetTable updates the default table name used for managing migrations.
func SetTable(s string) { defaultOpts.table = s }

// SetSQLDirectory configures the default options to use the specified directory for SQL migration files.
func SetSQLDirectory(path string) {
	defaultOpts.migrationType = migrationTypeSQL
	defaultOpts.directory = path
}

func applyOptions(optionsFns []OptionsFunc) options {
	opts := defaultOpts
	for _, apply := range optionsFns {
		apply(&opts)
	}
	return opts
}
