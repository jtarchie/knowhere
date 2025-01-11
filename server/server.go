package server

import (
	"database/sql"
	"fmt"
	"log/slog"
	runtimeNative "runtime"
	"time"

	_ "github.com/jtarchie/sqlitezstd"
	"github.com/labstack/echo/v4"
)

type Server struct {
	port    int
	client  *sql.DB
	handler *echo.Echo
}

func New(
	port int,
	connectionString string,
	cors []string,
	allowCIDRs []string,
	timeout time.Duration,
) (*Server, error) {
	client, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not open database file: %w", err)
	}

	client.SetMaxOpenConns(runtimeNative.NumCPU())
	client.SetMaxIdleConns(runtimeNative.NumCPU())

	slog.Info(
		"server.config",
		slog.Int("port", port),
		slog.String("database_uri", connectionString),
	)

	handler := echo.New()
	handler.HideBanner = true
	handler.JSONSerializer = &DefaultJSONSerializer{}

	setupMiddleware(handler, cors, allowCIDRs)
	setupRoutes(handler, client, timeout)

	return &Server{
		client:  client,
		handler: handler,
		port:    port,
	}, nil
}

func (s *Server) Start() error {
	defer s.client.Close()

	bind := fmt.Sprintf("0.0.0.0:%d", s.port)

	slog.Info("server.started", slog.String("bind", "http://"+bind))

	err := s.handler.Start(bind)
	if err != nil {
		return fmt.Errorf("could not start http server: %w", err)
	}

	return nil
}
