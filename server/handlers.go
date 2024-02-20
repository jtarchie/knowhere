package server

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
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

//go:embed assets/*
var assets embed.FS

func setupAssets(handler *echo.Echo) error {
	fsys, err := fs.Sub(assets, "assets")
	if err != nil {
		return fmt.Errorf("could not filesystem for assets: %w", err)
	}

	assetHandler := http.FileServer(http.FS(fsys))
	handler.GET("/*", echo.WrapHandler(assetHandler))

	return nil
}

func setupRoutes(handler *echo.Echo, client *sql.DB) {
	handler.POST("/api/search", locationSearch(client))
}
