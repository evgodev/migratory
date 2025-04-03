package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/korfairo/migratory/internal/migrator"
	"github.com/korfairo/migratory/internal/require"
)

func TestReadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := map[string]struct {
		content string
		want    *Config
		err     error
	}{
		"valid postgres config": {
			content: `
directory: /path/to/directory
dsn: postgres://user:password@localhost:5432/my_db
table: migrations
`,
			want: &Config{
				Dir:     "/path/to/directory",
				DSN:     "postgres://user:password@localhost:5432/my_db",
				Table:   "migrations",
				Dialect: migrator.DialectPostgres,
			},
			err: nil,
		},
		"valid clickhouse config": {
			content: `
directory: /clickhouse/directory
dsn: clickhouse://user:password@localhost:8123/default
table: migrations
`,
			want: &Config{
				Dir:     "/clickhouse/directory",
				DSN:     "clickhouse://user:password@localhost:8123/default",
				Table:   "migrations",
				Dialect: migrator.DialectClickHouse,
			},
			err: nil,
		},
		"empty config": {
			content: "",
			want: &Config{
				Dir:     defaultConfig.Dir,
				DSN:     defaultConfig.DSN,
				Table:   defaultConfig.Table,
				Dialect: defaultConfig.Dialect,
			},
			err: nil,
		},
		"invalid config": {
			content: `
directory: [ /path/to/directory]
dsn = postgres://user:password@localhost:5432/mydb
table migrations
`,
			want: nil,
			err:  ErrUnmarshalFailure,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a temporary file with the test content
			filePath := filepath.Join(tmpDir, name+".yml")

			if test.content != "" {
				err := os.WriteFile(filePath, []byte(test.content), 0644)
				if err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
			} else {
				// Create an empty file
				f, err := os.Create(filePath)
				if err != nil {
					t.Fatalf("Failed to create empty test file: %v", err)
				}
				f.Close()
			}

			got, err := ReadConfig(filePath)

			require.Equal(t, got, test.want, "ReadConfig(...) config")
			require.ErrorIs(t, err, test.err, "ReadConfig(...)")
		})
	}

	// Test for nonexistent file separately
	t.Run("nonexistent file", func(t *testing.T) {
		nonexistentPath := filepath.Join(tmpDir, "nonexistent.yml")
		got, err := ReadConfig(nonexistentPath)

		require.Nil(t, got, "ReadConfig(...) config")
		require.ErrorIs(t, err, ErrReadConfigFile, "ReadConfig(...)")
	})
}
