package runtime

import (
	"fmt"

	"github.com/engelsjk/polygol"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
	"github.com/samber/lo"
)

type Bound struct {
	bound orb.Bound
}

func NewBound(bound orb.Bound) Bound {
	return Bound{bound: bound}
}

func (wb Bound) AsBB() string {
	return fmt.Sprintf("%f,%f,%f,%f", wb.bound.Min[0], wb.bound.Min[1], wb.bound.Max[0], wb.bound.Max[1])
}

func (wb Bound) Min() [2]float64 {
	return wb.bound.Min
}

func (wb Bound) Max() [2]float64 {
	return wb.bound.Max
}

func (wb Bound) Center() Point {
	return Point(wb.bound.Center())
}

func (wb Bound) Left() float64 {
	return wb.bound.Left()
}

func (wb Bound) Right() float64 {
	return wb.bound.Right()
}

func (wb Bound) Top() float64 {
	return wb.bound.Top()
}

func (wb Bound) Bottom() float64 {
	return wb.bound.Bottom()
}

func (wb Bound) Intersects(bounds Bound) bool {
	return wb.bound.Intersects(bounds.bound)
}

func (wb Bound) Extend(radius float64) Bound {
	return Bound{geo.BoundPad(wb.bound, radius)}
}

func (wb Bound) AsFeature(properties map[string]interface{}) *geojson.Feature {
	feature := geojson.NewFeature(wb.bound)

	for name, value := range properties {
		feature.Properties[name] = value
	}

	return feature
}

type Bounds []Bound

func (r Bounds) AsBound() Bound {
	union := r[0].bound
	for _, bound := range r {
		union = union.Union(bound.bound)
	}

	return Bound{union}
}

func (r Bounds) Union() orb.Geometry {
	polygons := lo.Map(r, func(result Bound, _ int) orb.Polygon {
		return result.bound.ToPolygon()
	})
	points := lo.Map(polygons, func(polygon orb.Polygon, _ int) polygol.Geom {
		return g2p(polygon)
	})
	geoms, _ := polygol.Union(points[0], points[1:]...)

	return p2g(geoms)
}

func (r Bounds) AsFeature(properties map[string]interface{}) *geojson.Feature {
	feature := geojson.NewFeature(r.Union())

	for name, value := range properties {
		feature.Properties[name] = value
	}

	return feature
}

func g2p(g orb.Geometry) polygol.Geom {
	var p polygol.Geom

	switch v := g.(type) {
	case orb.Polygon:
		p = make([][][][]float64, 1)
		p[0] = make([][][]float64, len(v))
		for i := range v { // rings
			p[0][i] = make([][]float64, len(v[i]))
			for j := range v[i] { // points
				pt := v[i][j]
				p[0][i][j] = []float64{pt.X(), pt.Y()}
			}
		}
	case orb.MultiPolygon:
		p = make([][][][]float64, len(v))
		for i := range v { // polygons
			p[i] = make([][][]float64, len(v[i]))
			for j := range v[i] { // rings
				p[i][j] = make([][]float64, len(v[i][j]))
				for k := range v[i][j] { // points
					pt := v[i][j][k]
					p[i][j][k] = []float64{pt.X(), pt.Y()}
				}
			}
		}
	}

	return p
}

// source: https://github.com/engelsjk/polygol/tree/4b38d812f2db0cb5ad25d919740693dfa3341382/examples#convversion-functions
func p2g(p [][][][]float64) orb.Geometry {

	g := make(orb.MultiPolygon, len(p))

	for i := range p {
		g[i] = make([]orb.Ring, len(p[i]))
		for j := range p[i] {
			g[i][j] = make([]orb.Point, len(p[i][j]))
			for k := range p[i][j] {
				pt := p[i][j][k]
				point := orb.Point{pt[0], pt[1]}
				g[i][j][k] = point
			}
		}
	}
	return g
}
