package marshal

import (
	"fmt"

	"github.com/paulmach/osm"
)

func Members(members osm.Members) string {
	if len(members) == 0 {
		return "[]"
	}

	payload := "["

	index := 0

	for _, member := range members {
		switch member.Type {
		case osm.TypeNode, osm.TypeWay, osm.TypeRelation:
			if 0 < index {
				payload += ","
			}

			payload += fmt.Sprintf("[%d,%q,%q]", member.Ref, member.Type, member.Role)
			index++
		case osm.TypeChangeset, osm.TypeNote, osm.TypeUser, osm.TypeBounds:
		}
	}

	payload += "]"

	return payload
}
