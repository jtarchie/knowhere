package marshal_test

import (
	"strings"

	"github.com/jtarchie/knowhere/marshal"
	"github.com/jtarchie/knowhere/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Results", func() {
	It("returns an empty object when there are no tags", func() {
		builder := &strings.Builder{}

		err := marshal.Results(builder, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(builder.String()).To(MatchJSON(`[]`))

		builder.Reset()

		err = marshal.Results(builder, []query.Result{})
		Expect(err).NotTo(HaveOccurred())
		Expect(builder).To(MatchJSON(`[]`))
	})

	It("returns array of JSON", func() {
		builder := &strings.Builder{}
		err := marshal.Results(builder, []query.Result{
			{},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(builder).To(MatchJSON(`[
			{
				"id": 0,
				"minLat": 0,
				"maxLat": 0,
				"minLon": 0,
				"maxLon": 0,
				"type": ""
			}
		]`))
	})
})
