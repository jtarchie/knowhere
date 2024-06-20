package commands

import (
	"fmt"
	"io"
	"time"

	"github.com/jtarchie/knowhere/server"
)

type Server struct {
	Cors           []string      `help:"list of URLs to allow CORs"`
	DB             string        `help:"sqlite file to server"      required:""`
	Port           int           `default:"8080"                    help:"port for the http server" required:""`
	RuntimeTimeout time.Duration `help:"the timeout for single runtime" default:"2s"`
}

func (s *Server) Run(_ io.Writer) error {
	server, err := server.New(
		s.Port,
		s.DB,
		s.Cors,
		s.RuntimeTimeout,
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
