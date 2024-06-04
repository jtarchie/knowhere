package runtime

import (
	"github.com/jtarchie/knowhere/query"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type WrappedResult struct {
	query.Result
}

func (wr *WrappedResult) Bbox() *WrappedBound {
	return &WrappedBound{
		orb.Bound{
			Min: orb.Point{wr.MinLon, wr.MinLat},
			Max: orb.Point{wr.MaxLon, wr.MaxLat},
		},
	}
}

func (wr *WrappedResult) AsFeature(properties map[string]interface{}) *geojson.Feature {
	feature := geojson.NewFeature(wr.Bbox().Center())

	feature.Properties["title"] = wr.Name
	for name, value := range properties {
		feature.Properties[name] = value
	}

	return feature
}
