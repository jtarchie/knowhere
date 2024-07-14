package commands

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/jtarchie/knowhere/query"
)

type Generate struct {
	Value string `arg:""`
}

func (q *Generate) Run(stdout io.Writer) error {
	sqlQuery, err := query.ToIndexedSQL(q.Value)
	if err != nil {
		return fmt.Errorf("could not parse query: %w", err)
	}

	slog.Debug("generate.query", "query", q.Value, "sql", sqlQuery)
	fmt.Fprintln(stdout, sqlQuery) //nolint

	return nil
}
