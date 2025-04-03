package store

import "fmt"

const schemaName = "public"

type postgresQueryBuilder struct{}

func (p *postgresQueryBuilder) MigrationsTableExists(tableName string) string {
	q := `SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = '%s' AND tablename  = '%s')`
	return fmt.Sprintf(q, schemaName, tableName)
}

func (p *postgresQueryBuilder) CreateMigrationsTable(tableName string) string {
	q := `CREATE TABLE %s.%s (
		id bigint PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		applied_at timestamp NOT NULL
	)`
	return fmt.Sprintf(q, schemaName, tableName)
}

func (p *postgresQueryBuilder) InsertMigration(tableName string) string {
	q := `INSERT INTO %s.%s (id, name, applied_at) VALUES ($1, $2, now())`
	return fmt.Sprintf(q, schemaName, tableName)
}

func (p *postgresQueryBuilder) DeleteMigration(tableName string) string {
	q := `DELETE FROM %s.%s WHERE id = $1`
	return fmt.Sprintf(q, schemaName, tableName)
}

func (p *postgresQueryBuilder) ListMigrations(tableName string) string {
	q := `SELECT id, name, applied_at FROM %s.%s ORDER BY id ASC`
	return fmt.Sprintf(q, schemaName, tableName)
}

func (p *postgresQueryBuilder) SelectLastMigrationID(tableName string) string {
	q := `SELECT id FROM %s.%s ORDER BY id DESC LIMIT 1`
	return fmt.Sprintf(q, schemaName, tableName)
}
