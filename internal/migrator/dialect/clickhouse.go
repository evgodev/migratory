package dialect

import (
	"fmt"
)

type Clickhouse struct{}

func (c *Clickhouse) MigrationsTableExists(tableName string) string {
	q := `EXISTS ` + tableName
	return q
}

func (c *Clickhouse) CreateMigrationsTable(tableName string) string {
	q := `CREATE TABLE %s (
		id Int64 PRIMARY KEY,
		name String NOT NULL,
		applied_at timestamp NOT NULL
	)
	ENGINE = MergeTree() PRIMARY KEY id;`
	return fmt.Sprintf(q, tableName)
}

func (c *Clickhouse) InsertMigration(tableName string) string {
	q := `INSERT INTO %s (id, name, applied_at) VALUES (?, ?, now())`
	return fmt.Sprintf(q, tableName)
}

func (c *Clickhouse) DeleteMigration(tableName string) string {
	q := `ALTER TABLE %s DELETE WHERE id = ? SETTINGS mutations_sync = 2;`
	return fmt.Sprintf(q, tableName)
}

func (c *Clickhouse) ListMigrations(tableName string) string {
	q := `SELECT id, name, applied_at FROM %s ORDER BY id ASC`
	return fmt.Sprintf(q, tableName)
}

func (c *Clickhouse) SelectLastMigrationID(tableName string) string {
	q := `SELECT id FROM %s ORDER BY id DESC LIMIT 1`
	return fmt.Sprintf(q, tableName)
}
