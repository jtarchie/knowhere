package server

import (
	"database/sql"
	_ "embed"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jtarchie/knowhere/services"
	"github.com/labstack/echo/v4"
)

func runtime(client *sql.DB) func(echo.Context) error {
	runtime := services.NewRuntime(client)

	return func(ctx echo.Context) error {
		body := &strings.Builder{}

		_, err := io.Copy(body, ctx.Request().Body)
		if err != nil {
			slog.Error("search.error", slog.String("error", err.Error()))

			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Could not read request body",
			})
		}
		defer ctx.Request().Body.Close()

		source := body.String()

		if source == "" {
			slog.Error("search.error", slog.String("error", "source was empty"))

			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "source was empty",
			})
		}

		if source == "" {
			slog.Error("search.error", slog.String("error", err.Error()))

			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "could not execute script",
			})
		}

		value, err := runtime.Execute(source)
		if err != nil {
			slog.Error("search.error", slog.String("error", err.Error()))

			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Could not evaluate script",
			})
		}

		return ctx.JSON(http.StatusOK, value.Export())
	}
}
