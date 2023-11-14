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

	sql := `
		SELECT
			*
		FROM
			entries e
		-- JOIN
		--	search s
		-- ON search.id = e.id
	`

	if 0 < len(ast.Types) {
		sql += " WHERE "

		conds := []string{}

		for _, t := range ast.Types {
			switch t {
			case NodeFilter:
				conds = append(conds, "e.osm_type = 'node'")
			case AreaFilter:
				conds = append(conds, "e.osm_type = 'area'")
			case WayFilter:
				conds = append(conds, "e.osm_type = 'way'")
			case RelationFilter:
				conds = append(conds, "e.osm_type = 'relation'")
			}
		}

		sql += "(" + strings.Join(conds, " OR ") + ")"
	}

	return sql, nil
}
