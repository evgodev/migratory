# migratory
## minimalistic database migration library and CLI tool

Supports PostgreSQL. Works as a library (go package) and CLI tool.

## As library

Works with **database/sql** standard package.

```shell
go get github.com/korfairo/migratory
```

### API

Register your `.go` migrations with functions:
```
- func AddMigration(up, down GoMigrateFn)
- func AddMigrationNoTx(up, down GoMigrateNoTxFn)
````
You can also use `.sql` migrations. Set directory with migrations files:

```
- func SetSQLDirectory(path string)
```

Or use `OptionsFunc` `WithSQLMigrationDir(d string)` in the next commands.

Manage your migrations with functions:

```
- func Up(db *sql.DB, opts ...OptionsFunc) (n int, err error)
- func Down(db *sql.DB, opts ...OptionsFunc) error
- func Redo(db *sql.DB, opts ...OptionsFunc) error
- func GetStatus(db *sql.DB, opts ...OptionsFunc) ([]MigrationResult, error)
- func GetDBVersion(db *sql.DB, opts ...OptionsFunc) (int64, error)
```

Or with your context:

```
- func UpContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc) (n int, err error)
- func DownContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc) error
- func RedoContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc) error
- func GetStatusContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc) ([]MigrationResult, error)
- func GetDBVersionContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc) (int64, error)
```

## As CLI tool

```shell
go install github.com/korfairo/migratory/cmd/migratory
```

### Usage

```
Usage:
  migratory [command]

Available Commands:
  create      Creates .sql or .go migration template
  dbversion   Shows the DB version (id of the last applied migration
  down        Rollback last applied migration
  help        Help about any command
  redo        Rollbacks and applies again last migration
  status      Shows migration statuses
  up          Up all unapplied migrations

Flags:
  -c, --config string   path to yaml config
  -d, --db string       database connection string
      --dir string      directory with .sql migration files (default ".")
  -h, --help            help for migratory
  -s, --schema string   name of database schema with migrations table (default "public")
  -t, --table string    name of migrations table (default "migrations")
```

You can find information about all commands and their usage with --help or -h flag.

All commands works with config (.yml file). Create a configuration file and pass its path with the `-c ./path/` flag.

Config example:
```yaml
directory: /path/to/directory
dsn: postgres://user:password@localhost:5432/my_db
schema: public
table: migrations
```