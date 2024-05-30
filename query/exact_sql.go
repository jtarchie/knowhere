package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

func ToExactSQL(query string) (string, error) {
	ast, err := Parse(query)
	if err != nil {
		return "", fmt.Errorf("could not parse query into SQL: %w", err)
	}

	var (
		builder strings.Builder
		prefix  string
	)

	prefixes, ok := ast.Directives["prefix"]
	if ok && len(prefixes) == 1 {
		prefix = prefixes[0] + "_"
	}

	allowedTags := ast.Tags

	builder.WriteString(fmt.Sprintf(`
		SELECT
			*
		FROM
			%sentries e
	`, prefix))

	builder.WriteString(" WHERE ( e.osm_type IN (")

	for index, t := range ast.Types {
		if 0 < index {
			builder.WriteString(",")
		}

		builder.WriteString(strconv.Itoa(int(t)))
	}

	builder.WriteString(") ) ")

	exists := lo.ContainsBy(allowedTags, func(tag FilterTag) bool {
		return tag.Op == OpEquals || tag.Op == OpExists
	})

	notExists := lo.ContainsBy(allowedTags, func(tag FilterTag) bool {
		return tag.Op == OpNotEquals || tag.Op == OpNotExists
	})

	if exists {
		index := 0

		for _, tag := range allowedTags {
			switch tag.Op {
			case OpEquals:
				builder.WriteString(" AND ( ")

				for index, lookup := range tag.Lookups {
					if 0 < index {
						builder.WriteString(" OR ")
					}

					if tag.Name != "" {
						builder.WriteString("e.tags->>'$.")
						builder.WriteString(tag.Name)
						builder.WriteString("' GLOB ")
					} else {
						builder.WriteString("e.tags")
						builder.WriteString(tag.Name)
						builder.WriteString(" GLOB ")
					}

					builder.WriteString("'")
					builder.WriteString(lookup)
					builder.WriteString("'")
				}

				builder.WriteString(" )")

				index++
			case OpExists:
				builder.WriteString(` AND ( e.tags->>'$.`)
				builder.WriteString(tag.Name)
				builder.WriteString(`' IS NOT NULL )`)

				index++
			case OpNotEquals, OpNotExists:
			}
		}
	}

	if notExists {
		builder.WriteString(" AND NOT ( ")

		index := 0

		for _, tag := range allowedTags {
			switch tag.Op {
			case OpNotEquals:
				if 0 < index {
					builder.WriteString(" AND ")
				}

				builder.WriteString("( ")

				for index, lookup := range tag.Lookups {
					if 0 < index {
						builder.WriteString(" OR ")
					}

					if tag.Name != "" {
						builder.WriteString("e.tags->>'$.")
						builder.WriteString(tag.Name)
						builder.WriteString("' GLOB ")
					} else {
						builder.WriteString("e.tags")
						builder.WriteString(tag.Name)
						builder.WriteString(" GLOB ")
					}

					builder.WriteString("'")
					builder.WriteString(lookup)
					builder.WriteString("'")
				}

				builder.WriteString(" )")

				index++
			case OpNotExists:
				if 0 < index {
					builder.WriteString(" AND ")
				}

				builder.WriteString(`( e.tags->>'$.`)
				builder.WriteString(tag.Name)
				builder.WriteString(`' IS NOT NULL )`)

				index++
			case OpEquals, OpExists:
			}
		}

		builder.WriteString(" )")
	}

	if ids, ok := ast.Directives["id"]; ok {
		builder.WriteString(` AND e.osm_id IN ( `)

		for index, id := range ids {
			if 0 < index {
				builder.WriteString(", ")
			}

			builder.WriteString(id)
		}

		builder.WriteString(` )`)
	}

	return builder.String(), nil
}
