package marshal_test

import (
	"testing"

	"github.com/jtarchie/knowhere/marshal"
	"github.com/paulmach/osm"
)

//nolint: gochecknoglobals
var json string

func BenchmarkTags(b *testing.B) {
	for n := 0; n < b.N; n++ {
		json = marshal.Tags(map[string]string{
			"a": "b",
			"c": "d",
		})
	}
}

func BenchmarkMembers(b *testing.B) {
	for n := 0; n < b.N; n++ {
		json = marshal.Members(osm.Members{
			osm.Member{
				Type: osm.TypeNode,
				Ref:  1,
				Role: "donotcare",
			},
			osm.Member{
				Type: osm.TypeWay,
				Ref:  2,
				Role: "donotcare",
			},
		})
	}
}

func BenchmarkWayNodes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		json = marshal.WayNodes(osm.WayNodes{
			osm.WayNode{ID: 1},
			osm.WayNode{ID: 2},
		})
	}
}
