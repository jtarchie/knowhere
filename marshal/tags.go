package marshal

import (
	"fmt"
	"strings"
)

func Tags(tags map[string]string) string {
	if len(tags) == 0 {
		return "{}"
	}

	builder := &strings.Builder{}
	builder.WriteByte('{')

	count := 0

	for key, value := range tags {
		marshalString(builder, key)
		builder.WriteByte(':')
		marshalString(builder, value)

		if count < len(tags)-1 {
			builder.WriteByte(',')
		}

		count++
	}

	builder.WriteByte('}')

	return builder.String()
}

func marshalString(builder *strings.Builder, str string) {
	builder.WriteByte('"')

	for i := range len(str) {
		char := str[i]
		switch char {
		case '\\', '"':
			builder.WriteByte('\\')
			builder.WriteByte(char)
		case '\n':
			builder.WriteString("\\n")
		case '\r':
			builder.WriteString("\\r")
		case '\t':
			builder.WriteString("\\t")
		case '\b':
			builder.WriteString("\\b")
		case '\f':
			builder.WriteString("\\f")
		default:
			const unicodeThreshold = 0x20
			if char < unicodeThreshold {
				builder.WriteString(fmt.Sprintf("\\u%04x", char))
			} else {
				builder.WriteByte(char)
			}
		}
	}

	builder.WriteByte('"')
}
