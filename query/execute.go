package query

import (
	"database/sql"
	"fmt"

	"github.com/wroge/scan"
)

type Result struct {
	ID      int64   `json:"id"`
	MinLat  float64 `json:"minLat"`
	MaxLat  float64 `json:"maxLat"`
	MinLon  float64 `json:"minLon"`
	MaxLon  float64 `json:"maxLon"`
	OsmType string  `json:"type"`
}

func Execute(client *sql.DB, search string) ([]Result, error) {
	sql, err := ToSQL(search)
	if err != nil {
		return nil, fmt.Errorf("could not parse the query: %w", err)
	}

	rows, err := client.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	columns := map[string]scan.Scanner[Result]{
		"id":       scan.Any(func(p *Result, id int64) { p.ID = id }),
		"minLat":   scan.Any(func(p *Result, minLat float64) { p.MinLat = minLat }),
		"maxLat":   scan.Any(func(p *Result, maxLat float64) { p.MaxLat = maxLat }),
		"minLon":   scan.Any(func(p *Result, minLon float64) { p.MinLon = minLon }),
		"maxLon":   scan.Any(func(p *Result, maxLon float64) { p.MaxLon = maxLon }),
		"osm_type": scan.Any(func(p *Result, osmType string) { p.OsmType = osmType }),
	}

	results, err := scan.All(rows, columns)
	if err != nil {
		return nil, fmt.Errorf("could not scan all rows: %w", err)
	}

	return results, nil
}
