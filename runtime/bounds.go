package runtime

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

type WrappedBound struct {
	orb.Bound
}

func (wr *WrappedBound) Intersects(bounds *WrappedBound) bool {
	return wr.Bound.Intersects(bounds.Bound)
}

// Extends a bounding box in kilometers in each direction.
// This is for best effort, not exact.
func (wb *WrappedBound) Extend(radius float64) *WrappedBound {
	return &WrappedBound{geo.BoundPad(wb.Bound, radius*1000)}
}
