package commands

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/jtarchie/knowhere/query"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"github.com/wroge/scan"
)

type Server struct {
	Port int    `default:"8080"               help:"port for the http server" required:""`
	DB   string `help:"sqlite file to server" required:""                     type:"existingfile"`
}

//go:embed assets/*
var assets embed.FS

type Result struct {
	ID      int64   `json:"id"`
	MinLat  float64 `json:"minLat"`
	MaxLat  float64 `json:"maxLat"`
	MinLon  float64 `json:"minLon"`
	MaxLon  float64 `json:"maxLon"`
	OsmType string  `json:"type"`
}

func (s *Server) Run() error {
	client, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_query_only=true&immutable=true&mode=ro", s.DB))
	if err != nil {
		return fmt.Errorf("could not open database file: %w", err)
	}
	defer client.Close()

	server := echo.New()

	server.Use(slogecho.New(slog.Default()))
	server.Use(middleware.Recover())

	fsys, _ := fs.Sub(assets, "assets")
	assetHandler := http.FileServer(http.FS(fsys))

	server.GET("/", echo.WrapHandler(assetHandler))

	server.POST("/api/search", func(ctx echo.Context) error {
		search := ctx.FormValue("search")
		if search == "" {
			//nolint: wrapcheck
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Required to have a `search` for the query",
			})
		}

		sql, err := query.ToSQL(search)
		if err != nil {
			slog.Error("parse.error", slog.String("search", search), slog.String("error", err.Error()))

			//nolint: wrapcheck
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("Could not parse `search`: %s", err.Error()),
			})
		}

		start := time.Now()

		rows, err := client.Query(sql)
		if rows.Err() != nil {
			slog.Error("query.error", slog.String("search", search), slog.String("error", err.Error()))

			//nolint: wrapcheck
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Could not read from database",
			})
		}

		if err != nil {
			slog.Error("query.error", slog.String("search", search), slog.String("error", err.Error()))

			//nolint: wrapcheck
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Could not read from database",
			})
		}

		columns := map[string]scan.Scanner[Result]{
			"id":       scan.Any(func(p *Result, id int64) { p.ID = id }),
			"minLat":   scan.Any(func(p *Result, minLat float64) { p.MinLat = minLat }),
			"maxLat":   scan.Any(func(p *Result, maxLat float64) { p.MaxLat = maxLat }),
			"minLon":   scan.Any(func(p *Result, minLon float64) { p.MinLon = minLon }),
			"maxLon":   scan.Any(func(p *Result, maxLon float64) { p.MaxLon = maxLon }),
			"osm_type": scan.Any(func(p *Result, osmType string) { p.OsmType = osmType }),
		}

		results, err := scan.All(rows, columns)
		if err != nil {
			slog.Info("scan.error", slog.String("search", search), slog.String("error", err.Error()))

			//nolint: wrapcheck
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Could not read from database",
			})
		}

		slog.Info(
			"query.complete",
			slog.String("search", search),
			slog.String("sql", sql),
			slog.Duration("took", time.Since(start)),
		)

		//nolint: wrapcheck
		return ctx.JSON(http.StatusOK, results)
	})

	err = server.Start(fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return fmt.Errorf("server stopped: %w", err)
	}

	return nil
}
