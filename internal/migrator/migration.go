package migrator

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrMigrationNotPrepared = errors.New("migration is not prepared")
	ErrNilMigrationExecutor = errors.New("migration executors is nil")
	ErrNilMigrationPreparer = errors.New("migration preparer is nil")
	ErrNilExecutorContainer = errors.New("migration preparer returned nil Executors")
)

type Migrations []Migration

type Migration struct {
	id   int64
	name string

	isPrepared bool
	preparer   Preparer

	executors Executors
}

func NewMigration(id int64, name string, executor ExecutorTx) Migration {
	return Migration{
		id:         id,
		name:       name,
		isPrepared: true,
		executors: Executors{
			useDB:      false,
			executorTx: executor,
		},
	}
}

func NewMigrationNoTx(id int64, name string, executorDB ExecutorDB) Migration {
	return Migration{
		id:         id,
		name:       name,
		isPrepared: true,
		executors: Executors{
			useDB:      true,
			executorDB: executorDB,
		},
	}
}

func NewMigrationWithPreparer(id int64, name string, preparer Preparer) Migration {
	return Migration{
		id:         id,
		name:       name,
		isPrepared: false,
		preparer:   preparer,
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
