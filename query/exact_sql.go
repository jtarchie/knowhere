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
		WHERE
	`, prefix))

	parts := []string{}

	if 0 < len(ast.Types) {
		asString := lo.Map(ast.Types, func(item FilterType, _ int) string {
			return strconv.Itoa(int(item))
		})

		parts = append(
			parts,
			fmt.Sprintf(
				`e.osm_type IN (%s)`,
				strings.Join(asString, ","),
			),
		)
	}

	groupedTags := lo.GroupBy(allowedTags, func(tag FilterTag) OpType {
		return tag.Op
	})

	for operation, tags := range groupedTags {
		switch operation {
		case OpEquals:
			for _, tag := range tags {
				tagPart := "e.tags"
				if tag.Name != "" {
					tagPart += "->>'$." + tag.Name + "'"
				}
				asString := lo.Map(tag.Lookups, func(item string, _ int) string {
					return fmt.Sprintf(
						`%s GLOB '%s'`,
						tagPart,
						item,
					)
				})

				parts = append(
					parts,
					fmt.Sprintf("( %s )", strings.Join(asString, " OR ")),
				)
			}
		case OpNotEquals:
			for _, tag := range tags {
				tagPart := "e.tags"
				if tag.Name != "" {
					tagPart += "->>'$." + tag.Name + "'"
				}
				asString := lo.Map(tag.Lookups, func(item string, _ int) string {
					return fmt.Sprintf(
						`%s NOT GLOB '%s'`,
						tagPart,
						item,
					)
				})

				parts = append(
					parts,
					fmt.Sprintf("( %s )", strings.Join(asString, " OR ")),
				)
			}
		case OpExists:
			for _, tag := range tags {
				parts = append(
					parts,
					fmt.Sprintf(
						"( e.tags->>'$.%s' IS NOT NULL )",
						tag.Name,
					),
				)
			}
		case OpNotExists:
			for _, tag := range tags {
				parts = append(
					parts,
					fmt.Sprintf(
						"( e.tags->>'$.%s' IS NULL )",
						tag.Name,
					),
				)
			}
		}
	}

	if ids, ok := ast.Directives["id"]; ok {
		parts = append(
			parts,
			fmt.Sprintf(
				`e.osm_id IN (%s)`,
				strings.Join(ids, ","),
			),
		)
	}

	builder.WriteString(strings.Join(parts, " AND "))

	return builder.String(), nil
}
