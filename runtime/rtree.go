package runtime

import (
	"github.com/tidwall/rtree"
)

type RTree struct {
	rtree.RTreeG[*WrappedResult]
}

func (r *RTree) Insert(bound *WrappedBound, element *WrappedResult) {
	r.RTreeG.Insert(bound.Min, bound.Max, element)
}

func (r *RTree) Within(bound *WrappedBound) bool {
	contains := false

	r.RTreeG.Search(bound.Min, bound.Max, func(min, max [2]float64, _ *WrappedResult) bool {
		contains = true

		// as long as one thing exists
		return false
	})

	return contains
}
