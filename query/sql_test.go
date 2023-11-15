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

	DescribeTable("query with filters", func(q string, expectedSQL string) {
		actualSQL, err := query.ToSQL(q)
		Expect(err).NotTo(HaveOccurred())

		Expect(pretty(actualSQL)).To(Equal(pretty(expectedSQL)))
	},
		Entry("nodes", "n", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node'`),
		Entry("ways", "w", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'way'`),
		Entry("area", "a", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'area'`),
		Entry("relation", "r", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'relation'`),
		Entry("nodes and area", "na", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node OR area'`),
		Entry("area and nodes", "an", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node OR area'`),
		Entry("all explicit", "nwar", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node OR area OR way OR relation'`),
		Entry("all implicit", "*", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node OR area OR way OR relation'`),
	)

	DescribeTable("query with tags", func(q string, expectedSQL string) {
		actualSQL, err := query.ToSQL(q)
		Expect(err).NotTo(HaveOccurred())

		Expect(pretty(actualSQL)).To(Equal(pretty(expectedSQL)))
	},
		Entry("single tag", "n[amenity=restaurant]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node' AND s.tags MATCH '( ("amenity restaurant") )'`),
		Entry("multiple tags", "n[amenity=restaurant][cuisine=sushi]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node' AND s.tags MATCH '( ("amenity restaurant") ) AND ( ("cuisine sushi") )'`),
		Entry("single tag with multiple values", "na[amenity=restaurant,pub,cafe]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node OR area' AND s.tags MATCH '( ("amenity restaurant") OR ("amenity pub") OR ("amenity cafe") )'`),
		Entry("single tag that exists", "na[amenity]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node OR area' AND s.tags MATCH '( "amenity" )'`),
		Entry("multiple tag that exists", "r[route][ref][network]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'relation' AND s.tags MATCH '( "route" ) AND ( "ref" ) AND ( "network" )'`),
		Entry("multiple tag that have value and exist", "r[amenity=restaurant][name]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'relation' AND s.tags MATCH '( ("amenity restaurant") ) AND ( "name" )'`),
		Entry("tag with not matcher", "na[amenity=coffee][name!=Starbucks]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node OR area' AND s.tags MATCH '( ("amenity coffee") )' AND s.tags MATCH NOT '( ("name Starbucks") )'`),
		Entry("tag should not exist", "na[amenity=coffee][!name]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node OR area' AND s.tags MATCH '( ("amenity coffee") )' AND s.tags MATCH NOT '( "name" )'`),
		Entry("everything", `narw[name][!amenity][name="*King*","*Queen*"]`, `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE s.osm_type MATCH 'node OR area OR way OR relation' AND s.tags MATCH '( "name" ) AND ( ("name *King*") OR ("name *Queen*") )' AND s.tags MATCH NOT '( "amenity" )'`),
	)
})
