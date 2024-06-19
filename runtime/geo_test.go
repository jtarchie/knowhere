package runtime_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jtarchie/knowhere/runtime"
)

var _ = Describe("Geo", func() {
	It("returns a color", func() {
		geo := runtime.Geo{}

		for i := 0; i < 100; i++ {
			color := geo.Color(i)
			Expect(color).NotTo(Equal(""))
		}
	})

	It("converts a slice of bound(s) to bounds", func() {
		geo := runtime.Geo{}

		bound := []runtime.Bound{
			{},
			{},
		}

		bounds := geo.AsBounds(bound...)
		Expect(bounds).To(HaveLen(2))
		Expect(bounds).To(ContainElements(bound[0], bound[1]))
	})

	It("converts a slice of result(s) to results", func() {
		geo := runtime.Geo{}

		result := []runtime.Result{
			{},
			{},
		}

		results := geo.AsResults(result...)
		Expect(results).To(HaveLen(2))
		Expect(results).To(ContainElements(result[0], result[1]))
	})
})
