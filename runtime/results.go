package runtime

import (
	"slices"

	"github.com/paulmach/orb/geo"
)

type Results []Result

func (r Results) Cluster(radius float64) Results {
	results := Results{}
	tree := &RTree{}

	// sort by biggest surface area first
	slices.SortStableFunc(r, func(a Result, b Result) int {
		areaA := geo.Area(a.Bound().bound)
		areaB := geo.Area(b.Bound().bound)

		if areaA < areaB {
			return 1
		}
		if areaA > areaB {
			return -1
		}

		return 0
	})

	for _, entry := range r {
		extended := entry.Bound().Extend(radius)

		if !tree.Within(extended) {
			results = append(results, entry)
			tree.Insert(extended, entry)
		}
	}

	return results
}

func (r Results) Overlap(b Results, originRadius float64, neighborRadius float64, size int) []Results {
	tree := b.AsTree(neighborRadius)

	results := []Results{}
	alreadyUsed := map[int64]struct{}{}
	size++

	for _, result := range r {
		// initially populate with result that is looking for neighbors
		nearby := Results{result}

		if _, ok := alreadyUsed[result.ID]; ok {
			// if there is any crossover with result A and B
			// don't search for anything already
			continue
		}

		extended := result.Bound().Extend(originRadius)
		tree.RTreeG.Search(extended.bound.Min, extended.bound.Max, func(min, max [2]float64, result Result) bool {
			if _, ok := alreadyUsed[result.ID]; !ok {
				// only find unique neighbors, don't share
				nearby = append(nearby, result)
			}

			return true
		})

		slices.SortStableFunc(nearby, func(a Result, b Result) int {
			return int(geo.Distance(a.Bound().bound.Center(), result.Bound().bound.Center()) - geo.Distance(b.Bound().bound.Center(), result.Bound().bound.Center()))
		})

		if len(nearby) >= int(size) {
			nearby = nearby[:size]
			results = append(results, nearby)

			for _, used := range nearby {
				alreadyUsed[used.ID] = struct{}{}
			}
		}
	}

	return results
}

func (r Results) AsTree(radius float64) *RTree {
	tree := &RTree{}

	for _, entry := range r {
		bbox := entry.Bound()
		if 0 < radius {
			bbox = bbox.Extend(radius)
		}

		tree.Insert(bbox, entry)
	}

	return tree
}
