package migrator

import (
	"fmt"
	"os"

	sqlexec "github.com/korfairo/migratory/internal/migrator/executors"
	"github.com/korfairo/migratory/internal/migrator/parser"
)

type sqlPreparer struct {
	filePath string
}

func newSQLPreparer(filePath string) *sqlPreparer {
	return &sqlPreparer{
		filePath: filePath,
	}
}

func (s sqlPreparer) Prepare() (*executors, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file at path %s: %w", s.filePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	parsed, err := parser.ParseMigration(file)
	if parsed == nil || err != nil {
		return nil, fmt.Errorf("failed to parse migration %s: %w", s.filePath, err)
	}

	var container *executors
	if parsed.DisableTransactionUp || parsed.DisableTransactionDown {
		executor := sqlexec.NewSQLExecutorNoTx(parsed.UpStatements, parsed.DownStatements)
		container = newExecutorDBContainer(executor)
	} else {
		executor := sqlexec.NewSQLExecutor(parsed.UpStatements, parsed.DownStatements)
		container = newExecutorTxContainer(executor)
	}

	return container, nil
}
