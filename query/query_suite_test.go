package query_test

import (
	"testing"

	"github.com/jtarchie/knowhere/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestQuery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Query Suite")
}

var _ = Describe("Building a query", func() {
	FDescribeTable("can parse types into AST", func(q string, types ...query.FilterType) {
		result, err := query.Parse(q)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(&query.AST{
			Types: types,
			Tags:  []query.FilterTag{},
		}))
	},
		Entry("nodes", "n", query.NodeFilter),
		Entry("ways", "w", query.WayFilter),
		Entry("area", "a", query.AreaFilter),
		Entry("nodes and area", "na", query.NodeFilter, query.AreaFilter),
		Entry("area and nodes", "an", query.NodeFilter, query.AreaFilter),
		Entry("nodes and ways", "nw", query.NodeFilter, query.WayFilter),
		Entry("ways and nodes", "wn", query.NodeFilter, query.WayFilter),
		Entry("ways and relations", "wr", query.WayFilter, query.RelationFilter),
		Entry("all explicit", "nwar", query.NodeFilter, query.AreaFilter, query.WayFilter, query.RelationFilter),
		Entry("all explicit", "*", query.NodeFilter, query.AreaFilter, query.WayFilter, query.RelationFilter),
		Entry("duplicate ways and nodes", "wwnn", query.NodeFilter, query.WayFilter),
	)

	It("errors with unrecognized type", func() {
		_, err := query.Parse("not")
		Expect(err).To(HaveOccurred())
	})

	It("can parse a single tag", func() {
		ast, err := query.Parse("a[amenity=restaurant]")
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Types: []query.FilterType{query.AreaFilter},
			Tags: []query.FilterTag{
				{
					Name:    "amenity",
					Lookups: []string{"restaurant"},
					Op:      query.OpEquals,
				},
			},
		}))
	})

	It("can parse multiple tags", func() {
		ast, err := query.Parse("na[amenity=restaurant][cuisine=sushi]")
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Types: []query.FilterType{
				query.NodeFilter,
				query.AreaFilter,
			},
			Tags: []query.FilterTag{
				{
					Name:    "amenity",
					Lookups: []string{"restaurant"},
					Op:      query.OpEquals,
				},
				{
					Name:    "cuisine",
					Lookups: []string{"sushi"},
					Op:      query.OpEquals,
				},
			},
		}))
	})

	It("can parse tags with values and existence", func() {
		ast, err := query.Parse("na[amenity=restaurant][cuisine=sushi][takeaway][website]")
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Types: []query.FilterType{
				query.NodeFilter,
				query.AreaFilter,
			},
			Tags: []query.FilterTag{
				{
					Name:    "amenity",
					Lookups: []string{"restaurant"},
					Op:      query.OpEquals,
				},
				{
					Name:    "cuisine",
					Lookups: []string{"sushi"},
					Op:      query.OpEquals,
				},
				{
					Name:    "takeaway",
					Lookups: []string{},
					Op:      query.OpExists,
				},
				{
					Name:    "website",
					Lookups: []string{},
					Op:      query.OpExists,
				},
			},
		}))
	})

	It("parses for tags that should not exist", func() {
		ast, err := query.Parse("w[highway=residential][!oneway]")
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Types: []query.FilterType{query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "highway",
					Lookups: []string{"residential"},
					Op:      query.OpEquals,
				},
				{
					Name:    "oneway",
					Lookups: []string{},
					Op:      query.OpNoExists,
				},
			},
		}))
	})

	It("can support multiple value on a single tag", func() {
		ast, err := query.Parse("na[amenity=restaurant,pub,cafe]")
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Types: []query.FilterType{query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "amenity",
					Lookups: []string{"residential", "pub", "cafe"},
					Op:      query.OpEquals,
				},
			},
		}))
	})
})
