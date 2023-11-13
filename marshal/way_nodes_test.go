package marshal_test

import (
	"github.com/jtarchie/knowhere/marshal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/osm"
)

var _ = Describe("WayNodes", func() {
	It("returns an empty array when there are no WayNodes", func() {
		payload := marshal.WayNodes(nil)
		Expect(payload).To(MatchJSON(`[]`))

		payload = marshal.WayNodes(osm.WayNodes{})
		Expect(payload).To(MatchJSON(`[]`))
	})

	It("returns IDs as an array of JSON", func() {
		payload := marshal.WayNodes(osm.WayNodes{
			osm.WayNode{
				ID: 1,
			},
		})
		Expect(payload).To(MatchJSON(`[1]`))

		payload = marshal.WayNodes(osm.WayNodes{
			osm.WayNode{ID: 1},
			osm.WayNode{ID: 2},
		})
		Expect(payload).To(MatchJSON(`[1, 2]`))

		payload = marshal.WayNodes(osm.WayNodes{
			osm.WayNode{ID: 1},
			osm.WayNode{ID: 2},
			osm.WayNode{ID: 3},
		})
		Expect(payload).To(MatchJSON(`[1, 2, 3]`))
	})
})
