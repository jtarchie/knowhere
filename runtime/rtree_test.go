package runtime_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"

	"github.com/jtarchie/knowhere/query"
	"github.com/jtarchie/knowhere/runtime"
)

var _ = Describe("RTree", func() {
	It("can lookup if something within bounds", func() {
		tree := &runtime.RTree{}

		bounds := &runtime.WrappedBound{
			orb.Bound{
				Min: [2]float64{0, 0},
				Max: [2]float64{100, 100},
			},
		}

		Expect(tree.Within(bounds)).To(BeFalse())

		tree.Insert(
			&runtime.WrappedBound{
				orb.Bound{
					Min: [2]float64{25, 25},
					Max: [2]float64{50, 50},
				},
			},
			&runtime.WrappedResult{},
		)

		Expect(tree.Within(bounds)).To(BeTrue())
	})

	It("returns nearby items", func() {
		tree := &runtime.RTree{}

		actual := []*runtime.WrappedResult{
			&runtime.WrappedResult{query.Result{Name: "1"}},
			&runtime.WrappedResult{query.Result{Name: "2"}},
		}

		tree.Insert(
			&runtime.WrappedBound{
				orb.Bound{
					Min: [2]float64{0, 0},
					Max: [2]float64{25, 525},
				},
			},
			actual[0],
		)

		tree.Insert(
			&runtime.WrappedBound{
				orb.Bound{
					Min: [2]float64{25, 0},
					Max: [2]float64{50, 25},
				},
			},
			actual[1],
		)

		results := tree.Nearby(
			&runtime.WrappedBound{
				orb.Bound{
					Min: [2]float64{10, 10},
					Max: [2]float64{20, 20},
				},
			},
			1,
		)

		Expect(results).To(HaveLen(1))
		Expect(results).To(ContainElements(actual[0]))
	})
})
