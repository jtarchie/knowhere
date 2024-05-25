package query

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

func ToExactSQL(query string) (string, error) {
	ast, err := Parse(query)
	if err != nil {
		return "", fmt.Errorf("could not parse query into SQL: %w", err)
	}

	var builder strings.Builder

	prefixDirective, ok := lo.Find(ast.Directives, func(directive FilterDirective) bool {
		return directive.Name == "prefix" && len(directive.Value) > 0
	})

	prefix := ""
	if ok {
		prefix = prefixDirective.Value + "_"
	}

	allowedTags := ast.Tags

	builder.WriteString(fmt.Sprintf(`
		SELECT
			*
		FROM
			%sentries e
	`, prefix))

	builder.WriteString(" WHERE ( ")

	for index, t := range ast.Types {
		if 0 < index {
			builder.WriteString(" OR ")
		}

		switch t {
		case NodeFilter:
			builder.WriteString("e.osm_type = 'node'")
		case AreaFilter:
			builder.WriteString("e.osm_type = 'area'")
		case WayFilter:
			builder.WriteString("e.osm_type = 'way'")
		case RelationFilter:
			builder.WriteString("e.osm_type = 'relation'")
		}
	}

	builder.WriteString(" ) ")

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

	return builder.String(), nil
}
