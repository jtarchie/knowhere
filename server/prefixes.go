package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Prefix struct {
	Name     string  `json:"name"`
	FullName string  `json:"full_name"`
	MinLat   float32 `json:"min_lat"`
	MaxLat   float32 `json:"max_lat"`
	MinLon   float32 `json:"min_lon"`
	MaxLon   float32 `json:"max_lon"`
}

func prefixes(client *sql.DB) func(echo.Context) error {
	return func(ctx echo.Context) error {
		rows, err := client.QueryContext(ctx.Request().Context(), `
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

		if rows.Err() != nil {
			slog.Error("prefixes.error", slog.String("error", rows.Err().Error()))

			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Results could not be processed",
			})
		}
		defer rows.Close()

		var (
			prefix   Prefix
			prefixes []Prefix
		)

		for rows.Next() {
			err = rows.Scan(&prefix.Name, &prefix.FullName, &prefix.MinLat, &prefix.MaxLat, &prefix.MinLon, &prefix.MaxLon)
			if err != nil {
				slog.Error("prefixes.error", slog.String("error", err.Error()))

				return ctx.JSON(http.StatusBadRequest, map[string]string{
					"error": "Results could not be processed",
				})
			}

			prefixes = append(prefixes, prefix)
		}
		defer rows.Close()

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"prefixes": prefixes,
		})
	}
}
