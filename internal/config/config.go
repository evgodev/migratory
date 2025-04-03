package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/korfairo/migratory/internal/migrator"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Dir   string `yaml:"directory"`
	DSN   string `yaml:"dsn"`
	Table string `yaml:"table"`

	Dialect string
}

var defaultConfig = Config{
	Dir:     ".",
	DSN:     "",
	Table:   "migrations",
	Dialect: "",
}

var (
	ErrReadConfigFile   = errors.New("failed to read config file")
	ErrUnmarshalFailure = errors.New("failed to unmarshal config")
)

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadConfigFile, err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnmarshalFailure, err)
	}

	cfg.Dialect = dialectFromDSN(cfg.DSN)

	setDefaultValues(cfg)

	return cfg, nil
}

func dialectFromDSN(dsn string) string {
	if len(dsn) == 0 {
		return ""
	}
	firstWord := dsn[:strings.Index(dsn, ":")]
	switch firstWord {
	case migrator.DialectPostgres, migrator.DialectClickHouse:
		return firstWord
	default:
		return ""
	}
}

func setDefaultValues(cfg *Config) {
	if cfg.Dir == "" {
		cfg.Dir = defaultConfig.Dir
	}

	if cfg.Table == "" {
		cfg.Table = defaultConfig.Table
	}
}
