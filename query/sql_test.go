package query_test

import (
	"database/sql"
	"regexp"
	"strings"

	"github.com/jtarchie/knowhere/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/mattn/go-sqlite3"
)

var _ = Describe("Build SQL from a query", Ordered, func() {
	var client *sql.DB

	BeforeAll(func() {
		var err error

		client, err = sql.Open("sqlite3", ":memory:")
		Expect(err).NotTo(HaveOccurred())

		_, err = client.Exec(`
			CREATE TABLE entries (
				id, osm_type, osm_id, tags, minLat, maxLat, minLon, maxLon
			);
			CREATE VIRTUAL TABLE
				search
			USING
				fts5(tags, osm_type, osm_id, minLat, maxLat, minLon, maxLon, content = 'entries', tokenize="porter", content_rowid='id');
			CREATE TABLE test_entries (
				id, osm_type, osm_id, tags, minLat, maxLat, minLon, maxLon
			);
			CREATE VIRTUAL TABLE
				test_search
			USING
				fts5(tags, osm_type, osm_id, minLat, maxLat, minLon, maxLon, content = 'test_entries', tokenize="porter", content_rowid='id');
		`)
		Expect(err).NotTo(HaveOccurred())
	})

	pretty := func(sql string) string {
		space := regexp.MustCompile(`\s+`)

		return strings.TrimSpace(space.ReplaceAllString(sql, " "))
	}

	Describe("ToIndexedSQL", func() {
		DescribeTable("query with filters", func(q string, parts ...string) {
			actualSQL, err := query.ToIndexedSQL(q)
			Expect(err).NotTo(HaveOccurred())

			for _, part := range parts {
				Expect(pretty(actualSQL)).To(ContainSubstring(pretty(part)))
			}

			_, err = client.Exec(actualSQL)
			Expect(err).NotTo(HaveOccurred())
		},
			Entry("nodes", "n", `SELECT rowid AS id, * FROM search s WHERE s.osm_type IN (1)`),
			Entry("ways", "w", `SELECT rowid AS id, * FROM search s WHERE s.osm_type IN (2)`),
			Entry("relation", "r", `SELECT rowid AS id, * FROM search s WHERE s.osm_type IN (3)`),
			Entry("nodes and way", "nw", `s.osm_type IN (1,2)`),
			Entry("way and nodes", "wn", `s.osm_type IN (1,2)`),
			Entry("all explicit", "nwr", `s.osm_type IN (1,2,3)`),
			Entry("all implicit", "*", `s.osm_type IN (1,2,3)`),
		)

		DescribeTable("query with tags", func(q string, parts ...string) {
			actualSQL, err := query.ToIndexedSQL(q)
			Expect(err).NotTo(HaveOccurred())

			for _, part := range parts {
				Expect(pretty(actualSQL)).To(ContainSubstring(pretty(part)))
			}

			_, err = client.Exec(actualSQL)
			Expect(err).NotTo(HaveOccurred())
		},
			Entry("single tag", "n[amenity=restaurant]", `SELECT rowid AS id, * FROM search s WHERE s.osm_type IN (1) AND ( s.tags->>'$.amenity' = 'restaurant' ) AND s.tags MATCH '( "amenity" AND ( "restaurant" ) )'`),
			Entry("all tags", `nrw[*="*King*","*Queen*"]`, `s.osm_type IN (1,2,3)`, `s.tags MATCH '( "*King*" OR "*Queen*" )'`),
			Entry("all tags with negative", `n[*="cafe"][*!="Starbucks"]`, `s.tags MATCH '( "cafe" ) NOT ( "Starbucks" )'`),
			Entry("partial match", "nw[name!~Starbucks][name=~coffee]", `NOT ( "name" AND ( "Starbucks" ) )`, `( "name" AND ( "coffee" ) )`, `( LOWER(s.tags->>'$.name') NOT GLOB '*starbucks*' )`, `( LOWER(s.tags->>'$.name') GLOB '*coffee*' )`),
			Entry("partial match", `nw[name=~"Coffee Cafe*"]`, `( "name" AND ( "Coffee Cafe"* ) )`, `( LOWER(s.tags->>'$.name') GLOB '*coffee cafe**' )`),
			Entry("multiple tags", "n[amenity=restaurant][cuisine=sushi]", `( "amenity" AND ( "restaurant" ) )`, `( "cuisine" AND ( "sushi" ) )`),
			Entry("single tag with multiple values", "nw[amenity=restaurant,pub,cafe]", `( "amenity" AND ( "restaurant" OR "pub" OR "cafe" ) )`),
			Entry("single tag with multiple values with quotes", `nw[amenity="restaurant","pub","cafe"]`, `( "amenity" AND ( "restaurant" OR "pub" OR "cafe" ) )`),
			Entry("single tag with multiple values with and without quotes", `nw[amenity=Bobs Burgers,"Starbucks"]`, `( "amenity" AND ( "Bobs Burgers" OR "Starbucks" ) )`),
			Entry("single tag that exists", "nw[name]", `( "name" )`, `( s.tags->>'$.name' IS NOT NULL )`),
			Entry("multiple tag that exists", "r[route][ref][network]", `( "route" ) AND ( "ref" ) AND ( "network" )`),
			Entry("multiple tag that have value and exist", "r[amenity=restaurant][name]", `( "amenity" AND ( "restaurant" ) )`, `( "name" )`),
			Entry("tag with not matcher", "nw[amenity=coffee][name!=Starbucks]", `( "amenity" AND ( "coffee" ) ) NOT ( "name" AND ( "Starbucks" ) )`),
			Entry("tag should not exist", "nw[amenity=coffee][!name]", `( "amenity" AND ( "coffee" ) ) NOT ( "name" )`),
			Entry("everything", `nrw[name][!amenity][name="*King*","*Queen*"]`, `( "name" )`, `( "name" AND ( "*King*" OR "*Queen*" ) )`, `NOT ( "amenity" )`, `( s.tags->>'$.amenity' IS NULL )`, `( s.tags->>'$.name' IS NOT NULL )`),
			Entry("with table area", "n[amenity=restaurant](area=test)", `test_search`, `( "amenity" AND ( "restaurant" ) )`, `( s.tags->>'$.amenity' = 'restaurant' )`),
			Entry("with ids", "n(id=1,123,4567)", `s.osm_id IN (1,123,4567)`),
			Entry("with bounding box (bb=minLon,minLat,maxLon,maxLat)", "n(bb=1.10,2.20,11.11,99.99)", `1.10 <= s.minLon`, `2.20 <= s.minLat`, `s.maxLon <= 11.11`, `s.maxLat <= 99.99`),
			Entry("with greater than", "n[pop>100]", `s.tags->>'$.pop' > 100`, `( "pop" )`),
			Entry("with greater than equal", "n[pop>=100]", `s.tags->>'$.pop' >= 100`, `( "pop" )`),
			Entry("with less than", "n[pop<100]", `s.tags->>'$.pop' < 100`, `( "pop" )`),
			Entry("with less than equal", "n[pop<=100]", `s.tags->>'$.pop' <= 100`, `( "pop" )`),
		)
	})
})
