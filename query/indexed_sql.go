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

	builder.WriteString(fmt.Sprintf(`
		SELECT
			rowid AS id, *
		FROM
			%ssearch s
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
				`s.osm_type IN (%s)`,
				strings.Join(asString, ","),
			),
		)
	}

	groupedTags := lo.GroupBy(allowedTags, func(tag FilterTag) OpType {
		return tag.Op
	})

	equalParts := []string{}
	notParts := []string{}

	for op := range OpLessThanEquals {
		tags, ok := groupedTags[op]
		if !ok {
			continue
		}

		switch op {
		case OpEquals:
			for _, tag := range tags {
				asString := lo.Map(tag.Lookups, func(item string, _ int) string {
					return fmt.Sprintf("%q", item)
				})

				if tag.Name == "" {
					equalParts = append(
						equalParts,
						fmt.Sprintf(
							"( %s )",
							strings.Join(asString, " OR "),
						),
					)
				} else {
					equalParts = append(
						equalParts,
						fmt.Sprintf(
							"( %q AND ( %s ) )",
							tag.Name,
							strings.Join(asString, " OR "),
						),
					)
				}
			}
		case OpNotEquals:
			for _, tag := range tags {
				asString := lo.Map(tag.Lookups, func(item string, _ int) string {
					return fmt.Sprintf("%q", item)
				})

				if tag.Name == "" {
					notParts = append(
						notParts,
						fmt.Sprintf(
							"( %s )",
							strings.Join(asString, " OR "),
						),
					)
				} else {
					notParts = append(
						notParts,
						fmt.Sprintf(
							"( %q AND ( %s ) )",
							tag.Name,
							strings.Join(asString, " OR "),
						),
					)
				}
			}
		case OpExists:
			for _, tag := range tags {
				equalParts = append(
					equalParts,
					fmt.Sprintf(
						`( "%s" )`,
						tag.Name,
					),
				)
			}
		case OpNotExists:
			for _, tag := range tags {
				notParts = append(
					notParts,
					fmt.Sprintf(
						`( "%s" )`,
						tag.Name,
					),
				)
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
			fmt.Sprintf(
				"s.tags MATCH '%s'",
				equals,
			),
		)
	}

	if ids, ok := ast.Directives["id"]; ok {
		parts = append(
			parts,
			fmt.Sprintf(
				`s.osm_id IN (%s)`,
				strings.Join(ids, ","),
			),
		)
	}

	builder.WriteString(strings.Join(parts, " AND "))
	builder.WriteString(" ORDER BY rank")

	return builder.String(), nil
}
