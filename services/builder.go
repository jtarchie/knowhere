package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paulmach/osm"
)

type Builder struct {
	osmPath string
	dbPath  string
}

func NewBuilder(
	osmPath string,
	dbPath string,
) *Builder {
	return &Builder{
		osmPath: osmPath,
		dbPath:  dbPath,
	}
}

func (b *Builder) Execute() error {
	slog.Info("db.open", slog.String("filename", b.dbPath))

	// open the sql database
	client, err := sql.Open("sqlite3", b.dbPath)
	if err != nil {
		return fmt.Errorf("could not open database file: %w", err)
	}
	defer client.Close()

	slog.Info("db.schema.create", slog.String("filename", b.dbPath))

	_, err = client.Exec(`
		CREATE TABLE entries (
			id       INTEGER PRIMARY KEY AUTOINCREMENT,
			osm_id   INTEGER NOT NULL,
			osm_type INTEGER NOT NULL,
			minLat   FLOAT,
			maxLat   FLOAT,
			minLon   FLOAT,
			maxLon   FLOAT,
			tags     JSON,
			refs     JSON
		);
	`)
	if err != nil {
		return fmt.Errorf("could not execute schema: %w", err)
	}

	importer := NewImporter(b.osmPath)

	slog.Info("db.import", slog.String("filename", b.dbPath))

	err = importer.Execute(
		func(node *osm.Node) error {
			contents, err := json.Marshal(node.TagMap())
			if err != nil {
				return fmt.Errorf("could not marshal tag map for node %d: %w", node.ID, err)
			}

			_, err = client.Exec(
				`INSERT INTO entries
					(osm_id, osm_type, minLat, maxLat, minLon, maxLon, tags)
						VALUES
					(?, ?, ?, ?, ?, ?, ?);
			`, node.ID, "node", node.Lat, node.Lat, node.Lon, node.Lon, string(contents))
			if err != nil {
				return fmt.Errorf("could not insert node: %w", err)
			}

			return nil
		},
		func(way *osm.Way) error {
			contents, err := json.Marshal(way.TagMap())
			if err != nil {
				return fmt.Errorf("could not marshal tag map for node %d: %w", way.ID, err)
			}

			_, err = client.Exec(
				`INSERT INTO entries
					(osm_id, osm_type, tags)
						VALUES
					(?, ?, ?);
			`, way.ID, "way", string(contents))
			if err != nil {
				return fmt.Errorf("could not insert node: %w", err)
			}

			return nil
		},
		func(relation *osm.Relation) error {
			contents, err := json.Marshal(relation.TagMap())
			if err != nil {
				return fmt.Errorf("could not marshal tag map for node %d: %w", relation.ID, err)
			}

			_, err = client.Exec(
				`INSERT INTO entries
					(osm_id, osm_type, tags)
						VALUES
					(?, ?, ?);
			`, relation.ID, "way", string(contents))
			if err != nil {
				return fmt.Errorf("could not insert node: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("could not import into the database: %w", err)
	}

	slog.Info("db.import.complete", slog.String("filename", b.dbPath))

	// insert nodes, ways, and relations from osm pbf file
	// close db
	return nil
}
