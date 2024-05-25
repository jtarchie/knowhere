package query

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

func ToIndexedSQL(query string) (string, error) {
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
		JOIN
			%ssearch s
		ON
			s.rowid = e.id
	`, prefix, prefix))

	builder.WriteString(" WHERE ")

	builder.WriteString("s.osm_type MATCH '")

	for index, t := range ast.Types {
		if 0 < index {
			builder.WriteString(" OR ")
		}

		switch t {
		case NodeFilter:
			builder.WriteString("node")
		case AreaFilter:
			builder.WriteString("area")
		case WayFilter:
			builder.WriteString("way")
		case RelationFilter:
			builder.WriteString("relation")
		}
	}

	builder.WriteString("' ")

	exists := lo.ContainsBy(allowedTags, func(tag FilterTag) bool {
		return tag.Op == OpEquals || tag.Op == OpExists
	})

	notExists := lo.ContainsBy(allowedTags, func(tag FilterTag) bool {
		return tag.Op == OpNotEquals || tag.Op == OpNotExists
	})

	if exists || notExists {
		builder.WriteString("AND s.tags MATCH '")
	}

	if exists {
		index := 0

		for _, tag := range allowedTags {
			switch tag.Op {
			case OpEquals:
				if 0 < index {
					builder.WriteString(" AND ")
				}

				builder.WriteString("( ")

				for index, lookup := range tag.Lookups {
					if 0 < index {
						builder.WriteString(" OR ")
					}

					builder.WriteString(`("`)

					if tag.Name != "" {
						builder.WriteString(tag.Name)
						builder.WriteString(" ")
					}

					builder.WriteString(lookup)
					builder.WriteString(`")`)
				}

				builder.WriteString(" )")

				index++
			case OpExists:
				if 0 < index {
					builder.WriteString(" AND ")
				}

				builder.WriteString(`( "`)
				builder.WriteString(tag.Name)
				builder.WriteString(`" )`)

				index++
			case OpNotEquals, OpNotExists:
			}
		}
	}

	if notExists {
		builder.WriteString(" NOT ")

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

					builder.WriteString(`("`)

					if tag.Name != "" {
						builder.WriteString(tag.Name)
						builder.WriteString(" ")
					}

					builder.WriteString(lookup)
					builder.WriteString(`")`)
				}

				builder.WriteString(" )")

				index++
			case OpNotExists:
				builder.WriteString(`( "`)
				builder.WriteString(tag.Name)
				builder.WriteString(`" )`)
			case OpEquals, OpExists:
			}
		}
	}

	if exists || notExists {
		builder.WriteString("'")
	}

	return builder.String(), nil
}
