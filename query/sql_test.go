package query_test

import (
	"regexp"
	"strings"

	"github.com/jtarchie/knowhere/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Build SQL from a query", func() {
	pretty := func(sql string) string {
		space := regexp.MustCompile(`\s+`)

		return strings.TrimSpace(space.ReplaceAllString(sql, " "))
	}

	Describe("ToExactSQL", func() {
		DescribeTable("query with filters", func(q string, expectedSQL string) {
			actualSQL, err := query.ToExactSQL(q)
			Expect(err).NotTo(HaveOccurred())

			Expect(pretty(actualSQL)).To(Equal(pretty(expectedSQL)))
		},
			Entry("nodes", "n", `SELECT * FROM entries e WHERE ( e.osm_type = 1 )`),
			Entry("ways", "w", `SELECT * FROM entries e WHERE ( e.osm_type = 2 )`),
			Entry("relation", "r", `SELECT * FROM entries e WHERE ( e.osm_type = 3 )`),
			Entry("nodes and way", "nw", `SELECT * FROM entries e WHERE ( e.osm_type = 1 OR e.osm_type = 2 )`),
			Entry("way and nodes", "wn", `SELECT * FROM entries e WHERE ( e.osm_type = 1 OR e.osm_type = 2 )`),
			Entry("all explicit", "nwr", `SELECT * FROM entries e WHERE ( e.osm_type = 1 OR e.osm_type = 2 OR e.osm_type = 3 )`),
			Entry("all implicit", "*", `SELECT * FROM entries e WHERE ( e.osm_type = 1  OR e.osm_type = 2 OR e.osm_type = 3 )`),
		)

		DescribeTable("query with tags", func(q string, expectedSQL string) {
			actualSQL, err := query.ToExactSQL(q)
			Expect(err).NotTo(HaveOccurred())

			Expect(pretty(actualSQL)).To(Equal(pretty(expectedSQL)))
		},
			Entry("single tag", "n[amenity=restaurant]", `SELECT * FROM entries e WHERE ( e.osm_type = 1 ) AND ( e.tags->>'$.amenity' GLOB 'restaurant' )`),
			Entry("all tags", `nrw[*="*King*","*Queen*"]`, `SELECT * FROM entries e WHERE ( e.osm_type = 1 OR e.osm_type = 2 OR e.osm_type = 3 ) AND ( e.tags GLOB '*King*' OR e.tags GLOB '*Queen*' )`),
			Entry("all tags with negative", `n[*="cafe"][*!="Starbucks"]`, `SELECT * FROM entries e WHERE ( e.osm_type = 1 ) AND ( e.tags GLOB 'cafe' ) AND NOT ( ( e.tags GLOB 'Starbucks' ) )`),
			Entry("multiple tags", "n[amenity=restaurant][cuisine=sushi]", `SELECT * FROM entries e WHERE ( e.osm_type = 1 ) AND ( e.tags->>'$.amenity' GLOB 'restaurant' ) AND ( e.tags->>'$.cuisine' GLOB 'sushi' )`),
			Entry("single tag with multiple values", "nw[amenity=restaurant,pub,cafe]", `SELECT * FROM entries e WHERE ( e.osm_type = 1 OR e.osm_type = 2 ) AND ( e.tags->>'$.amenity' GLOB 'restaurant' OR e.tags->>'$.amenity' GLOB 'pub' OR e.tags->>'$.amenity' GLOB 'cafe' )`),
			Entry("single tag that exists", "nw[amenity]", `SELECT * FROM entries e WHERE ( e.osm_type = 1 OR e.osm_type = 2 ) AND ( e.tags->>'$.amenity' IS NOT NULL )`),
			Entry("multiple tag that exists", "r[route][ref][network]", `SELECT * FROM entries e WHERE ( e.osm_type = 3 ) AND ( e.tags->>'$.route' IS NOT NULL ) AND ( e.tags->>'$.ref' IS NOT NULL ) AND ( e.tags->>'$.network' IS NOT NULL )`),
			Entry("multiple tag that have value and exist", "r[amenity=restaurant][name]", `SELECT * FROM entries e WHERE ( e.osm_type = 3 ) AND ( e.tags->>'$.amenity' GLOB 'restaurant' ) AND ( e.tags->>'$.name' IS NOT NULL )`),
			Entry("tag with not matcher", "nw[amenity=coffee][name!=Starbucks]", `SELECT * FROM entries e WHERE ( e.osm_type = 1 OR e.osm_type = 2 ) AND ( e.tags->>'$.amenity' GLOB 'coffee' ) AND NOT ( ( e.tags->>'$.name' GLOB 'Starbucks' ) )`),
			Entry("tag should not exist", "nw[amenity=coffee][!name]", `SELECT * FROM entries e WHERE ( e.osm_type = 1 OR e.osm_type = 2 ) AND ( e.tags->>'$.amenity' GLOB 'coffee' ) AND NOT ( ( e.tags->>'$.name' IS NOT NULL ) )`),
			Entry("everything", `nrw[name][!amenity][name="*King*","*Queen*"]`, `SELECT * FROM entries e WHERE ( e.osm_type = 1 OR e.osm_type = 2 OR e.osm_type = 3 ) AND ( e.tags->>'$.name' IS NOT NULL ) AND ( e.tags->>'$.name' GLOB '*King*' OR e.tags->>'$.name' GLOB '*Queen*' ) AND NOT ( ( e.tags->>'$.amenity' IS NOT NULL ) )`),
			Entry("with table prefix", "n[amenity=restaurant](prefix=test)", `SELECT * FROM test_entries e WHERE ( e.osm_type = 1 ) AND ( e.tags->>'$.amenity' GLOB 'restaurant' )`),
			Entry("with ids", "n(id=1,123,4567)", `SELECT * FROM entries e WHERE ( e.osm_type = 1 ) AND e.osm_id IN ( 1, 123, 4567 )`),
		)
	})

	Describe("ToIndexedSQL", func() {
		DescribeTable("query with filters", func(q string, expectedSQL string) {
			actualSQL, err := query.ToIndexedSQL(q)
			Expect(err).NotTo(HaveOccurred())

			Expect(pretty(actualSQL)).To(Equal(pretty(expectedSQL)))
		},
			Entry("nodes", "n", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 )`),
			Entry("ways", "w", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 2 )`),
			Entry("relation", "r", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 3 )`),
			Entry("nodes and way", "nw", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 OR e.osm_type = 2 )`),
			Entry("way and nodes", "wn", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 OR e.osm_type = 2 )`),
			Entry("all explicit", "nwr", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 OR e.osm_type = 2 OR e.osm_type = 3 )`),
			Entry("all implicit", "*", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 OR e.osm_type = 2 OR e.osm_type = 3 )`),
		)

		DescribeTable("query with tags", func(q string, expectedSQL string) {
			actualSQL, err := query.ToIndexedSQL(q)
			Expect(err).NotTo(HaveOccurred())

			Expect(pretty(actualSQL)).To(Equal(pretty(expectedSQL)))
		},
			Entry("single tag", "n[amenity=restaurant]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 ) AND s.tags MATCH '( ("amenity restaurant") )'`),
			Entry("all tags", `nrw[*="*King*","*Queen*"]`, `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 OR e.osm_type = 2 OR e.osm_type = 3 ) AND s.tags MATCH '( ("*King*") OR ("*Queen*") )'`),
			Entry("all tags with negative", `n[*="cafe"][*!="Starbucks"]`, `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 ) AND s.tags MATCH '( ("cafe") ) NOT ( ("Starbucks") )'`),
			Entry("multiple tags", "n[amenity=restaurant][cuisine=sushi]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 ) AND s.tags MATCH '( ("amenity restaurant") ) AND ( ("cuisine sushi") )'`),
			Entry("single tag with multiple values", "nw[amenity=restaurant,pub,cafe]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 OR e.osm_type = 2 ) AND s.tags MATCH '( ("amenity restaurant") OR ("amenity pub") OR ("amenity cafe") )'`),
			Entry("single tag that exists", "nw[amenity]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 OR e.osm_type = 2 ) AND s.tags MATCH '( "amenity" )'`),
			Entry("multiple tag that exists", "r[route][ref][network]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 3 ) AND s.tags MATCH '( "route" ) AND ( "ref" ) AND ( "network" )'`),
			Entry("multiple tag that have value and exist", "r[amenity=restaurant][name]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 3 ) AND s.tags MATCH '( ("amenity restaurant") ) AND ( "name" )'`),
			Entry("tag with not matcher", "nw[amenity=coffee][name!=Starbucks]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 OR e.osm_type = 2 ) AND s.tags MATCH '( ("amenity coffee") ) NOT ( ("name Starbucks") )'`),
			Entry("tag should not exist", "nw[amenity=coffee][!name]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 OR e.osm_type = 2 ) AND s.tags MATCH '( ("amenity coffee") ) NOT ( "name" )'`),
			Entry("everything", `nrw[name][!amenity][name="*King*","*Queen*"]`, `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 OR e.osm_type = 2 OR e.osm_type = 3 ) AND s.tags MATCH '( "name" ) AND ( ("name *King*") OR ("name *Queen*") ) NOT ( "amenity" )'`),
			Entry("with table prefix", "n[amenity=restaurant](prefix=test)", `SELECT * FROM test_entries e JOIN test_search s ON s.rowid = e.id WHERE ( e.osm_type = 1 ) AND s.tags MATCH '( ("amenity restaurant") )'`),
			Entry("with ids", "n(id=1,123,4567)", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE ( e.osm_type = 1 ) AND e.osm_id IN ( 1, 123, 4567 )`),
		)
	})
})
