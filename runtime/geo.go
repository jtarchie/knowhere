package runtime

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/dop251/goja"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/jtarchie/knowhere/query"
	"github.com/paulmach/orb"
	"github.com/samber/lo"
)

type Geo struct {
	vm      *goja.Runtime
	timeout time.Duration
	client  *sql.DB
}

func (g *Geo) Rtree() *RTree {
	return &RTree{}
}

func (g *Geo) Query(qs string) Results {
	ctx, cancel := context.WithTimeout(context.TODO(), g.timeout)
	defer cancel()

	results, err := query.Execute(ctx, g.client, qs, query.ToIndexedSQL)
	if err != nil {
		slog.Error("geo.query", "query", qs, "err", err.Error())
		g.vm.Interrupt(fmt.Sprintf("could not execute query: %q", qs))

		return nil
	}

	return lo.Map(results, func(result query.Result, _ int) Result {
		return Result{result}
	})
}

type Prefix struct {
	Name     string  `db:"name"`
	FullName string  `db:"full_name"`
	MinLat   float64 `db:"minLat"`
	MaxLat   float64 `db:"maxLat"`
	MinLon   float64 `db:"minLon"`
	MaxLon   float64 `db:"maxLon"`
}

func (g *Geo) Prefixes() []Prefix {
	ctx, cancel := context.WithTimeout(context.TODO(), g.timeout)
	defer cancel()

	var results []Prefix

	err := sqlscan.Select(
		ctx,
		g.client,
		&results,
		`SELECT
				name, full_name, minLat, maxLat, minLon, maxLon
			FROM
				prefixes`,
	)
	if err != nil {
		slog.Error("geo.prefixes", "err", err.Error())
		g.vm.Interrupt("could not read prefixes")
	}

	return results
}

func (g *Geo) AsResults(results ...Result) Results {
	return Results(results)
}

func (g *Geo) AsBounds(bounds ...Bound) Bounds {
	return Bounds(bounds)
}

func (g *Geo) AsPoint(lat, lng float64) Point {
	return Point{lng, lat}
}

func (g *Geo) Color(index int) string {
	colorPalette := []string{
		"#E69F00", // Orange
		"#56B4E9", // Sky Blue
		"#009E73", // Bluish Green
		"#F0E442", // Yellow
		"#0072B2", // Blue
		"#D55E00", // Vermillion
		"#CC79A7", // Reddish Purple
		"#8DD3C7", // Light Blue-Green
		"#FDB462", // Soft Orange
		"#B3DE69", // Light Green
		"#FFED6F", // Light Yellow
		"#6A3D9A", // Deep Purple
		"#B15928", // Brownish-Orange
		"#44AA99", // Teal
		"#117733", // Dark Green
		"#999933", // Olive Green
		"#AA4499", // Purple
		"#DDCC77", // Light Tan
		"#882255", // Dark Red
		"#332288", // Dark Blue
	}

	return colorPalette[index%len(colorPalette)]
}
