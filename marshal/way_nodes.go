package marshal

import (
	"strconv"
	"strings"

	"github.com/paulmach/osm"
)

func WayNodes(wayNodes osm.WayNodes) string {
	if len(wayNodes) == 0 {
		return "[]"
	}

	builder := &strings.Builder{}

	builder.WriteByte('[')

	for count, value := range wayNodes {
		builder.WriteString(strconv.FormatInt(int64(value.ID), 10))

		if count < len(wayNodes)-1 {
			builder.WriteByte(',')
		}
	}

	builder.WriteByte(']')

	return builder.String()
}
