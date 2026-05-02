package cli

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database tools",
}

var dbPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the database path",
	RunE:  runDBPath,
}

var dbQueryCmd = &cobra.Command{
	Use:   "query [sql]",
	Short: "Execute a SQL query",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  runDBQuery,
}

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	RunE:  runDBMigrate,
}

var dbFormat string

func init() {
	dbCmd.AddCommand(dbPathCmd)
	dbCmd.AddCommand(dbQueryCmd)
	dbCmd.AddCommand(dbMigrateCmd)

	dbQueryCmd.Flags().StringVar(&dbFormat, "format", "tsv", "Output format (json, tsv)")
	rootCmd.AddCommand(dbCmd)
}

func getDBPath() string {
	home, _ := os.UserHomeDir()
	return home + "/.local/share/freecode/freecode.db"
}

func runDBPath(cmd *cobra.Command, args []string) error {
	path := getDBPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Database does not exist yet")
		fmt.Printf("Path: %s\n", path)
		return nil
	}
	fmt.Println(path)
	return nil
}

func runDBQuery(cmd *cobra.Command, args []string) error {
	path := getDBPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Database does not exist")
		return nil
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if len(args) == 0 {
		fmt.Println("Interactive mode not supported. Provide a SQL query.")
		fmt.Printf("Example: freecode db query \"SELECT * FROM sessions LIMIT 5\"\n")
		return nil
	}

	query := args[0]
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	if dbFormat == "json" {
		fmt.Println("[")
		first := true
		for rows.Next() {
			values := make([]interface{}, len(cols))
			valuePtrs := make([]interface{}, len(cols))
			for i := range values {
				valuePtrs[i] = &values[i]
			}
			if err := rows.Scan(valuePtrs...); err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			if first {
				first = false
			} else {
				fmt.Println(",")
			}
			fmt.Print("{")
			for i, col := range cols {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("\"%s\": ", col)
				switch v := values[i].(type) {
				case []byte:
					fmt.Printf("\"%s\"", string(v))
				case nil:
					fmt.Print("null")
				default:
					fmt.Printf("\"%v\"", v)
				}
			}
			fmt.Print("}")
		}
		fmt.Println("\n]")
	} else {
		fmt.Println("TSV output:")
		for rows.Next() {
			values := make([]interface{}, len(cols))
			valuePtrs := make([]interface{}, len(cols))
			for i := range values {
				valuePtrs[i] = &values[i]
			}
			if err := rows.Scan(valuePtrs...); err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}
			for i, v := range values {
				if i > 0 {
					fmt.Print("\t")
				}
				switch val := v.(type) {
				case []byte:
					fmt.Print(string(val))
				case nil:
					fmt.Print("NULL")
				default:
					fmt.Print(val)
				}
			}
			fmt.Println()
		}
	}

	return nil
}

func runDBMigrate(cmd *cobra.Command, args []string) error {
	path := getDBPath()

	dir := os.Getenv("FREECODE_HOME")
	if dir == "" {
		home, _ := os.UserHomeDir()
		dir = home + "/.local/share/freecode"
	}

	migrationsDir := dir + "/migrations"

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		fmt.Println("No migrations directory found")
		return nil
	}

	fmt.Printf("Database: %s\n", path)
	fmt.Printf("Migrations: %s\n", migrationsDir)
	fmt.Println("\nMigration system not yet implemented")
	return nil
}