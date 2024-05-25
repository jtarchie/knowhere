package server

import (
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dop251/goja"
	"github.com/jtarchie/knowhere/query"
	"github.com/labstack/echo/v4"
)

func runtime(client *sql.DB) func(echo.Context) error {
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

		vm := goja.New() //nolint: varnamelen

		err = vm.Set("execute", func(qs string) interface{} {
			results, err := query.Execute(client, qs)
			if err != nil {
				return map[string]string{
					"error": "The query could not be executed.",
				}
			}

			return results
		})

		if source == "" {
			slog.Error("search.error", slog.String("error", err.Error()))

			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "could not execute script",
			})
		}

		value, err := vm.RunString(source)
		if err != nil {
			slog.Error("search.error", slog.String("error", err.Error()))

			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Could not evaluate script",
			})
		}

		results, ok := value.Export().([]interface{})
		if !ok {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Final evaluation not result object",
			})
		}

		return ctx.JSON(http.StatusOK, results)
	}
}
