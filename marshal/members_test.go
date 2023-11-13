package marshal_test

import (
	"github.com/jtarchie/knowhere/marshal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/osm"
)

var _ = Describe("Members", func() {
	It("returns an empty array when there are no Members", func() {
		payload := marshal.Members(nil)
		Expect(payload).To(MatchJSON(`[]`))

		payload = marshal.Members(osm.Members{})
		Expect(payload).To(MatchJSON(`[]`))
	})

	It("returns key-value pairs as JSON", func() {
		payload := marshal.Members(osm.Members{
			osm.Member{
				Type: osm.TypeNode,
				Ref:  1,
				Role: "donotcare",
			},
		})
		Expect(payload).To(MatchJSON(`[
			[1,"node","donotcare"]
		]`))

		payload = marshal.Members(osm.Members{
			osm.Member{
				Type: osm.TypeNode,
				Ref:  1,
				Role: "donotcare",
			},
			osm.Member{
				Type: osm.TypeWay,
				Ref:  2,
				Role: "donotcare",
			},
		})
		Expect(payload).To(MatchJSON(`[
			[1,"node","donotcare"],
			[2,"way","donotcare"]
		]`))

		payload = marshal.Members(osm.Members{
			osm.Member{
				Type: osm.TypeNode,
				Ref:  1,
				Role: "donotcare",
			},
			osm.Member{
				Type: osm.TypeBounds,
				Ref:  2,
				Role: "donotcare",
			},
			osm.Member{
				Type: osm.TypeWay,
				Ref:  2,
				Role: "donotcare",
			},
		})
		Expect(payload).To(MatchJSON(`[
			[1,"node","donotcare"],
			[2,"way","donotcare"]
		]`))

		payload = marshal.Members(osm.Members{
			osm.Member{
				Type: osm.TypeNode,
				Ref:  1,
				Role: "donotcare",
			},
			osm.Member{
				Type: osm.TypeWay,
				Ref:  2,
				Role: "donotcare",
			},
			osm.Member{
				Type: osm.TypeBounds,
				Ref:  2,
				Role: "donotcare",
			},
		})
		Expect(payload).To(MatchJSON(`[
			[1,"node","donotcare"],
			[2,"way","donotcare"]
		]`))

		payload = marshal.Members(osm.Members{
			osm.Member{
				Type: osm.TypeBounds,
				Ref:  2,
				Role: "donotcare",
			},
			osm.Member{
				Type: osm.TypeNode,
				Ref:  1,
				Role: "donotcare",
			},
			osm.Member{
				Type: osm.TypeWay,
				Ref:  2,
				Role: "donotcare",
			},
		})
		Expect(payload).To(MatchJSON(`[
			[1,"node","donotcare"],
			[2,"way","donotcare"]
		]`))
	})
})
