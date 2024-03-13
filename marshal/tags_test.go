package marshal_test

import (
	"github.com/jtarchie/knowhere/marshal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tags", func() {
	It("returns an empty object when there are no tags", func() {
		payload := marshal.Tags(nil, nil)
		Expect(payload).To(MatchJSON(`{}`))

		payload = marshal.Tags(map[string]string{}, nil)
		Expect(payload).To(MatchJSON(`{}`))
	})

	It("returns key-value pairs as JSON", func() {
		payload := marshal.Tags(map[string]string{
			"a": "1",
		}, nil)
		Expect(payload).To(MatchJSON(`{"a":"1"}`))

		payload = marshal.Tags(map[string]string{
			"a": "1",
			"2": "b",
		}, nil)
		Expect(payload).To(MatchJSON(`{"a":"1", "2": "b"}`))

		payload = marshal.Tags(map[string]string{
			"a": "1",
			"2": "b",
			"c": "3",
		}, nil)
		Expect(payload).To(MatchJSON(`{"a":"1", "2": "b", "c": "3"}`))
	})

	When("allowed tags are provided", func() {
		It("only allows those tags", func() {
			payload = marshal.Tags(map[string]string{
				"a": "1",
				"2": "b",
				"c": "3",
			}, map[string]struct{}{
				"a": {},
			})
			Expect(payload).To(MatchJSON(`{"a":"1"}`))
		})
	})
})
