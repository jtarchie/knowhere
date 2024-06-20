package commands

import (
	"fmt"
	"io"

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

	fmt.Fprintln(stdout, sqlQuery) //nolint

	return nil
}
