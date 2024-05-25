package commands

import (
	"fmt"

	"github.com/jtarchie/knowhere/query"
)

type Query struct {
	Value string `arg:""`
}

func (q *Query) Run() error {
	sqlQuery, err := query.ToIndexedSQL(q.Value)
	if err != nil {
		return fmt.Errorf("could not parse query: %w", err)
	}

	fmt.Println(sqlQuery) //nolint

	return nil
}
