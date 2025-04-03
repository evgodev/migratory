package migrator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	separator       = "_"
	fileNamePattern = "*.sql"
)

var (
	ErrNoMigrationFiles = errors.New("migration files *.sql not found")
	ErrDuplicatedID     = errors.New("duplicated migrations ID detected")
	ErrGlobMigrations   = errors.New("unable to search migrations in directory")
	ErrDirectoryCheck   = errors.New("unable to check directory existence")
	ErrNoSeparator      = errors.New("no separator found in file name")
	ErrParseID          = errors.New("unable to parse ID in file name")
)

// ParseMigrationFileName parses a given migration file name into its ID and name.
func ParseMigrationFileName(fileName string) (id int64, migrationName string, err error) {
	base := filepath.Base(fileName)
	nameWithoutExt := strings.TrimSuffix(base, filepath.Ext(fileName))

	separatorIdx := strings.Index(nameWithoutExt, separator)
	if separatorIdx < 0 {
		return 0, "", ErrNoSeparator
	}

	id, err = strconv.ParseInt(nameWithoutExt[:separatorIdx], 10, 64)
	if err != nil {
		return 0, "", ErrParseID
	}

	migrationName = nameWithoutExt[separatorIdx+1:]

	return id, migrationName, nil
}

// SeekMigrations identifies and parses migration files in the given directory using the provided file system interface.
// Returns a sorted list of migration objects or an error if the directory or files are invalid.
func SeekMigrations(dir string) (Migrations, error) {
	if _, err := os.Stat(dir); err != nil {
		return nil, errors.Join(ErrDirectoryCheck, err)
	}

	migrationFiles, err := findMigrationFiles(dir)
	if err != nil {
		return nil, err
	}

	migrations, err := parseMigrationFiles(migrationFiles)
	if err != nil {
		return nil, err
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID() < migrations[j].ID()
	})

	return migrations, nil
}

func findMigrationFiles(dir string) ([]string, error) {
	migrationFiles, err := filepath.Glob(filepath.Join(dir, fileNamePattern))
	if err != nil {
		return nil, errors.Join(ErrGlobMigrations, err)
	}

	if len(migrationFiles) == 0 {
		return nil, ErrNoMigrationFiles
	}

	return migrationFiles, nil
}

func parseMigrationFiles(filePaths []string) (Migrations, error) {
	var migrations Migrations
	uniqueIDMap := make(map[int64]struct{}, len(filePaths))

	for _, filePath := range filePaths {
		id, name, err := ParseMigrationFileName(filePath)
		if err != nil {
			return nil, fmt.Errorf("file %s doesn't match the migration pattern: %w", filePath, err)
		}

		if _, exists := uniqueIDMap[id]; exists {
			return nil, fmt.Errorf("migration ID %d is duplicated: %w", id, ErrDuplicatedID)
		}

		uniqueIDMap[id] = struct{}{}
		migrations = append(migrations,
			NewMigrationWithPreparer(id, name, newSQLPreparer(filePath)))
	}

	return migrations, nil
}
