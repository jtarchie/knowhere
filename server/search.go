package server

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/jtarchie/knowhere/query"
	r "github.com/jtarchie/knowhere/runtime"
	"github.com/labstack/echo/v4"
	"github.com/paulmach/orb/geojson"
	"github.com/samber/lo"
)

func locationSearch(client *sql.DB) func(echo.Context) error {
	return func(ctx echo.Context) error {
		search := ctx.FormValue("search")
		if search == "" {
			return response(ctx, http.StatusBadRequest, map[string]string{
				"error": "Required to have a `search` for the query",
			})
		}

		results, err := query.Execute(ctx.Request().Context(), client, search, query.ToIndexedSQL)
		if err != nil {
			slog.Error("search.error", slog.String("error", err.Error()))

			return response(ctx, http.StatusBadRequest, map[string]string{
				"error": "Results could not be processed",
			})
		}

		features := lo.Map(results, func(result query.Result, _ int) *geojson.Feature {
			wrapper := r.Result{Result: result}
			return wrapper.AsFeature(nil)
		})

		ctx.Response().Header().Set("Cache-Control", "public, max-age=1800")
		ctx.Response().Header().Set("Expires", time.Now().Add(30*time.Minute).Format(http.TimeFormat))

		return response(ctx, http.StatusOK, geojson.FeatureCollection{
			Features: features,
		})
	}
}
