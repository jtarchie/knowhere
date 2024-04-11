package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

func setupMiddleware(handler *echo.Echo, cors []string) {
	handler.Use(slogecho.New(slog.Default()))
	handler.Use(middleware.Recover())
	handler.Use(middleware.Gzip())

	if 0 < len(cors) {
		handler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: cors,
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}))
	}
}

func setupRoutes(handler *echo.Echo, client *sql.DB) {
	handler.GET("/api/search", locationSearch(client))
	handler.GET("/api/prefixes", prefixes(client))
}
