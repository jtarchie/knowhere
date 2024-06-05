package runtime

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/dop251/goja"
	"github.com/jtarchie/knowhere/query"
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
		slog.Error("execute.failed", "query", qs, "err", err.Error())
		g.vm.Interrupt(fmt.Sprintf("could not execute query: %q", qs))
	}

	return lo.Map(results, func(result query.Result, _ int) Result {
		return Result{result}
	})
}

func (g *Geo) AsResults(results ...Result) Results {
	return Results(results)
}

func (g *Geo) AsBounds(bounds ...Bound) Bounds {
	return Bounds(bounds)
}
