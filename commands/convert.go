package commands

import (
	"fmt"
	"io"

	"github.com/jtarchie/knowhere/services"
)

type Convert struct {
	OSM         string   `help:"osm pbf file to build the sqlite file from" required:""                                               type:"existingfile"`
	DB          string   `help:"db filename to import data to"              required:""`
	Prefix      string   `help:"will add this prefix to all table names"    required:""`
	AllowedTags []string `default:"*"                                       help:"a list of allowed tags, all other will be filtered"`
}

func (b *Convert) Run(_ io.Writer) error {
	builder := services.NewConverter(b.OSM, b.DB, b.Prefix, b.AllowedTags)

	err := builder.Execute()
	if err != nil {
		return fmt.Errorf("could not build database from OSM: %w", err)
	}

	return nil
}
