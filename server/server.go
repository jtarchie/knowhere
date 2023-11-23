package server

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Server struct {
	port    int
	client  *sql.DB
	handler *echo.Echo
}

func New(
	port int,
	db string,
) (*Server, error) {
	client, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_query_only=true&immutable=true&mode=ro", db))
	if err != nil {
		return nil, fmt.Errorf("could not open database file: %w", err)
	}

	handler := echo.New()

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
	err := s.handler.Start(fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("could not start http server: %w", err)
	}

	return nil
}
