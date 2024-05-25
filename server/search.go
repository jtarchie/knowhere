package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/jtarchie/knowhere/query"
	"github.com/labstack/echo/v4"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func locationSearch(client *sql.DB) func(echo.Context) error {
	return func(ctx echo.Context) error {
		search := ctx.FormValue("search")
		if search == "" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Required to have a `search` for the query",
			})
		}

		results, err := query.Execute(client, search, query.ToIndexedSQL)
		if err != nil {
			slog.Error("search.error", slog.String("error", err.Error()))

			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Results could not be processed",
			})
		}

		features := []*geojson.Feature{}

		for _, result := range results {
			var feature geojson.Feature

			feature.ID = result.ID
			feature.Geometry = orb.Point{result.MinLon, result.MinLat}
			feature.Type = "Feature"

			feature.Properties = map[string]interface{}{
				"name": result.Name,
			}

			features = append(features, &feature)
		}

		return ctx.JSON(http.StatusOK, geojson.FeatureCollection{
			Features: features,
		})
	}
}
