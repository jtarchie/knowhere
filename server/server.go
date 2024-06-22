package server

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jtarchie/sqlitezstd"
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
	cors []string,
	timeout time.Duration,
) (*Server, error) {
	connectionString := fmt.Sprintf("file:%s?_query_only=true&immutable=true&mode=ro", dbFilename)

	if strings.Contains(dbFilename, ".zst") {
		err := sqlitezstd.Init()
		if err != nil {
			return nil, fmt.Errorf("could not load sqlite zstd vfs: %w", err)
		}

		connectionString += "&vfs=zstd"
	}

	client, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not open database file: %w", err)
	}

	_, err = client.Exec(`
		pragma temp_store = memory;
		pragma mmap_size = 268435456; -- 256 MB
		PRAGMA cache_size = 2000;
		PRAGMA busy_timeout = 5000;
	`)
	if err != nil {
		return nil, fmt.Errorf("could not setup pragmas: %w", err)
	}

	slog.Info(
		"server.config",
		slog.Int("port", port),
		slog.String("db", dbFilename),
	)

	handler := echo.New()
	handler.HideBanner = true
	handler.JSONSerializer = &DefaultJSONSerializer{}

	setupMiddleware(handler, cors)
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
