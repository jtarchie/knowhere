package runtime

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

type Bound struct {
	orb.Bound
}

func (wr *Bound) Intersects(bounds *Bound) bool {
	return wr.Bound.Intersects(bounds.Bound)
}

// Extends a bounding box in kilometers in each direction.
// This is for best effort, not exact.
func (wb *Bound) Extend(radius float64) *Bound {
	return &Bound{geo.BoundPad(wb.Bound, radius*1000)}
}
