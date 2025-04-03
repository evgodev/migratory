package main

import (
	"github.com/korfairo/migratory/internal/cli"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/lib/pq"
)

func main() {
	cli.Execute()
}
