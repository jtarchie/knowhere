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
			osm_type TEXT NOT NULL,
			minLat   FLOAT,
			maxLat   FLOAT,
			minLon   FLOAT,
			maxLon   FLOAT,
			tags     JSON,
			refs     JSON
		);

		CREATE TABLE refs (
			parent_id    INTEGER NOT NULL,
			osm_id 	     INTEGER NOT NULL,
			osm_type     TEXT NOT NULL
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

	refInsert, err := transaction.Prepare(`
		INSERT INTO refs (parent_id, osm_id, osm_type) VALUES (?, ?, ?);
	`)
	if err != nil {
		return fmt.Errorf("could not create prepared statement for ref insert: %w", err)
	}

	err = importer.Execute(
		func(node *osm.Node) error {
			tags := marshal.Tags(node.TagMap())
			_, err := insert.Exec(
				node.ID,
				osm.TypeNode,
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
			row, err := insert.Exec(
				way.ID,
				osm.TypeWay,
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

			for _, node := range way.Nodes {
				id, _ := row.LastInsertId()
				_, _ = refInsert.Exec(id, node.ID, osm.TypeNode)
			}

			return nil
		},
		func(relation *osm.Relation) error {
			row, err := insert.Exec(
				relation.ID,
				osm.TypeRelation,
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

			for _, member := range relation.Members {
				switch member.Type { //nolint: exhaustive
				case osm.TypeNode, osm.TypeWay:
					id, _ := row.LastInsertId()
					_, _ = refInsert.Exec(id, member.Ref, member.Type)
				}
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

	slog.Info("db.bounding_boxes.init", slog.String("filename", b.dbPath))

	_, err = client.Exec(`
		CREATE INDEX ref_ids ON refs(parent_id);
		CREATE INDEX osm_types ON entries(osm_type, osm_id);

		WITH ways AS (
			SELECT
					e.id
			FROM
					entries e
			WHERE
					e.osm_type = 'way'
		), bb AS (
			SELECT
					w.id AS id,
					MIN(n.minLat) AS minLat,
					MAX(n.maxLat) AS maxLat,
					MIN(n.minLon) AS minLon,
					MAX(n.maxLon) AS maxLon
			FROM
					ways w
			JOIN
					refs r
			ON
					r.parent_id = w.id
			JOIN
					entries n
			ON
					n.osm_type = 'node' AND n.osm_id = r.osm_id
			GROUP BY
					w.id
		)
		UPDATE
				entries
		SET
				minLat = bb.minLat,
				maxLat = bb.maxLat,
				minLon = bb.minLon,
				maxLon = bb.maxLon
		FROM
				bb
		WHERE
				entries.id = bb.id;

		WITH relations AS (
			SELECT
					e.id
			FROM
					entries e
			WHERE
					e.osm_type = 'relation'
		), bb AS (
			SELECT
					w.id AS id,
					MIN(n.minLat) AS minLat,
					MAX(n.maxLat) AS maxLat,
					MIN(n.minLon) AS minLon,
					MAX(n.maxLon) AS maxLon
			FROM
					relations w
			JOIN
					refs r
			ON
					r.parent_id = w.id
			JOIN
					entries n
			ON
					n.osm_type = r.osm_type AND n.osm_id = r.osm_id
			GROUP BY
					w.id
		)
		UPDATE
				entries
		SET
				minLat = bb.minLat,
				maxLat = bb.maxLat,
				minLon = bb.minLon,
				maxLon = bb.maxLon
		FROM
				bb
		WHERE
				entries.id = bb.id;

		DROP INDEX ref_ids;
		DROP INDEX osm_types;
		DROP TABLE refs;
	`)
	if err != nil {
		return fmt.Errorf("could not add bounding boxes: %w", err)
	}

	slog.Info("db.bounding_boxes.complete", slog.String("filename", b.dbPath))

	slog.Info("db.fts.init", slog.String("filename", b.dbPath))

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
	`)
	if err != nil {
		return fmt.Errorf("could build full text: %w", err)
	}

	slog.Info("db.fts.complete", slog.String("filename", b.dbPath))

	slog.Info("db.optimize.init", slog.String("filename", b.dbPath))

	_, err = client.Exec(`
		INSERT INTO
			search(search)
		VALUES
			('optimize');

		vacuum;
		pragma optimize;
	`)
	if err != nil {
		return fmt.Errorf("could not optimize: %w", err)
	}

	slog.Info("db.optimize.complete", slog.String("filename", b.dbPath))

	return nil
}
