package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

func ToIndexedSQL(query string) (string, error) {
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

	builder.WriteString(`
		SELECT
			rowid AS id, *
		FROM
			` + prefix + `search s
		WHERE
	`)

	parts := []string{}

	if 0 < len(ast.Types) {
		asString := lo.Map(ast.Types, func(item FilterType, _ int) string {
			return strconv.Itoa(int(item))
		})

		parts = append(
			parts,
			"s.osm_type IN ("+strings.Join(asString, ",")+")",
		)
	}

	groupedTags := lo.GroupBy(allowedTags, func(tag FilterTag) OpType {
		return tag.Op
	})

	equalParts := []string{}
	notParts := []string{}

	for operation, tags := range groupedTags {
		switch operation {
		case OpEquals:
			for _, tag := range tags {
				asString := lo.Map(tag.Lookups, func(item string, _ int) string {
					return `"` + item + `"`
				})

				if tag.Name == "" {
					equalParts = append(
						equalParts,
						"( "+strings.Join(asString, " OR ")+" )",
					)
				} else {
					equalParts = append(
						equalParts,
						`( "`+tag.Name+`" AND ( `+strings.Join(asString, " OR ")+" ) )",
					)

					asString = lo.Map(tag.Lookups, func(item string, _ int) string {
						return "s.tags->>'$." + tag.Name + "' = '" + item + "'"
					})
					parts = append(parts,
						"( "+strings.Join(asString, " OR ")+" )",
					)
				}
			}
		case OpNotEquals:
			for _, tag := range tags {
				asString := lo.Map(tag.Lookups, func(item string, _ int) string {
					return `"` + item + `"`
				})

				if tag.Name == "" {
					notParts = append(
						notParts,
						"( "+strings.Join(asString, " OR ")+" )",
					)
				} else {
					notParts = append(
						notParts,
						`( "`+tag.Name+`" AND ( `+strings.Join(asString, " OR ")+" ) )",
					)
				}
			}
		case OpExists:
			for _, tag := range tags {
				equalParts = append(
					equalParts,
					`( "`+tag.Name+`" )`,
				)

				parts = append(parts, "( s.tags->>'$."+tag.Name+"' IS NOT NULL )")
			}
		case OpNotExists:
			for _, tag := range tags {
				notParts = append(
					notParts,
					`( "`+tag.Name+`" )`,
				)

				parts = append(parts, "( s.tags->>'$."+tag.Name+"' IS NULL )")
			}
		case OpGreaterThan, OpGreaterThanEquals, OpLessThan, OpLessThanEquals:
			for _, tag := range tags {
				equalParts = append(
					equalParts,
					`( "`+tag.Name+`" )`,
				)

				parts = append(
					parts,
					"( s.tags->>'$."+tag.Name+"' "+operation.String()+" "+tag.Lookups[0]+" )",
				)
			}
		case OpContains:
			for _, tag := range tags {
				asString := lo.Map(tag.Lookups, func(item string, _ int) string {
					if item[len(item)-1] == '*' {
						return `"` + item[0:len(item)-1] + `"*`
					}

					return `"` + item + `"`
				})

				equalParts = append(
					equalParts,
					`( "`+tag.Name+`" AND ( `+strings.Join(asString, " OR ")+" ) )",
				)

				for _, lookup := range tag.Lookups {
					parts = append(
						parts,
						"( LOWER(s.tags->>'$."+tag.Name+"') GLOB '*"+strings.ToLower(lookup)+"*' )",
					)
				}
			}
		case OpNotContains:
			for _, tag := range tags {
				asString := lo.Map(tag.Lookups, func(item string, _ int) string {
					return `"` + item + `"`
				})

				notParts = append(
					notParts,
					`( "`+tag.Name+`" AND ( `+strings.Join(asString, " OR ")+" ) )",
				)

				for _, lookup := range tag.Lookups {
					parts = append(
						parts,
						"( LOWER(s.tags->>'$."+tag.Name+"') NOT GLOB '*"+strings.ToLower(lookup)+"*' )",
					)
				}
			}
		}
	}

	if 0 < len(equalParts) {
		equals := strings.Join(equalParts, " AND ")
		if 0 < len(notParts) {
			equals += " NOT " + strings.Join(notParts, " AND ")
		}
		parts = append(
			parts,
			"s.tags MATCH '"+equals+"'",
		)
	}

	if ids, ok := ast.Directives["id"]; ok {
		parts = append(
			parts,
			"s.osm_id IN ("+strings.Join(ids, ",")+")",
		)
	}

	if bbs, ok := ast.Directives["bb"]; ok {
		if 4 == len(bbs) {
			parts = append(
				parts,
				bbs[0]+" <= s.minLon",
				bbs[1]+" <= s.minLat",
				"s.maxLon <= "+bbs[2],
				"s.maxLat <= "+bbs[3],
			)
		}
	}

	for index, part := range parts {
		if 0 < index {
			builder.WriteString(" AND ")
		}
		builder.WriteString(part)
	}
	return builder.String(), nil
}
