package runtime

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/dop251/goja"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/iancoleman/strcase"
	"github.com/jtarchie/knowhere/address"
	"github.com/jtarchie/knowhere/query"
	"github.com/samber/lo"
)

type Query struct {
	vm      *goja.Runtime
	timeout time.Duration
	client  *sql.DB
}

func (g *Query) Execute(qs string) Results {
	slog.Debug("query.execute", "query", qs)

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

func (q *Query) Union(queries ...string) Results {
	slog.Debug("query.execute", "query", queries)

	ctx, cancel := context.WithTimeout(context.TODO(), q.timeout)
	defer cancel()

	results, err := query.Union(ctx, q.client, queries...)
	if err != nil {
		slog.Error("query.union", "query", queries, "err", err.Error())
		q.vm.Interrupt("could not union queries")

		return nil
	}

	return lo.Map(results, func(result query.Result, _ int) Result {
		return Result{result}
	})
}

func (g *Query) FromAddress(fullAddress string, area string) Results {
	parts, ok := address.Parse(fullAddress, true)
	if !ok {
		return Results{}
	}

	if area == "" {
		area = strcase.ToSnake(parts["state"])
	}

	return g.Union(
		`nwr[addr:housenumber=~"`+parts["house_number"]+`"][addr:street=~"`+parts["road"]+`*"][addr:city=~"`+parts["city"]+`"](area="`+area+`")`,
		`nwr[addr:street=~"`+parts["road"]+`*"][addr:city=~"`+parts["city"]+`"](area="`+area+`")`,
		`nwr[addr:housenumber=~"`+parts["house_number"]+`"][addr:street=~"`+parts["road"]+`*"](area="`+area+`")`,
		`nwr[name=~"`+parts["road"]+`*"][highway=residential](area="`+area+`")`,
	)
}

type Prefix struct {
	Name     string  `db:"name"`
	FullName string  `db:"full_name"`
	MinLat   float64 `db:"minLat"`
	MaxLat   float64 `db:"maxLat"`
	MinLon   float64 `db:"minLon"`
	MaxLon   float64 `db:"maxLon"`
}

func (g *Query) Areas() []Prefix {
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
				areas`,
	)
	if err != nil {
		slog.Error("query.areas", "err", err.Error())
		g.vm.Interrupt("could not read areas")
	}

	return results
}
