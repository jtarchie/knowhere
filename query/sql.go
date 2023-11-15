package query

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

func ToSQL(query string) (string, error) {
	ast, err := Parse(query)
	if err != nil {
		return "", fmt.Errorf("could not parse query into SQL: %w", err)
	}

	var builder strings.Builder

	builder.WriteString(`
		SELECT
			*
		FROM
			entries e
	`)

	if 0 < len(ast.Tags) {
		builder.WriteString(`
			JOIN
				search s
			ON
				s.rowid = e.id
		`)
	}

	builder.WriteString(" WHERE ")

	if 0 < len(ast.Types) {
		builder.WriteString("(")

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

		builder.WriteString(") ")
	}

	if 0 < len(ast.Tags) {
		exists := lo.ContainsBy(ast.Tags, func(tag FilterTag) bool {
			return tag.Op == OpEquals || tag.Op == OpExists
		})

		if exists {
			builder.WriteString("AND s.tags MATCH '")
		}

		for index, tag := range ast.Tags {
			if 0 < index {
				builder.WriteString(" AND ")
			}

			switch tag.Op {
			case OpEquals:
				builder.WriteString("( ")

				for index, lookup := range tag.Lookups {
					if 0 < index {
						builder.WriteString(" OR ")
					}

					builder.WriteString(`("`)
					builder.WriteString(tag.Name)
					builder.WriteString(" ")
					builder.WriteString(lookup)
					builder.WriteString(`")`)
				}

				builder.WriteString(" )")
			case OpExists:
				builder.WriteString(`( "`)
				builder.WriteString(tag.Name)
				builder.WriteString(`" )`)
			case OpNotEquals, OpNotExists:
			}
		}

		builder.WriteString("'")
	}

	return builder.String(), nil
}
