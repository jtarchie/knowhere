package commands

import (
	"fmt"

	"github.com/jtarchie/knowhere/services"
)

type Build struct {
	OSM string `help:"osm pbf file to build the sqlite file from" required:"" type:"existingfile"`
	DB  string `help:"db filename to import data from"            required:""`
}

func (b *Build) Run() error {
	builder := services.NewBuilder(b.OSM, b.DB)

	err := builder.Execute()
	if err != nil {
		return fmt.Errorf("could not build database from OSM: %w", err)
	}

	return nil
}
