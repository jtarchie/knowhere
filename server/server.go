package server

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type Server struct {
	port    int
	client  *sql.DB
	handler *echo.Echo
}

func New(
	port int,
	dbFilename string,
) (*Server, error) {
	client, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_query_only=true&immutable=true&mode=ro", dbFilename))
	if err != nil {
		return nil, fmt.Errorf("could not open database file: %w", err)
	}

	slog.Info(
		"server.config",
		slog.Int("port", port),
		slog.String("db", dbFilename),
	)

	handler := echo.New()
	handler.HideBanner = true
	handler.JSONSerializer = DefaultJSONSerializer{}

	setupMiddleware(handler)

	err = setupAssets(handler)
	if err != nil {
		return nil, fmt.Errorf("could not attach assets: %w", err)
	}

	setupRoutes(handler, client)

	return &Server{
		client:  client,
		handler: handler,
		port:    port,
	}, nil
}

func (s *Server) Start() error {
	bind := fmt.Sprintf("0.0.0.0:%d", s.port)

	slog.Info("server.started", slog.String("bind", fmt.Sprintf("http://%s", bind)))

	err := s.handler.Start(bind)
	if err != nil {
		return fmt.Errorf("could not start http server: %w", err)
	}

	return nil
}
