package services

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"github.com/schollz/progressbar/v3"
)

type Importer struct {
	filename string
}

func NewImporter(
	filename string,
) *Importer {
	return &Importer{
		filename: filename,
	}
}

func (i *Importer) Execute(
	nfn func(*osm.Node) error,
	wfn func(*osm.Way) error,
	rfn func(*osm.Relation) error,
) error {
	file, err := os.Open(i.filename)
	if err != nil {
		return fmt.Errorf("could not open osm pbf file: %w", err)
	}
	defer file.Close()

	stat, err := os.Stat(i.filename)
	if err != nil {
		return fmt.Errorf("could stat the osm pbf file: %w", err)
	}

	scanner := osmpbf.New(context.Background(), file, runtime.NumCPU())
	defer scanner.Close()

	bar := progressbar.DefaultBytes(
		stat.Size(),
		"importing",
	)

	for scanner.Scan() {
		switch object := scanner.Object().(type) {
		case *osm.Node:
			err := nfn(object)
			if err != nil {
				return fmt.Errorf("could not import node %d: %w", object.ID, err)
			}
		case *osm.Way:
			err := wfn(object)
			if err != nil {
				return fmt.Errorf("could not import way %d: %w", object.ID, err)
			}
		case *osm.Relation:
			err := rfn(object)
			if err != nil {
				return fmt.Errorf("could not import relation %d: %w", object.ID, err)
			}
		}

		_ = bar.Set64(scanner.FullyScannedBytes())
	}

	_ = bar.Finish()

	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("scanner had error: %w", err)
	}

	return nil
}
