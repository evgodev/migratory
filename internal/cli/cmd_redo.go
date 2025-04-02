package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var redoCmd = &cobra.Command{
	Use:   "redo [--dir <path>] [-d <db-string>] [-s <schema>] [-t <table>]",
	Short: "Rollbacks and applies again last migration",
	Long: `The "redo" command rolls back the last applied migration, then applies it again.
Command creates migrations table if not exists.`,
	Example: `migratory redo -c /etc/config.yml
migratory redo -d postgresql://role:password@127.0.0.1:5432/database --dir example/migrations/
migratory redo -d postgresql://role:password@127.0.0.1:5432/database --dir migrations/ -t my_migrations_table`,
	Run: func(_ *cobra.Command, _ []string) {
		if err := rollback(config.Dir, config.Schema, config.Table, true); err != nil {
			fmt.Printf("unable to redo migration: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("success: last migration reapplied")
	},
}

func init() {
	rootCmd.AddCommand(redoCmd)
}
