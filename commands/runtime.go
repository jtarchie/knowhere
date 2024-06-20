package commands

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jtarchie/knowhere/services"
	"github.com/jtarchie/sqlitezstd"
)

type Runtime struct {
	Filename       *os.File      `help:"script to execute" required:""`
	DB             string        `help:"sqlite file to server"      required:""`
	RuntimeTimeout time.Duration `help:"the timeout for single runtime" default:"2s"`
}

func (r *Runtime) Run(stdout io.Writer) error {
	connectionString := fmt.Sprintf("file:%s?_query_only=true&immutable=true&mode=ro", r.DB)

	if strings.Contains(r.DB, ".zst") {
		err := sqlitezstd.Init()
		if err != nil {
			return fmt.Errorf("could not load sqlite zstd vfs: %w", err)
		}

		connectionString += "&vfs=zstd"
	}

	client, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return fmt.Errorf("could not open database file: %w", err)
	}

	runtime := services.NewRuntime(client, r.RuntimeTimeout)

	contents, err := io.ReadAll(r.Filename)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	value, err := runtime.Execute(string(contents))
	if err != nil {
		return fmt.Errorf("could not execute script: %w", err)
	}

	contents, err = json.MarshalIndent(value.Export(), "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal value as JSON: %w", err)
	}

	fmt.Fprintln(stdout, string(contents))

	return nil
}
