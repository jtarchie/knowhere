package services

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jtarchie/knowhere/marshal"
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

	client.SetMaxOpenConns(1)

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

	slog.Info("db.import.init", slog.String("filename", b.dbPath))

	transaction, err := client.Begin()
	if err != nil {
		return fmt.Errorf("could not create a transaction: %w", err)
	}

	defer func() {
		_ = transaction.Rollback()
	}()

	insert, err := transaction.Prepare(`
	INSERT INTO entries
		(osm_id, osm_type, minLat, maxLat, minLon, maxLon, tags, refs)
			VALUES
		(?, ?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		return fmt.Errorf("could not create prepared statement for insert: %w", err)
	}

	err = importer.Execute(
		func(node *osm.Node) error {
			tags := marshal.Tags(node.TagMap())
			_, err = insert.Exec(
				node.ID,
				"node",
				node.Lat,
				node.Lat,
				node.Lon,
				node.Lon,
				tags,
				nil,
			)
			if err != nil {
				return fmt.Errorf("could not insert node: %w", err)
			}

			return nil
		},
		func(way *osm.Way) error {
			_, err = insert.Exec(
				way.ID,
				"way",
				nil,
				nil,
				nil,
				nil,
				marshal.Tags(way.TagMap()),
				marshal.WayNodes(way.Nodes),
			)
			if err != nil {
				return fmt.Errorf("could not insert node: %w", err)
			}

			return nil
		},
		func(relation *osm.Relation) error {
			_, err = insert.Exec(
				relation.ID,
				"relation",
				nil,
				nil,
				nil,
				nil,
				marshal.Tags(relation.TagMap()),
				marshal.Members(relation.Members),
			)
			if err != nil {
				return fmt.Errorf("could not insert node: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("could not import into the database: %w", err)
	}

	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	slog.Info("db.import.complete", slog.String("filename", b.dbPath))

	slog.Info("db.optimize.init", slog.String("filename", b.dbPath))

	_, err = client.Exec(`
		CREATE VIRTUAL TABLE
			search
		USING
			fts5(osm_type, tags, content = 'entries');

		WITH tags AS (
			SELECT
				entries.id AS id,
				entries.osm_type AS osm_type,
				json_each.key || ' ' || json_each.value AS kv
			FROM
				entries,
				json_each(entries.tags)
		)
		INSERT INTO
			search(rowid, osm_type, tags)
		SELECT
			id,
			osm_type,
			GROUP_CONCAT(kv, ' ')
		FROM
			tags
		GROUP BY
			id, osm_type;
		
		pragma vacuum;
		pragma optimize;
	`)
	if err != nil {
		return fmt.Errorf("could not optimize database: %w", err)
	}

	slog.Info("db.optimize.complete", slog.String("filename", b.dbPath))

	return nil
}
