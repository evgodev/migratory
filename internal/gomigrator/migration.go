package gomigrator

import (
	"context"
	"database/sql"
	"errors"
)

type Migrations []Migration

type Migration struct {
	id   int64
	name string

	isPrepared bool
	preparer   Preparer

	executor ExecutorContainer
}

type ExecutorContainer struct {
	useDB      bool
	executorTx ExecutorTx
	executorDB ExecutorDB
}

type ExecutorTx interface {
	UpTx(ctx context.Context, tx *sql.Tx) error
	DownTx(ctx context.Context, tx *sql.Tx) error
}

type ExecutorDB interface {
	Up(ctx context.Context, db *sql.DB) error
	Down(ctx context.Context, db *sql.DB) error
}

type Preparer interface {
	Prepare() (*ExecutorContainer, error)
}

func NewExecutorTxContainer(executorTx ExecutorTx) *ExecutorContainer {
	return &ExecutorContainer{
		useDB:      false,
		executorTx: executorTx,
	}
}

func NewExecutorDBContainer(executorDB ExecutorDB) *ExecutorContainer {
	return &ExecutorContainer{
		useDB:      true,
		executorDB: executorDB,
	}
}

func (e ExecutorContainer) UseDB() bool {
	return e.useDB
}

func (e ExecutorContainer) ExecutorTx() ExecutorTx {
	return e.executorTx
}

func (e ExecutorContainer) ExecutorDB() ExecutorDB {
	return e.executorDB
}

func NewMigration(id int64, name string, executor ExecutorTx) Migration {
	return Migration{
		id:         id,
		name:       name,
		isPrepared: true,
		executor: ExecutorContainer{
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
		executor: ExecutorContainer{
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

var (
	ErrMigrationNotPrepared = errors.New("migration is not prepared")
	ErrNilMigrationExecutor = errors.New("migration executor is nil")
	ErrNilMigrationPreparer = errors.New("migration preparer is nil")
	ErrNilExecutorContainer = errors.New("migration preparer returned nil ExecutorContainer")
)

func (m *Migration) UpTx(ctx context.Context, tx *sql.Tx) error {
	if !m.isPrepared {
		return ErrMigrationNotPrepared
	}

	if m.executor.ExecutorTx() == nil {
		return ErrNilMigrationExecutor
	}

	return m.executor.ExecutorTx().UpTx(ctx, tx)
}

func (m *Migration) DownTx(ctx context.Context, tx *sql.Tx) error {
	if !m.isPrepared {
		return ErrMigrationNotPrepared
	}

	if m.executor.ExecutorTx() == nil {
		return ErrNilMigrationExecutor
	}

	return m.executor.ExecutorTx().DownTx(ctx, tx)
}

func (m *Migration) UpDB(ctx context.Context, db *sql.DB) error {
	if !m.isPrepared {
		return ErrMigrationNotPrepared
	}

	if m.executor.ExecutorDB() == nil {
		return ErrNilMigrationExecutor
	}

	return m.executor.ExecutorDB().Up(ctx, db)
}

func (m *Migration) DownDB(ctx context.Context, db *sql.DB) error {
	if !m.isPrepared {
		return ErrMigrationNotPrepared
	}

	if m.executor.ExecutorDB() == nil {
		return ErrNilMigrationExecutor
	}

	return m.executor.ExecutorDB().Down(ctx, db)
}

func (m *Migration) ChooseExecutor() (noTx bool, err error) {
	if err = m.ensureIsPrepared(); err != nil {
		return false, err
	}

	return m.executor.UseDB(), nil
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

	m.executor = *executorController
	m.isPrepared = true

	return nil
}
