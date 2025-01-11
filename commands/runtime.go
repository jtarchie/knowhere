package commands

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/jtarchie/knowhere/services"
	_ "github.com/jtarchie/sqlitezstd"
)

type Runtime struct {
	Filename       *os.File      `help:"script to execute" required:""`
	DB             string        `help:"sqlite file to server"      required:""`
	RuntimeTimeout time.Duration `help:"the timeout for single runtime" default:"2s"`
}

func (r *Runtime) Run(stdout io.Writer) error {
	client, err := sql.Open("sqlite3", connectionString(r.DB))
	if err != nil {
		return fmt.Errorf("could not open database file: %w", err)
	}

	runtime := services.NewRuntime(client, r.RuntimeTimeout)

	contents, err := io.ReadAll(r.Filename)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	slog.Info("script.start")

	startTime := time.Now()

	value, err := runtime.Execute(string(contents))
	if err != nil {
		return fmt.Errorf("could not execute script: %w", err)
	}

	contents, err = json.MarshalIndent(value.Export(), "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal value as JSON: %w", err)
	}

	fmt.Fprintln(stdout, string(contents))

	slog.Info("script.end", "elapsed", time.Since(startTime).String())

	return nil
}
