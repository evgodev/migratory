package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	Postgres   Dialect = "postgres"
	ClickHouse Dialect = "clickhouse"
)

var (
	ErrUnsupportedDialect = errors.New("unsupported dialect")
	ErrNoRows             = errors.New("no rows in migrations table")
)

type Store struct {
	tableName    string
	queryManager QueryBuilder
}

type Dialect = string

type QueryBuilder interface {
	MigrationsTableExists(tableName string) string
	CreateMigrationsTable(tableName string) string
	InsertMigration(tableName string) string
	DeleteMigration(tableName string) string
	ListMigrations(tableName string) string
	SelectLastMigrationID(tableName string) string
}

type MigrationResult struct {
	ID        int64
	Name      string
	AppliedAt time.Time
}

type database interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func New(dbDialect, tableName string) (*Store, error) {
	var queryBuilder QueryBuilder

	switch dbDialect {
	case Postgres:
		queryBuilder = &postgresQueryBuilder{}
	case ClickHouse:
		queryBuilder = &clickhouseQueryBuilder{}
	default:
		return nil, ErrUnsupportedDialect
	}

	return &Store{
		queryManager: queryBuilder,
		tableName:    tableName,
	}, nil
}

func (s Store) MigrationsTableExists(ctx context.Context, db database) (bool, error) {
	q := s.queryManager.MigrationsTableExists(s.tableName)
	row := db.QueryRowContext(ctx, q)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to scan row: %w", err)
	}

	return exists, nil
}

func (s Store) CreateMigrationsTable(ctx context.Context, db database) error {
	q := s.queryManager.CreateMigrationsTable(s.tableName)
	_, err := db.ExecContext(ctx, q)
	return err
}

func (s Store) InsertMigration(ctx context.Context, db database, migrationName string, id int64) error {
	q := s.queryManager.InsertMigration(s.tableName)
	_, err := db.ExecContext(ctx, q, id, migrationName)
	return err
}

func (s Store) DeleteMigration(ctx context.Context, db database, id int64) error {
	q := s.queryManager.DeleteMigration(s.tableName)
	_, err := db.ExecContext(ctx, q, id)
	return err
}

func (s Store) SelectLastID(ctx context.Context, db database) (int64, error) {
	q := s.queryManager.SelectLastMigrationID(s.tableName)
	row := db.QueryRowContext(ctx, q)

	var id int64
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, ErrNoRows
		}
		return -1, fmt.Errorf("failed to scan migration id: %w", err)
	}
	if err := row.Err(); err != nil {
		return -1, fmt.Errorf("an error occurred during row scanning: %w", err)
	}

	return id, nil
}

func (s Store) ListMigrations(ctx context.Context, db database) ([]MigrationResult, error) {
	q := s.queryManager.ListMigrations(s.tableName)
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query listMigrations: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var migrations []MigrationResult
	for rows.Next() {
		var id int64
		var name string
		var appliedAt time.Time
		if err = rows.Scan(&id, &name, &appliedAt); err != nil {
			return nil, fmt.Errorf("failed to scan migration result: %w", err)
		}
		migrations = append(migrations, MigrationResult{
			ID:        id,
			Name:      name,
			AppliedAt: appliedAt,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("an error occurred during iteration through sql rows: %w", err)
	}

	return migrations, nil
}
