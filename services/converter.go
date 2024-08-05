package services

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jtarchie/knowhere/marshal"
	"github.com/jtarchie/knowhere/query"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcloughlin/geohash"
	"github.com/paulmach/osm"
	"github.com/samber/lo"
	"github.com/valyala/fasttemplate"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Converter struct {
	allowedTags []string
	dbPath      string
	name        string
	osmPath     string
	prefix      string
	rtree       bool
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func NewConverter(
	osmPath string,
	dbPath string,
	prefix string,
	tags []string,
	rtree bool,
) *Converter {
	var allowedTags []string

	if len(tags) > 0 && tags[0] != "*" {
		allowedTags = tags
	}

	caser := cases.Title(language.English, cases.NoLower)
	name := caser.String(nonAlphanumericRegex.ReplaceAllString(prefix, " "))

	return &Converter{
		allowedTags: allowedTags,
		dbPath:      dbPath,
		name:        name,
		osmPath:     osmPath,
		prefix:      strcase.ToSnake(strings.ToLower(prefix)),
		rtree:       rtree,
	}
}

const precision = 100_000

func (b *Converter) Sprintf(template string) string {
	t := fasttemplate.New(template, "{{", "}}")

	return t.ExecuteString(map[string]interface{}{
		"name":   b.name,
		"prefix": b.prefix,
	})
}

func (b *Converter) clientExecute(client *sql.DB, statement string) error {
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

func (b *Converter) Execute() error {
	slog.Info("db.open", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	// open the sql database
	client, err := sql.Open("sqlite3_geohash", b.dbPath)
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
			osm_type INTEGER NOT NULL,
			minLat   REAL,
			maxLat   REAL,
			minLon   REAL,
			maxLon   REAL,
			tags     BLOB,
			refs     BLOB
		) STRICT;

		CREATE TABLE {{prefix}}_refs (
			parent_id    INTEGER NOT NULL,
			osm_id 	     INTEGER NOT NULL,
			osm_type     INTEGER NOT NULL
		) STRICT;

		CREATE TABLE IF NOT EXISTS prefixes (
			id        INTEGER PRIMARY KEY AUTOINCREMENT,
			name      TEXT NOT NULL,
			full_name TEXT NOT NULL,
			minLat    REAL,
			maxLat    REAL,
			minLon    REAL,
			maxLon    REAL
		) STRICT;
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
		(?, ?, ?, ?, ?, ?, jsonb(?), jsonb(?));
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
				query.NodeFilter, // node
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
				query.WayFilter, // way
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

			id, _ := row.LastInsertId()

			for _, node := range way.Nodes {
				_, err = refInsert.Exec(id, node.ID, query.NodeFilter)
				if err != nil {
					return fmt.Errorf("could not create ref for way %d for node %d: %w", id, node.ID, err)
				}
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
				query.RelationFilter, // relation
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

			id, _ := row.LastInsertId()

			for _, member := range relation.Members {
				switch member.Type { //nolint: exhaustive
				case osm.TypeNode:
					_, err = refInsert.Exec(id, member.Ref, query.NodeFilter)
					if err != nil {
						return fmt.Errorf("could not create ref for relation %d for node %d: %w", id, member.Ref, err)
					}
				case osm.TypeWay:
					_, err = refInsert.Exec(id, member.Ref, query.WayFilter)
					if err != nil {
						return fmt.Errorf("could not create ref for relation %d for way %d: %w", id, member.Ref, err)
					}
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
		CREATE UNIQUE INDEX {{prefix}}_osm_types ON {{prefix}}_entries(osm_type, osm_id);

		WITH ways AS (
			SELECT
					e.id
			FROM
					{{prefix}}_entries e
			WHERE
					e.osm_type = 2 -- way
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
					n.osm_type = 1 AND n.osm_id = r.osm_id -- node
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
					e.osm_type = 3 -- relation
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
			ORDER BY
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
		DELETE FROM {{prefix}}_entries WHERE
			minLon IS NULL OR
			maxLon IS NULL OR
			minLat IS NULL OR
			maxLat IS NULL OR
			tags = jsonb('{}');
		DROP INDEX {{prefix}}_ref_ids;
		DROP TABLE {{prefix}}_refs;

		-- calculate geohash and add to tags
		UPDATE {{prefix}}_entries SET tags = jsonb_set(tags, '$.geohash', geohash(minLat, maxLat, minLon, maxLon));
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
			fts5(tags, osm_id, minLat, maxLat, minLon, maxLon, osm_type, content = '{{prefix}}_entries', tokenize="porter", content_rowid='id');

		WITH tags AS (
			SELECT
				{{prefix}}_entries.id AS id,
				json_each.key || ' ' || json_each.value AS kv
			FROM
				{{prefix}}_entries,
				json_each({{prefix}}_entries.tags)
		)
		INSERT INTO
			{{prefix}}_search(rowid, tags)
		SELECT
			id,
			GROUP_CONCAT(kv, ' ')
		FROM
			tags
		GROUP BY
			id;
	`)
	if err != nil {
		return fmt.Errorf("could build full text: %w", err)
	}

	slog.Info("db.fts.complete", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

	if b.rtree {
		slog.Info("db.rtree.init", slog.String("filename", b.dbPath), slog.String("prefix", b.prefix))

		err = b.clientExecute(client, `
			CREATE VIRTUAL TABLE IF NOT EXISTS {{prefix}}_rtree USING rtree(
				id INTEGER PRIMARY KEY,
				minLon REAL,
				maxLon REAL,
				minLat REAL,
				maxLat REAL
			);

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
	}

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

func geohashFunc(minLat, maxLat, minLon, maxLon float64) string {
	return geohash.Encode(
		(minLat+maxLat)/2.0,
		(minLon+maxLon)/2.0,
	)
}

func init() {
	sql.Register("sqlite3_geohash", &sqlite3.SQLiteDriver{
		Extensions: []string{},
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			if err := conn.RegisterFunc("geohash", geohashFunc, true); err != nil {
				return err
			}
			return nil
		},
	})
}
