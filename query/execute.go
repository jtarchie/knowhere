package query

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

type Result struct {
	ID      int64   `db:"id"`
	MaxLat  float64 `db:"maxLat"`
	MaxLon  float64 `db:"maxLon"`
	MinLat  float64 `db:"minLat"`
	MinLon  float64 `db:"minLon"`
	Name    string  `db:"name"`
	OsmID   int64   `db:"osm_id"`
	OsmType int64   `db:"osm_type"`
}

func Execute(client *sql.DB, search string, fun func(string) (string, error)) ([]Result, error) {
	sqlQuery, err := fun(search)
	if err != nil {
		return nil, fmt.Errorf("could not parse the query: %w", err)
	}

	var results []Result

	err = sqlscan.Select(
		context.TODO(),
		client,
		&results,
		fmt.Sprintf("SELECT id, osm_id, osm_type, minLon, minLat, maxLon, maxLat, IFNULL(tags->>'$.name', 'Unknown') as name FROM (%s)", sqlQuery),
	)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	return results, nil
}
