package commands

import (
	"fmt"
	"io"

	"github.com/jtarchie/knowhere/services"
)

type Convert struct {
	AllowedTags []string `default:"*"                                       help:"a list of allowed tags, all other will be filtered"`
	DB          string   `help:"db filename to import data to"              required:""`
	OSM         string   `help:"osm pbf file to build the sqlite file from" required:""                                               type:"existingfile"`
	Prefix      string   `help:"will add this area to all table names"    required:""`
	Rtree       bool     `help:"enable rtree index"`
	OptimizeDB  bool     `help:"optimize the database after import" default:"true"`
}

func (b *Convert) Run(_ io.Writer) error {
	builder := services.NewConverter(b.OSM, b.DB, b.Prefix, b.AllowedTags, b.Rtree, b.OptimizeDB)

	err := builder.Execute()
	if err != nil {
		return fmt.Errorf("could not build database from OSM: %w", err)
	}

	return nil
}
