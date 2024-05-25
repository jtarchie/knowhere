package server

import (
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/jtarchie/knowhere/query"
	"github.com/labstack/echo/v4"
)

//go:embed turf.js
var turfJSSource string

func runtime(client *sql.DB) func(echo.Context) error {
	vmPool := sync.Pool{
		New: func() any {
			vm := goja.New() //nolint: varnamelen

			_, err := vm.RunString(turfJSSource)
			if err != nil {
				panic(fmt.Sprintf("could not warmup the VM: %s", err))
			}

			return vm
		},
	}

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

		vm := vmPool.Get().(*goja.Runtime) //nolint
		defer vmPool.Put(vm)

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

		value, err := vm.RunString(fmt.Sprintf(`
			(function() {
				%s
			})()
		`, source))
		if err != nil {
			slog.Error("search.error", slog.String("error", err.Error()))

			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Could not evaluate script",
			})
		}

		return ctx.JSON(http.StatusOK, value)
	}
}
