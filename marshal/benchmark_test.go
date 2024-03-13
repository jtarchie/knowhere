package marshal_test

import (
	"encoding/json"
	"testing"

	"github.com/jtarchie/knowhere/marshal"
	"github.com/paulmach/osm"
)

//nolint: gochecknoglobals
var payload string

func BenchmarkControlTags(b *testing.B) {
	for n := 0; n < b.N; n++ {
		//nolint: errchkjson
		_, _ = json.Marshal(map[string]string{
			"amenity": "cafe",
			"name":    "Starbucks",
		})
	}
}

func BenchmarkTags(b *testing.B) {
	for n := 0; n < b.N; n++ {
		payload = marshal.Tags(map[string]string{
			"amenity": "cafe",
			"name":    "Starbucks",
		}, nil)
	}
}

func BenchmarkMembers(b *testing.B) {
	for n := 0; n < b.N; n++ {
		payload = marshal.Members(osm.Members{
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
		payload = marshal.WayNodes(osm.WayNodes{
			osm.WayNode{ID: 1},
			osm.WayNode{ID: 2},
		})
	}
}
