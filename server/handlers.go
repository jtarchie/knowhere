package server

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

func setupMiddleware(handler *echo.Echo, cors []string, allowedCIDRs []string) {
	handler.Use(slogecho.New(slog.Default()))
	handler.Use(middleware.Recover())
	handler.Use(middleware.Gzip())

	if 0 < len(cors) && cors[0] != "*" {
		slog.Info("cors.setup", "cors", cors)
		handler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: cors,
		}))
	}

	if 0 < len(allowedCIDRs) && allowedCIDRs[0] != "0.0.0.0/0" {
		slog.Info("cidrs.setup", "allows", allowedCIDRs)
		handler.Use(CIDRAllow(allowedCIDRs...))
	}
}

func setupRoutes(
	handler *echo.Echo,
	client *sql.DB,
	timeout time.Duration,
) {
	handler.GET("/api/search", locationSearch(client))
	handler.GET("/api/areas", areas(client))
	handler.Any("/api/runtime", runtime(client, timeout))
	handler.GET("/up", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	assetHandler := assetHandler()
	handler.GET("/", echo.WrapHandler(assetHandler))
	handler.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))
}
