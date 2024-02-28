package commands

import (
	"fmt"

	"github.com/jtarchie/knowhere/server"
)

type Server struct {
	Port int      `default:"8080"                    help:"port for the http server" required:""`
	DB   string   `help:"sqlite file to server"      required:""`
	Cors []string `help:"list of URLs to allow CORs"`
}

func (s *Server) Run() error {
	server, err := server.New(
		s.Port,
		s.DB,
		s.Cors,
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
