package runtime

import (
	"github.com/jtarchie/knowhere/query"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type Result struct {
	query.Result
}

func (wr Result) Bound() Bound {
	return Bound{
		bound: orb.Bound{
			Min: orb.Point{wr.MinLon, wr.MinLat},
			Max: orb.Point{wr.MaxLon, wr.MaxLat},
		},
	}
}

func (wr Result) AsFeature(properties map[string]interface{}) *geojson.Feature {
	feature := geojson.NewFeature(wr.Bound().bound.Center())

	feature.Properties["title"] = wr.Name()
	feature.Properties["id"] = wr.ID
	feature.Properties["type"] = wr.OsmType

	for name, value := range properties {
		feature.Properties[name] = value
	}

	return feature
}
