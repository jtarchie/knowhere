package commands

import (
	"fmt"

	"github.com/jtarchie/knowhere/services"
)

type Build struct {
	OSM    string `help:"osm pbf file to build the sqlite file from" required:"" type:"existingfile"`
	DB     string `help:"db filename to import data from"            required:""`
	Prefix string `help:"will add this prefix to all table names"    required:""`
}

func (b *Build) Run() error {
	builder := services.NewBuilder(b.OSM, b.DB, b.Prefix)

	err := builder.Execute()
	if err != nil {
		return fmt.Errorf("could not build database from OSM: %w", err)
	}

	return nil
}
