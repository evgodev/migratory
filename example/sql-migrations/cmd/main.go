package main

import (
	"database/sql"
	"log"

	"github.com/evgodev/migratory"
	_ "github.com/lib/pq"
)

func main() {
	// Connect to your database
	db, err := sql.Open(
		"postgres",
		"postgres://postgres:password@localhost:5432/test?sslmode=disable",
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Set the SQL migrations directory
	migratory.SetSQLDirectory("./migrations")
	migratory.SetDialect(migratory.Postgres)

	// Apply all migrations using SQL files from the directory
	count, err := migratory.Up(db)
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Printf("Applied %d migrations\n", count)

	// Get the current migration status
	status, err := migratory.GetStatus(db)
	if err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}

	log.Println("Migration Status:")
	for _, m := range status {
		log.Printf("- [%d] %s: %t %s\n", m.ID, m.Name, m.IsApplied, m.AppliedAt)
	}
}
