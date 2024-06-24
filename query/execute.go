package query

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

type Result struct {
	ID      int64      `db:"id"       js:"id"`
	MaxLat  float64    `db:"maxLat"   js:"max_lat"`
	MaxLon  float64    `db:"maxLon"   js:"max_lon"`
	MinLat  float64    `db:"minLat"   js:"min_lat"`
	MinLon  float64    `db:"minLon"   js:"min_lon"`
	Name    string     `db:"name"     js:"name"`
	OsmID   int64      `db:"osm_id"   js:"osm_id"`
	OsmType FilterType `db:"osm_type" js:"osm_type"`
}

func Execute(ctx context.Context, client *sql.DB, search string, fun func(string) (string, error)) ([]Result, error) {
	sqlQuery, err := fun(search)
	if err != nil {
		return nil, fmt.Errorf("could not parse the query: %w", err)
	}

	var results []Result

	err = sqlscan.Select(
		ctx,
		client,
		&results,
		fmt.Sprintf("SELECT id, osm_id, osm_type, minLon, minLat, maxLon, maxLat, IFNULL(tags->>'$.name', 'Unknown') as name FROM (%s)", sqlQuery),
	)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	return results, nil
}
