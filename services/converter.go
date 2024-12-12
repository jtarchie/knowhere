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
	area        string
	rtree       bool
	optimize    bool
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func NewConverter(
	osmPath string,
	dbPath string,
	area string,
	tags []string,
	rtree bool,
	optimize bool,
) *Converter {
	var allowedTags []string

	if len(tags) > 0 && tags[0] != "*" {
		allowedTags = tags
	}

	caser := cases.Title(language.English, cases.NoLower)
	name := caser.String(nonAlphanumericRegex.ReplaceAllString(area, " "))

	return &Converter{
		allowedTags: allowedTags,
		dbPath:      dbPath,
		name:        name,
		osmPath:     osmPath,
		area:        strcase.ToSnake(strings.ToLower(area)),
		rtree:       rtree,
		optimize:    optimize,
	}
}

const precision = 100_000

func (b *Converter) Sprintf(template string) string {
	t := fasttemplate.New(template, "{{", "}}")

	return t.ExecuteString(map[string]interface{}{
		"name": b.name,
		"area": b.area,
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
	slog.Info("db.open", slog.String("filename", b.dbPath), slog.String("area", b.area))

	// open the sql database
	client, err := sql.Open("sqlite3_geohash", b.dbPath)
	if err != nil {
		return fmt.Errorf("could not open database file: %w", err)
	}
	defer client.Close()

	client.SetMaxOpenConns(1)

	// Speed up inserts by adjusting PRAGMA settings
	_, err = client.Exec(`
		PRAGMA synchronous = OFF;
		PRAGMA journal_mode = MEMORY;
		PRAGMA temp_store = MEMORY;
	`)
	if err != nil {
		return fmt.Errorf("could not set PRAGMA settings: %w", err)
	}

	slog.Info("db.schema.create", slog.String("filename", b.dbPath), slog.String("area", b.area))

	err = b.clientExecute(client, `
		CREATE TABLE {{area}}_entries (
			id       INTEGER PRIMARY KEY AUTOINCREMENT,
			osm_id   INTEGER NOT NULL,
			osm_type INTEGER NOT NULL,
			minLat   REAL,
			maxLat   REAL,
			minLon   REAL,
			maxLon   REAL,
			tags     BLOB
		) STRICT;

		CREATE TABLE {{area}}_refs (
			parent_id    INTEGER NOT NULL,
			osm_id 	     INTEGER NOT NULL,
			osm_type     INTEGER NOT NULL
		) STRICT;

		CREATE TABLE IF NOT EXISTS areas (
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

	slog.Info("db.import.init", slog.String("filename", b.dbPath), slog.String("area", b.area))

	transaction, err := client.Begin()
	if err != nil {
		return fmt.Errorf("could not create a transaction: %w", err)
	}

	defer func() {
		_ = transaction.Rollback()
	}()

	insert, err := transaction.Prepare(b.Sprintf(`
	INSERT INTO {{area}}_entries
		(osm_id, osm_type, minLat, maxLat, minLon, maxLon, tags)
			VALUES
		(?, ?, ?, ?, ?, ?, jsonb(?));
	`))
	if err != nil {
		return fmt.Errorf("could not create prepared statement for insert: %w", err)
	}

	refInsert, err := transaction.Prepare(b.Sprintf(`
		INSERT INTO {{area}}_refs (parent_id, osm_id, osm_type) VALUES (?, ?, ?);
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

	slog.Info("db.import.complete", slog.String("filename", b.dbPath), slog.String("area", b.area))

	slog.Info("db.bounding_boxes.init", slog.String("filename", b.dbPath), slog.String("area", b.area))

	err = b.clientExecute(client, `
		CREATE INDEX {{area}}_ref_ids ON {{area}}_refs(parent_id);
		CREATE UNIQUE INDEX {{area}}_osm_types ON {{area}}_entries(osm_type, osm_id);

		WITH ways AS (
			SELECT
					e.id
			FROM
					{{area}}_entries e
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
					{{area}}_refs r
			ON
					r.parent_id = w.id
			JOIN
					{{area}}_entries n
			ON
					n.osm_type = 1 AND n.osm_id = r.osm_id -- node
			GROUP BY
					w.id
		)
		UPDATE
				{{area}}_entries
		SET
				minLat = bb.minLat,
				maxLat = bb.maxLat,
				minLon = bb.minLon,
				maxLon = bb.maxLon
		FROM
				bb
		WHERE
				{{area}}_entries.id = bb.id;

		WITH relations AS (
			SELECT
					e.id
			FROM
					{{area}}_entries e
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
					{{area}}_refs r
			ON
					r.parent_id = w.id
			JOIN
					{{area}}_entries n
			ON
					n.osm_type = r.osm_type AND n.osm_id = r.osm_id
			GROUP BY
					w.id
			ORDER BY
					w.id
		)
		UPDATE
				{{area}}_entries
		SET
				minLat = bb.minLat,
				maxLat = bb.maxLat,
				minLon = bb.minLon,
				maxLon = bb.maxLon
		FROM
				bb
		WHERE
				{{area}}_entries.id = bb.id;

		-- useful for calculating boundaries,
		-- not useful for searching against
		DELETE FROM {{area}}_entries WHERE
			minLon IS NULL OR
			maxLon IS NULL OR
			minLat IS NULL OR
			maxLat IS NULL OR
			tags = jsonb('{}');
		DROP INDEX {{area}}_ref_ids;
		DROP TABLE {{area}}_refs;

		-- calculate geohash and add to tags
		UPDATE {{area}}_entries SET tags = jsonb_set(tags, '$.geohash', geohash(minLat, maxLat, minLon, maxLon));
	`)
	if err != nil {
		return fmt.Errorf("could not add bounding boxes: %w", err)
	}

	slog.Info("db.bounding_boxes.complete", slog.String("filename", b.dbPath), slog.String("area", b.area))

	slog.Info("db.fts.init", slog.String("filename", b.dbPath), slog.String("area", b.area))

	err = b.clientExecute(client, `
		CREATE VIRTUAL TABLE
			{{area}}_search
		USING
			fts5(tags, osm_id, minLat, maxLat, minLon, maxLon, osm_type, content = '{{area}}_entries', tokenize="unicode61", content_rowid='id');

		WITH tags AS (
			SELECT
				{{area}}_entries.id AS id,
				json_each.key || ' ' || json_each.value AS kv
			FROM
				{{area}}_entries,
				json_each({{area}}_entries.tags)
		)
		INSERT INTO
			{{area}}_search(rowid, tags)
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

	slog.Info("db.fts.complete", slog.String("filename", b.dbPath), slog.String("area", b.area))

	if b.rtree {
		slog.Info("db.rtree.init", slog.String("filename", b.dbPath), slog.String("area", b.area))

		err = b.clientExecute(client, `
			CREATE VIRTUAL TABLE IF NOT EXISTS {{area}}_rtree USING rtree(
				id INTEGER PRIMARY KEY,
				minLon REAL,
				maxLon REAL,
				minLat REAL,
				maxLat REAL
			);

			INSERT INTO
				{{area}}_rtree(id, minLon, maxLon, minLat, maxLat)
			SELECT
				id,
				minLon, maxLon,
				minLat, maxLat
			FROM
			{{area}}_entries;
		`)
		if err != nil {
			return fmt.Errorf("could build full text: %w", err)
		}

		slog.Info("db.rtree.complete", slog.String("filename", b.dbPath), slog.String("area", b.area))
	}

	slog.Info("db.optimize.init", slog.String("filename", b.dbPath), slog.String("area", b.area))

	err = b.clientExecute(client, `
		pragma page_size = 16384;
		
		INSERT INTO
			{{area}}_search({{area}}_search)
		VALUES
			('optimize');

		INSERT INTO areas(name, full_name, minLat, maxLat, minLon, maxLon)
			SELECT
				'{{area}}', -- Assuming {{area}} is replaced with the actual area value you want to insert
				'{{name}}', -- Assuming {{name}} is replaced with the actual name value you want to insert
				MIN(minLat),
				MAX(maxLat),
				MIN(minLon),
				MAX(maxLon)
			FROM {{area}}_entries;
	`)
	if err != nil {
		return fmt.Errorf("could not optimize data: %w", err)
	}

	slog.Info("db.optimize.complete", slog.String("filename", b.dbPath), slog.String("area", b.area))

	if b.optimize {
		slog.Info("db.optimize.init", slog.String("filename", b.dbPath))
		err = b.clientExecute(client, `
			PRAGMA temp_store = FILE;
			VACUUM;
			PRAGMA optimize;
		`)
		if err != nil {
			return fmt.Errorf("could not optimize db: %w", err)
		}
		slog.Info("db.optimize.completed", slog.String("filename", b.dbPath))
	}

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
