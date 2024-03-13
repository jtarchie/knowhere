package marshal

import (
	"fmt"
	"strings"
)

func Tags(tags map[string]string, allowedTags map[string]struct{}) string {
	if len(tags) == 0 {
		return "{}"
	}

	builder := &strings.Builder{}
	builder.WriteByte('{')

	count := 0

	totalCount := len(tags)
	if len(allowedTags) > 0 {
		totalCount = len(allowedTags)
	}

	for key, value := range tags {
		if _, ok := allowedTags[key]; !ok && len(allowedTags) > 0 {
			continue
		}

		marshalString(builder, key)
		builder.WriteByte(':')
		marshalString(builder, value)

		if count < totalCount-1 {
			builder.WriteByte(',')
		}

		count++
	}

	builder.WriteByte('}')

	return builder.String()
}

func marshalString(builder *strings.Builder, str string) {
	builder.WriteByte('"')

	for i := 0; i < len(str); i++ {
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
			//nolint: gomnd
			if char < 0x20 {
				builder.WriteString(fmt.Sprintf("\\u%04x", char))
			} else {
				builder.WriteByte(char)
			}
		}
	}
	builder.WriteByte('"')
}
