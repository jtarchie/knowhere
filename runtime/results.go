package runtime

import (
	"slices"

	"github.com/paulmach/orb/geo"
)

type Results []Result

func (r Results) Cluster(radius float64) Results {
	results := Results{}
	tree := &RTree{}

	slices.SortStableFunc(r, func(a Result, b Result) int {
		areaA := geo.Area(a.Bbox().Bound)
		areaB := geo.Area(b.Bbox().Bound)

		if areaA < areaB {
			return 1
		}
		if areaA > areaB {
			return -1
		}

		return 0
	})

	for _, entry := range r {
		extended := entry.Bbox().Extend(radius)

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
	alreadyUsed := map[Result]struct{}{}
	size++

	for _, result := range r {
		// initially populate with result that is looking for neighbors
		nearby := Results{result}

		if _, ok := alreadyUsed[result]; ok {
			// if there is any crossover with result A and B
			// don't search for anything already
			continue
		}

		extended := result.Bbox().Extend(originRadius)
		tree.Search(extended.Min, extended.Max, func(min, max [2]float64, result Result) bool {
			if _, ok := alreadyUsed[result]; !ok {
				// only find unique neighbors, don't share
				nearby = append(nearby, result)
			}

			return true
		})

		slices.SortStableFunc(nearby, func(a Result, b Result) int {
			return int(geo.Distance(a.Bbox().Center(), result.Bbox().Center()) - geo.Distance(b.Bbox().Center(), result.Bbox().Center()))
		})

		if len(nearby) >= int(size) {
			nearby = nearby[:size]
			results = append(results, nearby)

			for _, used := range nearby {
				alreadyUsed[used] = struct{}{}
			}
		}
	}

	return results
}

func (r Results) AsTree(radius float64) *RTree {
	tree := &RTree{}

	for _, entry := range r {
		bbox := entry.Bbox()
		if 0 < radius {
			bbox = bbox.Extend(radius)
		}

		tree.Insert(bbox, entry)
	}

	return tree
}
