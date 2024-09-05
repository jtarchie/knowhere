package runtime

import (
	"github.com/paulmach/orb/geo"
)

type Geo struct{}

func (g *Geo) Rtree() *RTree {
	return &RTree{}
}

func (g *Geo) AsResults(results ...Result) Results {
	return Results(results)
}

func (g *Geo) AsBounds(bounds ...Bound) Bounds {
	return Bounds(bounds)
}

func (g *Geo) AsPoint(lat, lng float64) Point {
	return Point{lng, lat}
}

func (g *Geo) Distance(p1, p2 Bound) float64 {
	return geo.Distance(p1.bound.Center(), p2.bound.Center())
}
