package runtime

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type Point orb.Point

func (r Point) AsBound() Bound {
	return Bound{orb.Point(r).Bound()}
}

func (r Point) AsFeature(properties map[string]interface{}) *geojson.Feature {
	feature := geojson.NewFeature(orb.Point(r))

	for name, value := range properties {
		feature.Properties[name] = value
	}

	return feature
}
