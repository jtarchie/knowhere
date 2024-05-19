package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/labstack/echo/v4"
)

type Prefix struct {
	Name     string  `db:"name"      json:"name"`
	FullName string  `db:"full_name" json:"full_name"`
	MinLat   float32 `db:"minLat"    json:"min_lat"`
	MaxLat   float32 `db:"maxLat"    json:"max_lat"`
	MinLon   float32 `db:"minLon"    json:"min_lon"`
	MaxLon   float32 `db:"maxLon"    json:"max_lon"`
}

func prefixes(client *sql.DB) func(echo.Context) error {
	return func(ctx echo.Context) error {
		var prefixes []Prefix
		
		err := sqlscan.Select(
			ctx.Request().Context(),
			client,
			&prefixes, `
			SELECT
				name, full_name, minLat, maxLat, minLon, maxLon
			FROM
				prefixes
		`)
		if err != nil {
			slog.Error("prefixes.error", slog.String("error", err.Error()))

			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Results could not be processed",
			})
		}

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"prefixes": prefixes,
		})
	}
}
