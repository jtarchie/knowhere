package services

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jtarchie/knowhere/marshal"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paulmach/osm"
	"github.com/samber/lo"
	"github.com/valyala/fasttemplate"
)

type Builder struct {
	allowedTags []string
	dbPath      string
	name        string
	osmPath     string
	prefix      string
}

func NewBuilder(
	osmPath string,
	dbPath string,
	prefix string,
	tags []string,
) *Builder {
	var allowedTags []string

	if len(tags) > 0 && tags[0] != "*" {
		allowedTags = tags
	}

	return &Builder{
		allowedTags: allowedTags,
		dbPath:      dbPath,
		osmPath:     osmPath,
		prefix:      strcase.ToSnake(strings.ToLower(prefix)),
		name:        prefix,
	}
}

const precision = 100_000

func (b *Builder) Sprintf(template string) string {
	t := fasttemplate.New(template, "{{", "}}")

	return t.ExecuteString(map[string]interface{}{
		"name":   b.name,
		"prefix": b.prefix,
	})
}

func (b *Builder) clientExecute(client *sql.DB, statement string) error {
	queries := strings.Split(statement, ";\n")

	for _, query := range queries {
		query = strings.TrimSpace(b.Sprintf(query))
		if query == "" {
			continue
		}

		slog.Info("db.execute", slog.String("filename", b.dbPath), slog.String("query", query))

		_, err := client.Exec(query)
		if err != nil {
			return fmt.Errorf("could execute query: %w", err)
		}
	}

	return nil
}

func (b *Builder) Execute() error {
	slog.Info("db.open", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	// open the sql database
	client, err := sql.Open("sqlite3", b.dbPath)
	if err != nil {
		return fmt.Errorf("could not open database file: %w", err)
	}
	defer client.Close()

	client.SetMaxOpenConns(1)

	slog.Info("db.schema.create", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	err = b.clientExecute(client, `
		CREATE TABLE {{prefix}}_entries (
			id       INTEGER PRIMARY KEY AUTOINCREMENT,
			osm_id   INTEGER NOT NULL,
			osm_type TEXT NOT NULL,
			minLat   FLOAT,
			maxLat   FLOAT,
			minLon   FLOAT,
			maxLon   FLOAT,
			tags     JSONB,
			refs     JSONB
		);

		CREATE TABLE {{prefix}}_refs (
			parent_id    INTEGER NOT NULL,
			osm_id 	     INTEGER NOT NULL,
			osm_type     TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS prefixes (
			id        INTEGER PRIMARY KEY AUTOINCREMENT,
			name      TEXT NOT NULL,
			full_name TEXT NOT NULL,
			minLat    FLOAT,
			maxLat    FLOAT,
			minLon    FLOAT,
			maxLon    FLOAT
		);

		CREATE VIRTUAL TABLE IF NOT EXISTS {{prefix}}_rtree USING rtree(
			id INTEGER PRIMARY KEY,
			minLon REAL,
			maxLon REAL,
			minLat REAL,
			maxLat REAL
	);
	`)
	if err != nil {
		return fmt.Errorf("could not execute schema: %w", err)
	}

	importer := NewImporter(b.osmPath)

	slog.Info("db.import.init", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	transaction, err := client.Begin()
	if err != nil {
		return fmt.Errorf("could not create a transaction: %w", err)
	}

	defer func() {
		_ = transaction.Rollback()
	}()

	insert, err := transaction.Prepare(b.Sprintf(`
	INSERT INTO {{prefix}}_entries
		(osm_id, osm_type, minLat, maxLat, minLon, maxLon, tags, refs)
			VALUES
		(?, ?, ?, ?, ?, ?, ?, ?);
	`))
	if err != nil {
		return fmt.Errorf("could not create prepared statement for insert: %w", err)
	}

	refInsert, err := transaction.Prepare(b.Sprintf(`
		INSERT INTO {{prefix}}_refs (parent_id, osm_id, osm_type) VALUES (?, ?, ?);
	`))
	if err != nil {
		return fmt.Errorf("could not create prepared statement for ref insert: %w", err)
	}

	err = importer.Execute(
		func(node *osm.Node) error {
			filteredTags := node.TagMap()
			if 0 < len(b.allowedTags) {
				filteredTags = lo.PickByKeys(node.TagMap(), b.allowedTags)
			}

			_, err := insert.Exec(
				node.ID,
				osm.TypeNode,
				math.Round(node.Lat*precision)/precision,
				math.Round(node.Lat*precision)/precision,
				math.Round(node.Lon*precision)/precision,
				math.Round(node.Lon*precision)/precision,
				marshal.Tags(filteredTags),
				nil,
			)
			if err != nil {
				return fmt.Errorf("could not insert node: %w", err)
			}

			return nil
		},
		func(way *osm.Way) error {
			filteredTags := way.TagMap()
			if 0 < len(b.allowedTags) {
				filteredTags = lo.PickByKeys(way.TagMap(), b.allowedTags)
			}

			row, err := insert.Exec(
				way.ID,
				osm.TypeWay,
				nil,
				nil,
				nil,
				nil,
				marshal.Tags(filteredTags),
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
			filteredTags := relation.TagMap()
			if 0 < len(b.allowedTags) {
				filteredTags = lo.PickByKeys(relation.TagMap(), b.allowedTags)
			}

			row, err := insert.Exec(
				relation.ID,
				osm.TypeRelation,
				nil,
				nil,
				nil,
				nil,
				marshal.Tags(filteredTags),
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

	slog.Info("db.import.complete", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	slog.Info("db.bounding_boxes.init", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	err = b.clientExecute(client, `
		CREATE INDEX {{prefix}}_ref_ids ON {{prefix}}_refs(parent_id);
		CREATE INDEX {{prefix}}_osm_types ON {{prefix}}_entries(osm_type, osm_id);

		WITH ways AS (
			SELECT
					e.id
			FROM
					{{prefix}}_entries e
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
					{{prefix}}_refs r
			ON
					r.parent_id = w.id
			JOIN
					{{prefix}}_entries n
			ON
					n.osm_type = 'node' AND n.osm_id = r.osm_id
			GROUP BY
					w.id
		)
		UPDATE
				{{prefix}}_entries
		SET
				minLat = bb.minLat,
				maxLat = bb.maxLat,
				minLon = bb.minLon,
				maxLon = bb.maxLon
		FROM
				bb
		WHERE
				{{prefix}}_entries.id = bb.id;

		WITH relations AS (
			SELECT
					e.id
			FROM
					{{prefix}}_entries e
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
					{{prefix}}_refs r
			ON
					r.parent_id = w.id
			JOIN
					{{prefix}}_entries n
			ON
					n.osm_type = r.osm_type AND n.osm_id = r.osm_id
			GROUP BY
					w.id
		)
		UPDATE
				{{prefix}}_entries
		SET
				minLat = bb.minLat,
				maxLat = bb.maxLat,
				minLon = bb.minLon,
				maxLon = bb.maxLon
		FROM
				bb
		WHERE
				{{prefix}}_entries.id = bb.id;

		-- useful for calculating boundaries,
		-- not useful for searching against
		DELETE FROM {{prefix}}_entries WHERE tags ='{}';
		DROP INDEX {{prefix}}_ref_ids;
		DROP INDEX {{prefix}}_osm_types;
		DROP TABLE {{prefix}}_refs;
	`)
	if err != nil {
		return fmt.Errorf("could not add bounding boxes: %w", err)
	}

	slog.Info("db.bounding_boxes.complete", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	slog.Info("db.fts.init", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	err = b.clientExecute(client, `
		CREATE VIRTUAL TABLE
			{{prefix}}_search
		USING
			fts5(osm_type, tags, content = '{{prefix}}_entries', tokenize="trigram");

		WITH tags AS (
			SELECT
				{{prefix}}_entries.id AS id,
				{{prefix}}_entries.osm_type AS osm_type,
				json_each.key || ' ' || json_each.value AS kv
			FROM
				{{prefix}}_entries,
				json_each({{prefix}}_entries.tags)
		)
		INSERT INTO
			{{prefix}}_search(rowid, osm_type, tags)
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

	slog.Info("db.fts.complete", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	slog.Info("db.rtree.init", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	err = b.clientExecute(client, `
		INSERT INTO
			{{prefix}}_rtree(id, minLon, maxLon, minLat, maxLat)
		SELECT
			id,
			minLon, maxLon,
			minLat, maxLat
		FROM
		{{prefix}}_entries;
	`)
	if err != nil {
		return fmt.Errorf("could build full text: %w", err)
	}

	slog.Info("db.rtree.complete", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	slog.Info("db.optimize.init", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	err = b.clientExecute(client, `
		PRAGMA page_size = 65536;
		PRAGMA cache_size = 4096;
		
		INSERT INTO
			{{prefix}}_search({{prefix}}_search)
		VALUES
			('optimize');

		INSERT INTO prefixes(name, full_name, minLat, maxLat, minLon, maxLon)
			SELECT
				'{{prefix}}', -- Assuming {{prefix}} is replaced with the actual prefix value you want to insert
				'{{name}}', -- Assuming {{name}} is replaced with the actual name value you want to insert
				MIN(minLat),
				MAX(maxLat),
				MIN(minLon),
				MAX(maxLon)
			FROM {{prefix}}_entries;

		vacuum;
		pragma optimize;
	`)
	if err != nil {
		return fmt.Errorf("could not optimize: %w", err)
	}

	slog.Info("db.optimize.complete", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	return nil
}
