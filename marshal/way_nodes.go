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
		if 0 < count {
			builder.WriteByte(',')
		}

		builder.WriteString(strconv.FormatInt(int64(value.ID), 10))
	}

	builder.WriteByte(']')

	return builder.String()
}
