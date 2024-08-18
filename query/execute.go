package query

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/goccy/go-json"
	"log/slog"
	"strings"

	"github.com/georgysavva/scany/v2/sqlscan"
)

type Tags map[string]string

func (t *Tags) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), t)
}

type Result struct {
	ID      int64      `db:"id"       js:"id"`
	MaxLat  float64    `db:"maxLat"   js:"max_lat"`
	MaxLon  float64    `db:"maxLon"   js:"max_lon"`
	MinLat  float64    `db:"minLat"   js:"min_lat"`
	MinLon  float64    `db:"minLon"   js:"min_lon"`
	OsmID   int64      `db:"osm_id"   js:"osm_id"`
	OsmType FilterType `db:"osm_type" js:"osm_type"`
	Tags    Tags       `db:"tags"     js:"tags"`
}

func (r Result) Name() string {
	return r.Tags["name"]
}

func Execute(ctx context.Context, client *sql.DB, search string, fun func(string) (string, error)) ([]Result, error) {
	sqlQuery, err := fun(search)
	if err != nil {
		return nil, fmt.Errorf("could not parse the query: %w", err)
	}

	slog.Debug("query.execute", "sql", sqlQuery)

	var results []Result

	err = sqlscan.Select(
		ctx,
		client,
		&results,
		fmt.Sprintf(`
			SELECT
				id,
				osm_id,
				osm_type,
				minLon,
				minLat,
				maxLon,
				maxLat,
				json(tags) as tags
			FROM (%s)`, sqlQuery),
	)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	return results, nil
}

func Union(ctx context.Context, client *sql.DB, queries ...string) ([]Result, error) {
	sqlQueries := make([]string, 0, len(queries))

	for _, query := range queries {
		sqlQuery, err := ToIndexedSQL(query)
		if err != nil {
			return nil, fmt.Errorf("could not parse the query: %w", err)
		}

		sqlQueries = append(sqlQueries, fmt.Sprintf(`
		SELECT
			id,
			osm_id,
			osm_type,
			minLon,
			minLat,
			maxLon,
			maxLat,
			json_object(
				'name', tags->>'$.name'
			) as tags
		FROM (%s)`, sqlQuery))
	}

	unionSQL := strings.Join(sqlQueries, "\n UNION ALL \n")
	slog.Debug("query.union", "sql", unionSQL)

	var results []Result

	err := sqlscan.Select(
		ctx,
		client,
		&results,
		unionSQL,
	)
	if err != nil {
		return nil, fmt.Errorf("could not union query: %w", err)
	}

	return results, nil
}
