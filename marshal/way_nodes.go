package marshal

import (
	"fmt"

	"github.com/paulmach/osm"
)

func WayNodes(wayNodes osm.WayNodes) string {
	if len(wayNodes) == 0 {
		return "[]"
	}

	payload := "["

	for count, value := range wayNodes {
		payload += fmt.Sprintf("%d", value.ID)
		if count < len(wayNodes)-1 {
			payload += ","
		}
	}

	payload += "]"

	return payload
}
