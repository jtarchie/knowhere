package query_test

import (
	"github.com/jtarchie/knowhere/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// examples taken from:
// https://docs.geodesk.com/java/goql

var _ = Describe("Building a query", func() {
	DescribeTable("can parse types into AST", func(q string, types ...query.FilterType) {
		result, err := query.Parse(q)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      types,
			Tags:       []query.FilterTag{},
		}))
	},
		Entry("nodes", "n", query.NodeFilter),
		Entry("ways", "w", query.WayFilter),
		Entry("area", "r", query.RelationFilter),
		Entry("nodes and ways", "nw", query.NodeFilter, query.WayFilter),
		Entry("ways and nodes", "wn", query.NodeFilter, query.WayFilter),
		Entry("ways and relations", "wr", query.WayFilter, query.RelationFilter),
		Entry("all explicit", "nwr", query.NodeFilter, query.WayFilter, query.RelationFilter),
		Entry("all implicit", "*", query.NodeFilter, query.WayFilter, query.RelationFilter),
		Entry("duplicate ways and nodes", "wwnn", query.NodeFilter, query.WayFilter),
	)

	It("errors with unrecognized type", func() {
		_, err := query.Parse("not")
		Expect(err).To(HaveOccurred())
	})

	It("can parse a single tag", func() {
		ast, err := query.Parse("n[amenity=restaurant]")
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.NodeFilter},
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
		ast, err := query.Parse("nw[amenity=restaurant][cuisine=sushi]")
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types: []query.FilterType{
				query.NodeFilter,
				query.WayFilter,
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
		ast, err := query.Parse("nw[amenity=restaurant][cuisine=sushi][takeaway][website]")
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types: []query.FilterType{
				query.NodeFilter,
				query.WayFilter,
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
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "highway",
					Lookups: []string{"residential"},
					Op:      query.OpEquals,
				},
				{
					Name:    "oneway",
					Lookups: []string{},
					Op:      query.OpNotExists,
				},
			},
		}))
	})

	It("can support multiple value on a single tag", func() {
		ast, err := query.Parse("w[amenity=restaurant,pub,cafe]")
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "amenity",
					Lookups: []string{"restaurant", "pub", "cafe"},
					Op:      query.OpEquals,
				},
			},
		}))
	})

	It("errors on unbalanced brackets for tags", func() {
		_, err := query.Parse("w[amenity=restaurant,pub,cafe")
		Expect(err).To(HaveOccurred())
	})

	It("can support tags that do not equal values", func() {
		ast, err := query.Parse("w[highway][highway!=motorway,primary]")
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "highway",
					Lookups: []string{},
					Op:      query.OpExists,
				},
				{
					Name:    "highway",
					Lookups: []string{"motorway", "primary"},
					Op:      query.OpNotEquals,
				},
			},
		}))
	})

	It("supports any feature", func() {
		ast, err := query.Parse("*[!name]")
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.NodeFilter, query.WayFilter, query.RelationFilter},
			Tags: []query.FilterTag{
				{
					Name:    "name",
					Lookups: []string{},
					Op:      query.OpNotExists,
				},
			},
		}))
	})

	It("can support quoted values for a tag", func() {
		ast, err := query.Parse(`nw[amenity=pub][name="The King's Head"]`)
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.NodeFilter, query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "amenity",
					Lookups: []string{"pub"},
					Op:      query.OpEquals,
				},
				{
					Name:    "name",
					Lookups: []string{"The King's Head"},
					Op:      query.OpEquals,
				},
			},
		}))

		ast, err = query.Parse(`nw[amenity=pub][name="The King's Head","Another Value",Yep]`)
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.NodeFilter, query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "amenity",
					Lookups: []string{"pub"},
					Op:      query.OpEquals,
				},
				{
					Name:    "name",
					Lookups: []string{"The King's Head", "Another Value", "Yep"},
					Op:      query.OpEquals,
				},
			},
		}))
	})

	It("supports options for a globbing syntax", func() {
		ast, err := query.Parse(`nw[amenity=pub][name="*King*"]`)
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.NodeFilter, query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "amenity",
					Lookups: []string{"pub"},
					Op:      query.OpEquals,
				},
				{
					Name:    "name",
					Lookups: []string{"*King*"},
					Op:      query.OpEquals,
				},
			},
		}))
	})

	It("supports directives for the query", func() {
		ast, err := query.Parse(`nw[amenity=pub](area=colorado)`)
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Types: []query.FilterType{query.NodeFilter, query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "amenity",
					Lookups: []string{"pub"},
					Op:      query.OpEquals,
				},
			},
			Directives: map[string]query.FilterDirective{
				"area": []string{"colorado"},
			},
		}))
	})

	It("supports inequalities", func() {
		ast, err := query.Parse(`nw[pop>0]`)
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.NodeFilter, query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "pop",
					Lookups: []string{"0"},
					Op:      query.OpGreaterThan,
				},
			},
		}))

		ast, err = query.Parse(`nw[pop>=0]`)
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.NodeFilter, query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "pop",
					Lookups: []string{"0"},
					Op:      query.OpGreaterThanEquals,
				},
			},
		}))

		ast, err = query.Parse(`nw[pop<0]`)
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.NodeFilter, query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "pop",
					Lookups: []string{"0"},
					Op:      query.OpLessThan,
				},
			},
		}))

		ast, err = query.Parse(`nw[pop<=0]`)
		Expect(err).NotTo(HaveOccurred())
		Expect(ast).To(Equal(&query.AST{
			Directives: map[string]query.FilterDirective{},
			Types:      []query.FilterType{query.NodeFilter, query.WayFilter},
			Tags: []query.FilterTag{
				{
					Name:    "pop",
					Lookups: []string{"0"},
					Op:      query.OpLessThanEquals,
				},
			},
		}))
	})
})
