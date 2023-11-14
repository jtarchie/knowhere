package query

import (
	"fmt"
	"strings"
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

	if 0 < len(ast.Types) {
		builder.WriteString(" WHERE (")

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

		builder.WriteString(")")
	}

	return builder.String(), nil
}
