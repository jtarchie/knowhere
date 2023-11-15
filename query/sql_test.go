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
		Entry("nodes", "n", `SELECT * FROM entries e WHERE (e.osm_type = 'node')`),
		Entry("ways", "w", `SELECT * FROM entries e WHERE (e.osm_type = 'way')`),
		Entry("area", "a", `SELECT * FROM entries e WHERE (e.osm_type = 'area')`),
		Entry("relation", "r", `SELECT * FROM entries e WHERE (e.osm_type = 'relation')`),
		Entry("nodes and area", "na", `SELECT * FROM entries e WHERE (e.osm_type = 'node' OR e.osm_type = 'area')`),
		Entry("area and nodes", "an", `SELECT * FROM entries e WHERE (e.osm_type = 'node' OR e.osm_type = 'area')`),
		Entry("all explicit", "nwar", `SELECT * FROM entries e WHERE (e.osm_type = 'node' OR e.osm_type = 'area' OR e.osm_type = 'way' OR e.osm_type = 'relation')`),
		Entry("all implicit", "*", `SELECT * FROM entries e WHERE (e.osm_type = 'node' OR e.osm_type = 'area' OR e.osm_type = 'way' OR e.osm_type = 'relation')`),
	)

	DescribeTable("query with tags", func(q string, expectedSQL string) {
		actualSQL, err := query.ToSQL(q)
		Expect(err).NotTo(HaveOccurred())

		Expect(pretty(actualSQL)).To(Equal(pretty(expectedSQL)))
	},
		Entry("single tag", "n[amenity=restaurant]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE (e.osm_type = 'node') AND s.tags MATCH '( ("amenity restaurant") )'`),
		Entry("multiple tags", "n[amenity=restaurant][cuisine=sushi]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE (e.osm_type = 'node') AND s.tags MATCH '( ("amenity restaurant") ) AND ( ("cuisine sushi") )'`),
		Entry("single tag with multiple values", "na[amenity=restaurant,pub,cafe]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE (e.osm_type = 'node' OR e.osm_type = 'area') AND s.tags MATCH '( ("amenity restaurant") OR ("amenity pub") OR ("amenity cafe") )'`),
		Entry("single tag with multiple values", "na[amenity=restaurant,pub,cafe]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE (e.osm_type = 'node' OR e.osm_type = 'area') AND s.tags MATCH '( ("amenity restaurant") OR ("amenity pub") OR ("amenity cafe") )'`),
		Entry("single tag that exists", "na[amenity]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE (e.osm_type = 'node' OR e.osm_type = 'area') AND s.tags MATCH '( "amenity" )'`),
		Entry("multiple tag that exists", "r[route][ref][network]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE (e.osm_type = 'relation') AND s.tags MATCH '( "route" ) AND ( "ref" ) AND ( "network" )'`),
		Entry("multiple tag that have value and exist", "r[amenity=restaurant][name]", `SELECT * FROM entries e JOIN search s ON s.rowid = e.id WHERE (e.osm_type = 'relation') AND s.tags MATCH '( ("amenity restaurant") ) AND ( "name" )'`),
	)
})
