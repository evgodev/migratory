package migrator

import (
	"context"
	"database/sql"
	"errors"

	"github.com/evgodev/migratory/internal/migrator/executor"
)

var (
	ErrMigrationNotPrepared = errors.New("migration is not prepared")
	ErrNilMigrationExecutor = errors.New("migration executors is nil")
	ErrNilMigrationPreparer = errors.New("migration preparer is nil")
	ErrNilExecutorContainer = errors.New("migration preparer returned nil executors")
)

type Migrations []Migration

// Migration represents a database migration with a unique ID, name, and executors for transactional
// or non-transactional use. This type manages whether a migration is prepared for execution and
// supports lazy loading (SQL migrations are parsed only during the migration application process).
type Migration struct {
	id   int64
	name string

	isPrepared bool
	preparer   *sqlPreparer

	executors executors
}

func NewGoMigration(id int64, name string, up, down executor.GoMigrateFn) Migration {
	return Migration{
		id:         id,
		name:       name,
		isPrepared: true,
		executors: executors{
			useDB:      false,
			executorTx: executor.NewGoExecutor(up, down),
		},
	}
}

func NewGoMigrationNoTx(id int64, name string, up, down executor.GoMigrateNoTxFn) Migration {
	return Migration{
		id:         id,
		name:       name,
		isPrepared: true,
		executors: executors{
			useDB:      true,
			executorDB: executor.NewGoExecutorNoTx(up, down),
		},
	}
}

func NewSQLMigration(id int64, name, filePath string) Migration {
	return Migration{
		id:         id,
		name:       name,
		isPrepared: false,
		preparer:   newSQLPreparer(filePath),
	}
}

func (m *Migration) UpTx(ctx context.Context, tx *sql.Tx) error {
	if !m.isPrepared {
		return ErrMigrationNotPrepared
	}

	if m.executors.ExecutorTx() == nil {
		return ErrNilMigrationExecutor
	}

	return m.executors.ExecutorTx().UpTx(ctx, tx)
}

func (m *Migration) DownTx(ctx context.Context, tx *sql.Tx) error {
	if !m.isPrepared {
		return ErrMigrationNotPrepared
	}

	if m.executors.ExecutorTx() == nil {
		return ErrNilMigrationExecutor
	}

	return m.executors.ExecutorTx().DownTx(ctx, tx)
}

func (m *Migration) UpDB(ctx context.Context, db *sql.DB) error {
	if !m.isPrepared {
		return ErrMigrationNotPrepared
	}

	if m.executors.ExecutorDB() == nil {
		return ErrNilMigrationExecutor
	}

	return m.executors.ExecutorDB().Up(ctx, db)
}

func (m *Migration) DownDB(ctx context.Context, db *sql.DB) error {
	if !m.isPrepared {
		return ErrMigrationNotPrepared
	}

	if m.executors.ExecutorDB() == nil {
		return ErrNilMigrationExecutor
	}

	return m.executors.ExecutorDB().Down(ctx, db)
}

func (m *Migration) ChooseExecutor() (noTx bool, err error) {
	if err = m.ensureIsPrepared(); err != nil {
		return false, err
	}

	return m.executors.useDB, nil
}

func (m *Migration) ID() int64 {
	return m.id
}

func (m *Migration) Name() string {
	return m.name
}

func (m *Migration) ensureIsPrepared() error {
	if m.isPrepared {
		return nil
	}

	if m.preparer == nil {
		return ErrNilMigrationPreparer
	}

	return m.prepare()
}

func (m *Migration) prepare() error {
	executorController, err := m.preparer.Prepare()
	if err != nil {
		return err
	}

	if executorController == nil {
		return ErrNilExecutorContainer
	}

	m.executors = *executorController
	m.isPrepared = true

	return nil
}
