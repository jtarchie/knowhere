package services_test

import (
	"github.com/jtarchie/knowhere/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/osm"
)

var _ = Describe("When importing an OSM file", func() {
	It("loads nodes, ways, and relations", func() {
		importer := services.NewImporter("../fixtures/sample.pbf")
		var nodeCount, wayCount, relationCount uint64

		err := importer.Execute(
			func(_ *osm.Node) error {
				nodeCount++

				return nil
			},
			func(_ *osm.Way) error {
				wayCount++

				return nil
			},
			func(_ *osm.Relation) error {
				relationCount++

				return nil
			},
		)
		Expect(err).NotTo(HaveOccurred())

		Expect(nodeCount).To(BeEquivalentTo(290))
		Expect(wayCount).To(BeEquivalentTo(44))
		Expect(relationCount).To(BeEquivalentTo(5))
	})
})
