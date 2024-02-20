package query

import (
	"database/sql"
	"fmt"

	// "github.com/wroge/scan".
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func Execute(client *sql.DB, search string) ([]*geojson.Feature, error) {
	sql, err := ToSQL(search)
	if err != nil {
		return nil, fmt.Errorf("could not parse the query: %w", err)
	}

	rows, err := client.Query(fmt.Sprintf("SELECT id, minLon, minLat FROM (%s)", sql))
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	features := []*geojson.Feature{}

	for rows.Next() {
		var (
			feature geojson.Feature
			point   orb.Point
		)

		err := rows.Scan(&feature.ID, &point[0], &point[1])
		if err != nil {
			return nil, fmt.Errorf("could not load results: %w", err)
		}
		defer rows.Close()

		feature.Geometry = point
		feature.Type = "Feature"

		features = append(features, &feature)
	}

	return features, nil
}
