package runtime

import (
	"math"

	"github.com/jtarchie/knowhere/query"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type WrappedResult struct {
	query.Result
}

type WrappedBound struct {
	orb.Bound
}

func (wr *WrappedResult) Bbox() WrappedBound {
	return WrappedBound{
		orb.Bound{
			Min: orb.Point{wr.MinLon, wr.MinLat},
			Max: orb.Point{wr.MaxLon, wr.MaxLat},
		},
	}
}

func (wr *WrappedResult) AsFeature() *geojson.Feature {
	feature := geojson.NewFeature(wr.Bbox().Center())
	feature.Properties["title"] = wr.Name

	return feature
}

func (wr *WrappedBound) Intersects(bounds WrappedBound) bool {
	return wr.Bound.Intersects(bounds.Bound)
}

// Extends a bounding box in kilometers in each direction.
// This is for best effort, not exact.
func (wb *WrappedBound) Extend(radius float64) WrappedBound {
	bounds := orb.Bound{}
	kmInDegreesLat := 1 / 111.0 // 1 degree in km
	avgLat := math.Cos(bounds.Min[1] * math.Pi / 180)
	kmInDegreesLon := kmInDegreesLat / avgLat

	deltaLat := radius * kmInDegreesLat
	deltaLon := radius * kmInDegreesLon

	bounds.Min[0] = wb.Min[0] - deltaLon
	bounds.Max[0] = wb.Max[0] + deltaLon
	bounds.Min[1] = wb.Min[1] - deltaLat
	bounds.Max[1] = wb.Max[1] + deltaLat

	return WrappedBound{bounds}
}
