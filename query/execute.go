package query

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

type Result struct {
	ID     int64   `db:"id"`
	MinLon float64 `db:"minLon"`
	MaxLon float64 `db:"maxLon"`
	MinLat float64 `db:"minLat"`
	MaxLat float64 `db:"maxLat"`
	Name   string  `db:"name"`
}

func Execute(client *sql.DB, search string) ([]Result, error) {
	sqlQuery, err := ToSQL(search)
	if err != nil {
		return nil, fmt.Errorf("could not parse the query: %w", err)
	}

	var results []Result

	err = sqlscan.Select(
		context.TODO(),
		client,
		&results,
		fmt.Sprintf("SELECT id, minLon, minLat, maxLon, maxLat, IFNULL(tags->>'$.name', 'Unknown') as name FROM (%s)", sqlQuery),
	)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	return results, nil
}
