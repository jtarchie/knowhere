package commands

import (
	"fmt"
	"io"
	"time"

	"github.com/jtarchie/knowhere/server"
)

type Server struct {
	Cors           []string      `help:"list of URLs to allow CORs"`
	AllowCIDR      []string      `help:"list of CIDRs to allow through to server" default:"0.0.0.0/0"`
	DB             string        `help:"sqlite file to server"      required:""`
	Port           int           `default:"8080"                    help:"port for the http server" required:""`
	RuntimeTimeout time.Duration `help:"the timeout for single runtime" default:"2s"`
	CacheSize      int           `help:"SQLite cache size" default:"5000"`
}

func (s *Server) Run(_ io.Writer) error {
	server, err := server.New(
		s.Port,
		s.DB,
		s.Cors,
		s.AllowCIDR,
		s.RuntimeTimeout,
		s.CacheSize,
	)
	if err != nil {
		return fmt.Errorf("could not initialized the server: %w", err)
	}

	err = server.Start()
	if err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
