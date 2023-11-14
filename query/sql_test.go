package query_test

import (
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/jtarchie/knowhere/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Build SQL from a query", func() {
	pretty := func(sql string) string {
		parsed, err := parser.ParseOne(sql)
		Expect(err).NotTo(HaveOccurred())

		f := tree.DefaultPrettyCfg()

		return f.Pretty(parsed.AST)
	}

	DescribeTable("query with filters", func(q string, expectedSQL string) {
		actualSQL, err := query.ToSQL(q)
		Expect(err).NotTo(HaveOccurred())

		Expect(pretty(actualSQL)).To(Equal(pretty(expectedSQL)))
	},
		Entry("nodes", "n", `SELECT * FROM entries AS e WHERE (e.osm_type = 'node')`),
		Entry("ways", "w", `SELECT * FROM entries AS e WHERE (e.osm_type = 'way')`),
		Entry("area", "a", `SELECT * FROM entries AS e WHERE (e.osm_type = 'area')`),
		Entry("relation", "r", `SELECT * FROM entries AS e WHERE (e.osm_type = 'relation')`),
		Entry("nodes and area", "na", `SELECT * FROM entries AS e WHERE (e.osm_type = 'node' OR e.osm_type = 'area')`),
		Entry("area and nodes", "an", `SELECT * FROM entries AS e WHERE (e.osm_type = 'node' OR e.osm_type = 'area')`),
		Entry("all explicit", "nwar", `SELECT * FROM entries AS e WHERE (e.osm_type = 'node' OR e.osm_type = 'area' OR e.osm_type = 'way' OR e.osm_type = 'relation')`),
		Entry("all implicit", "*", `SELECT * FROM entries AS e WHERE (e.osm_type = 'node' OR e.osm_type = 'area' OR e.osm_type = 'way' OR e.osm_type = 'relation')`),
	)
})
