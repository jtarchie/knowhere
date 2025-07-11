package marshal

import (
	"strconv"
	"strings"

	"github.com/jtarchie/knowhere/query"
	"github.com/paulmach/osm"
)

func Members(members osm.Members) string {
	if len(members) == 0 {
		return "[]"
	}

	builder := &strings.Builder{}

	builder.WriteByte('[')

	index := 0

	for _, member := range members {
		switch member.Type {
		case osm.TypeNode, osm.TypeWay, osm.TypeRelation:
			if 0 < index {
				builder.WriteByte(',')
			}

			builder.WriteByte('[')
			builder.WriteString(strconv.FormatInt(member.Ref, 10))
			builder.WriteByte(',')
			switch member.Type {
			case osm.TypeNode:
				builder.WriteString(query.NodeFilter.String())
			case osm.TypeWay:
				builder.WriteString(query.WayFilter.String())
			case osm.TypeRelation:
				builder.WriteString(query.RelationFilter.String())
			}
			builder.WriteByte(',')
			marshalString(builder, member.Role)
			builder.WriteByte(']')

			index++
		case osm.TypeChangeset, osm.TypeNote, osm.TypeUser, osm.TypeBounds:
		}
	}

	builder.WriteByte(']')

	return builder.String()
}
