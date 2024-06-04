package runtime

import (
	"github.com/jtarchie/knowhere/query"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type Result struct {
	query.Result
}

func (wr *Result) Bbox() *Bound {
	return &Bound{
		orb.Bound{
			Min: orb.Point{wr.MinLon, wr.MinLat},
			Max: orb.Point{wr.MaxLon, wr.MaxLat},
		},
	}
}

func (wr *Result) AsFeature(properties map[string]interface{}) *geojson.Feature {
	feature := geojson.NewFeature(wr.Bbox().Center())

	feature.Properties["title"] = wr.Name
	for name, value := range properties {
		feature.Properties[name] = value
	}

	return feature
}

type Results []Result

func (r Results) Cluster(radius float64) Results {
	results := Results{}
	tree := &RTree{}

	for _, entry := range r {
		extended := entry.Bbox().Extend(radius)

		if !tree.Within(extended) {
			results = append(results, entry)
			tree.Insert(extended, &entry)
		}
	}

	return results
}

func (r Results) AsTree(radius float64) *RTree {
	tree := &RTree{}

	for _, entry := range r {
		extended := entry.Bbox().Extend(radius)

		tree.Insert(extended, &entry)
	}

	return tree
}
