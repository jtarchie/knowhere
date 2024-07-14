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
	"github.com/samber/lo"
)

type Query struct {
	vm      *goja.Runtime
	timeout time.Duration
	client  *sql.DB
}

func (g *Query) Execute(qs string) Results {
	ctx, cancel := context.WithTimeout(context.TODO(), g.timeout)
	defer cancel()

	results, err := query.Execute(ctx, g.client, qs, query.ToIndexedSQL)
	if err != nil {
		slog.Error("query.execute", "query", qs, "err", err.Error())
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

func (g *Query) Prefixes() []Prefix {
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
		slog.Error("query.prefixes", "err", err.Error())
		g.vm.Interrupt("could not read prefixes")
	}

	return results
}
