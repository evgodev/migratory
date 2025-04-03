package migrator

import (
	"fmt"
	"os"

	"github.com/korfairo/migratory/internal/migrator/executors"
	"github.com/korfairo/migratory/internal/migrator/parser"
)

type Preparer interface {
	Prepare() (*Executors, error)
}

type sqlPreparer struct {
	sourcePath string
}

func newSQLPreparer(sourceFilePath string) sqlPreparer {
	return sqlPreparer{
		sourcePath: sourceFilePath,
	}
}

func (s sqlPreparer) Prepare() (*Executors, error) {
	file, err := os.Open(s.sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file at path %s: %w", s.sourcePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	parsed, err := parser.ParseMigration(file)
	if parsed == nil || err != nil {
		return nil, fmt.Errorf("failed to parse migration %s: %w", s.sourcePath, err)
	}

	var container *Executors
	if parsed.DisableTransactionUp || parsed.DisableTransactionDown {
		executor := executors.NewSQLExecutorNoTx(parsed.UpStatements, parsed.DownStatements)
		container = newExecutorDBContainer(executor)
	} else {
		executor := executors.NewSQLExecutor(parsed.UpStatements, parsed.DownStatements)
		container = newExecutorTxContainer(executor)
	}

	return container, nil
}
