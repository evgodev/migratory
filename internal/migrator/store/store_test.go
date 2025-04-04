package store

import (
	"testing"

	"github.com/evgodev/migratory/internal/require"
)

func TestNewStore(t *testing.T) {
	type args struct {
		dbDialect string
		tableName string
	}
	tests := map[string]struct {
		args    args
		want    *Store
		wantErr error
	}{
		"existing dialect Posgres": {
			args: args{
				dbDialect: Postgres,
				tableName: "migrations",
			}, want: &Store{
				"migrations",
				&postgresQueryBuilder{},
			},
			wantErr: nil,
		},
		"existing dialect Clickhouse": {
			args: args{
				dbDialect: ClickHouse,
				tableName: "migrations",
			}, want: &Store{
				"migrations",
				&clickhouseQueryBuilder{},
			},
			wantErr: nil,
		},
		"unknown dialect": {
			args: args{
				dbDialect: "mysql",
				tableName: "migrations",
			}, want: nil,
			wantErr: ErrUnsupportedDialect,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := New(test.args.dbDialect, test.args.tableName)
			require.ErrorIs(t, err, test.wantErr, "New(...) error")
			require.Equal(t, got, test.want, "New(...) new Store")
		})
	}
}
