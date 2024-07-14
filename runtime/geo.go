package runtime

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
