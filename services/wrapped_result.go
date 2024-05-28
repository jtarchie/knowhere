package services

import (
	"github.com/jtarchie/knowhere/query"
	"github.com/paulmach/orb"
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
