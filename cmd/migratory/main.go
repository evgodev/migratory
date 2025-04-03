package main

import (
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/korfairo/migratory/internal/cli"
	_ "github.com/lib/pq"
)

func main() {
	cli.Execute()
}
