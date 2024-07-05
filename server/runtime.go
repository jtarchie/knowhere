package server

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/jtarchie/knowhere/services"
	"github.com/labstack/echo/v4"
	"github.com/paulmach/orb/geojson"
)

func runtime(
	client *sql.DB,
	timeout time.Duration,
) func(echo.Context) error {
	runtime := services.NewRuntime(client, timeout)

	geojson.CustomJSONMarshaler = &DefaultJSONSerializer{}

	return func(ctx echo.Context) error {
		source := ctx.FormValue("source")
		if source == "" {
			body := &strings.Builder{}

			_, err := io.Copy(body, ctx.Request().Body)
			if err != nil {
				slog.Error("runtime.error", slog.String("error", err.Error()))

				return response(ctx, http.StatusBadRequest, map[string]string{
					"error": "Could not read request body",
				})
			}
			defer ctx.Request().Body.Close()

			source = body.String()
		}

		if source == "" {
			slog.Error("runtime.error", slog.String("error", "source was empty"))

			return response(ctx, http.StatusBadRequest, map[string]string{
				"error": "source not provided in request body",
			})
		}

		value, err := runtime.Execute(source)
		if err != nil {
			slog.Error(
				"runtime.error",
				slog.String("error", err.Error()),
				slog.String("type", fmt.Sprintf("%#v", err)),
			)

			msg := "evaluation error"

			var exception *goja.Exception
			if errors.As(err, &exception) {
				msg += ": " + exception.Error()
			}

			var interrupted *goja.InterruptedError
			if errors.As(err, &interrupted) {
				msg += ": " + interrupted.Error()
			}

			return response(ctx, http.StatusBadRequest, map[string]string{
				"error": msg,
			})
		}

		return response(ctx, http.StatusOK, value.Export())
	}
}
