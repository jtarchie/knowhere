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
		
		// fmt.Printf("a = %f, b = %f\n", areaA, areaB)

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

func (r Results) Overlap(b Results, radius float64, size uint) []Results {
	tree := b.AsTree(0)

	results := []Results{}
	alreadyUsed := map[Result]struct{}{}

	for _, result := range r {
		var nearby Results

		if _, ok := alreadyUsed[result]; ok {
			continue
		}

		extended := result.Bbox().Extend(radius)
		tree.Search(extended.Min, extended.Max, func(min, max [2]float64, result Result) bool {
			nearby = append(nearby, result)

			return int(size) > len(nearby)
		})

		if len(nearby) == int(size) {
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
