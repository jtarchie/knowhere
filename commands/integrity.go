package commands

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"strings"

	_ "github.com/jtarchie/sqlitezstd"
)

type Integrity struct {
	UncompressedDB string `help:"path to the uncompressed SQLite database" required:""`
	CompressedDB   string `help:"path to the compressed SQLite database" required:""`
}

func (i *Integrity) Run(stdout io.Writer) error {
	uncompressedDB, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_query_only=true&immutable=true&mode=ro", i.UncompressedDB))
	if err != nil {
		return fmt.Errorf("could not open uncompressed database: %w", err)
	}
	defer uncompressedDB.Close()

	compressedDB, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_query_only=true&immutable=true&mode=ro&vfs=zstd", i.CompressedDB))
	if err != nil {
		return fmt.Errorf("could not open compressed database: %w", err)
	}
	defer compressedDB.Close()

	tables, err := getTables(uncompressedDB)
	if err != nil {
		return fmt.Errorf("could not get tables: %w", err)
	}

	for _, table := range tables {
		uncompressedCount, err := getTableCount(uncompressedDB, table)
		if err != nil {
			return fmt.Errorf("could not get count for uncompressed table %s: %w", table, err)
		}

		compressedCount, err := getTableCount(compressedDB, table)
		if err != nil {
			return fmt.Errorf("could not get count for compressed table %s: %w", table, err)
		}

		if uncompressedCount != compressedCount {
			slog.Error("integrity check failed",
				"table", table,
				"uncompressed_count", uncompressedCount,
				"compressed_count", compressedCount)
		} else {
			slog.Info("integrity check passed",
				"table", table,
				"count", uncompressedCount)
		}
	}

	fmt.Fprintln(stdout, "Integrity check completed")
	return nil
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		if !strings.HasPrefix(tableName, "sqlite_") {
			tables = append(tables, tableName)
		}
	}
	return tables, rows.Err()
}

func getTableCount(db *sql.DB, table string) (int, error) {
	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
	return count, err
}
