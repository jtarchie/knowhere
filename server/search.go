package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/jtarchie/knowhere/query"
	"github.com/labstack/echo/v4"
)

func routeSearch(client *sql.DB) func(echo.Context) error {
	return func(ctx echo.Context) error {
		search := ctx.FormValue("search")
		if search == "" {
			//nolint: wrapcheck
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Required to have a `search` for the query",
			})
		}

		results, err := query.Execute(client, search)
		if err != nil {
			slog.Error("search.error", slog.String("error", err.Error()))

			//nolint: wrapcheck
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Results could not be processed",
			})
		}

		//nolint: wrapcheck
		return ctx.JSON(http.StatusOK, results)
	}
}
