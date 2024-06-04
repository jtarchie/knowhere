package runtime

import (
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/clip"
	"github.com/paulmach/orb/geo"
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

func (r *RTree) Nearby(bound *WrappedBound, count uint) []*WrappedResult {
	results := make([]*WrappedResult, 0, count)

	r.RTreeG.Nearby(
		// rtree.BoxDist[float64, *WrappedResult](bound.Min, bound.Max, nil),
		fromWrappedDistOverlap(bound),
		func(min, max [2]float64, data *WrappedResult, dist float64) bool {
			results = append(results, data)

			count--

			return 0 < count
		},
	)

	return results
}

func fromWrappedDistOverlap(target *WrappedBound) func(min, max [2]float64, data *WrappedResult, item bool) float64 {
	return func(bMin, bMax [2]float64, _ *WrappedResult, _ bool) float64 {
		current := orb.Bound{
			Min: bMin,
			Max: bMax,
		}
		if target.Contains(current.LeftTop()) && target.Contains(current.RightBottom()) {
			return geo.Distance(target.Center(), current.Center())
		}

		clipped := clip.Bound(target.Bound, current)
		if clipped.IsEmpty() {
			return math.MaxFloat64
		}

		return geo.Distance(target.Center(), clipped.Center())
	}
}
