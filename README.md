# Migratory

A lightweight and flexible database migration library and CLI tool with minimal dependencies.

## Features

- Supports **PostgreSQL** and **ClickHouse**. You can easily add another DB, supported by `database/sql` standard go package
- Usable as a library (Go package) or CLI tool
- **Go migrations**: code-based migrations with support for transactional (`tx`) or non-transactional (`db`) connections in the library
- **SQL migrations**: `.sql` files with raw SQL statements, supported in both library and CLI
- Lazy SQL parsing: SQL migrations are parsed only when applied, ensuring efficient usage

## Library Usage

Install the package:

```shell
go get github.com/evgodev/migratory
```

#### Migration Registration

| Function | Description |
| --- | --- |
| `AddMigration(up, down GoMigrateFn)` | Registers a transaction-based migration with `up` and `down` functions. Each function receives a `tx *sql.Tx` parameter that represents a transaction, allowing you to execute multiple SQL operations as a single atomic unit. If any error occurs, the entire transaction is rolled back. Use this for operations that require data consistency or can safely be run within a transaction. |
| `AddMigrationNoTx(up, down GoMigrateNoTxFn)` | Registers a non-transactional migration with `up` and `down` functions. Each function receives a `db *sql.DB` parameter that represents a direct database connection without transaction support. This is useful for operations that cannot be run within a transaction in PostgreSQL (like some DDL operations, creating indexes, etc.) or when you need to manage transactions manually. |
| `SetSQLDirectory(path string)` | Sets the directory where `.sql` migration files are located. SQL migrations are automatically parsed and registered from this directory |

##### Migration Function Types

```go
// Transaction-based migration function definition
type GoMigrateFn func(ctx context.Context, tx *sql.Tx) error

// Non-transaction based migration function definition  
type GoMigrateNoTxFn func(ctx context.Context, db *sql.DB) error
```

##### Example Usage

**Code-based transaction migration:**
```go
migratory.AddMigration(
    // Up migration (using tx)
    func(ctx context.Context, tx *sql.Tx) error {
        _, err := tx.ExecContext(ctx, `
            CREATE TABLE users (
                id SERIAL PRIMARY KEY,
                username VARCHAR(100) NOT NULL UNIQUE
            )
        `)
        return err
    },
    // Down migration (using tx)
    func(ctx context.Context, tx *sql.Tx) error {
        _, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS users`)
        return err
    },
)
```

**Code-based non-transaction migration:**
```go
migratory.AddMigrationNoTx(
    // Up migration (using db directly)
    func(ctx context.Context, db *sql.DB) error {
        // Creating an index might need to be outside a transaction
        _, err := db.ExecContext(ctx, `
            CREATE INDEX idx_users_username ON users(username)
        `)
        return err
    },
    // Down migration (using db directly)
    func(ctx context.Context, db *sql.DB) error {
        _, err := db.ExecContext(ctx, `DROP INDEX IF EXISTS idx_users_username`)
        return err
    },
)
```

**SQL migration file:**
```sql
-- +migrate up
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);

-- Create an index on user_id for faster queries
CREATE INDEX idx_posts_user_id ON posts(user_id);

-- +migrate down
-- The down migration should revert everything done in the up migration
-- in reverse order (drop index first, then table)
DROP INDEX IF EXISTS idx_posts_user_id;
DROP TABLE IF EXISTS posts;
```

When using SQL migrations:
1. Save files with a numeric prefix (like `01_create_users.sql`, `02_create_posts.sql`) to control execution order
2. Use the `-- +migrate up` and `-- +migrate down` comments to separate migration sections
3. Set the SQL directory with `migratory.SetSQLDirectory("./migrations")`
4. Each SQL file is automatically registered as a migration


#### Migration Operations

| Function | Description |
|----------|-------------|
| `Up(db *sql.DB, opts ...OptionsFunc) (n int, err error)` | Applies all unapplied migrations |
| `Down(db *sql.DB, opts ...OptionsFunc) error` | Rolls back the last applied migration |
| `Redo(db *sql.DB, opts ...OptionsFunc) error` | Rolls back and reapplies the last migration |
| `GetStatus(db *sql.DB, opts ...OptionsFunc) ([]MigrationResult, error)` | Returns the status of all migrations |
| `GetDBVersion(db *sql.DB, opts ...OptionsFunc) (int64, error)` | Returns the current migration version (ID of the last applied migration) |

#### Context-Aware Operations

Each operation above has a context-aware equivalent:

| Function | Description |
|----------|-------------|
| `UpContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc)` | Context-aware version of `Up()` |
| `DownContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc)` | Context-aware version of `Down()` |
| `RedoContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc)` | Context-aware version of `Redo()` |
| `GetStatusContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc)` | Context-aware version of `GetStatus()` |
| `GetDBVersionContext(ctx context.Context, db *sql.DB, opts ...OptionsFunc)` | Context-aware version of `GetDBVersion()` |

## CLI Usage

Install the CLI tool:

```shell
go install github.com/evgodev/migratory/cmd/migratory
```

### Available Commands

| Command | Description |
|---------|-------------|
| `create` | Create .sql or .go migration template |
| `dbversion` | Show the DB version (ID of the last applied migration) |
| `down` | Rollback the last applied migration |
| `help` | Display help about any command |
| `redo` | Rollback and apply the last migration again |
| `status` | Show migration statuses |
| `up` | Apply all unapplied migrations |

### Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-c, --config string` | Path to YAML config | |
| `-d, --db string` | Database connection string | |
| `--dir string` | Directory with .sql migration files | `.` |
| `-h, --help` | Show help for migratory | |
| `-t, --table string` | Name of the migrations table | `migrations` |

You can find detailed information about all commands using the `--help` or `-h` flag.

### Configuration File

All commands can use a YAML configuration file. Create a configuration file and pass its path using the `-c ./path/` flag.

Example configuration:

```yaml
directory: /path/to/directory
dsn: postgres://user:password@localhost:5432/my_db
table: migrations
```

## Changelog
### [1.0.0] - 2025-01-28
- Initial release of Migratory.