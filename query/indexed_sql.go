package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mmcloughlin/geohash"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/samber/lo"
)

func ToIndexedSQL(query string) (string, error) {
	ast, err := Parse(query)
	if err != nil {
		return "", fmt.Errorf("could not parse query into SQL: %w", err)
	}

	var (
		builder strings.Builder
		area    string
	)

	areas, ok := ast.Directives["area"]
	if ok && len(areas) == 1 {
		area = areas[0] + "_"
	}

	allowedTags := ast.Tags

	builder.WriteString(`
		SELECT
			rowid AS id, *
		FROM
			` + area + `search s
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

	if bbs, ok := ast.Directives["bb"]; ok {
		if 4 == len(bbs) {
			parts = append(
				parts,
				bbs[0]+" <= s.minLon",
				bbs[1]+" <= s.minLat",
				"s.maxLon <= "+bbs[2],
				"s.maxLat <= "+bbs[3],
			)

			// find distance of lat/lng
			bounds := orb.Bound{
				Min: orb.Point{toFloat(bbs[0]), toFloat(bbs[1])},
				Max: orb.Point{toFloat(bbs[2]), toFloat(bbs[3])},
			}
			precision := boundsToGeohashPrecision(bounds)
			center := bounds.Center()

			// get all neighboring hashes
			hash := geohash.EncodeWithPrecision(center.Lat(), center.Lon(), precision)
			hashes := []string{hash}
			hashes = append(hashes, geohash.Neighbors(hash)...)

			// add to full text search
			asString := lo.Map(hashes, func(item string, _ int) string {
				return item + `*`
			})
			equalParts = append(
				equalParts,
				`( "geohash" AND ( `+strings.Join(asString, " OR ")+" ) )",
			)
		}
	}

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

	for index, part := range parts {
		if 0 < index {
			builder.WriteString(" AND ")
		}
		builder.WriteString(part)
	}
	return builder.String(), nil
}

func toFloat(s string) float64 {
	value, _ := strconv.ParseFloat(s, 64)
	return value
}

func boundsToGeohashPrecision(bounds orb.Bound) uint {
	geohashPrecisions := [][2]float64{
		{0.0372, 0.0186},
		{0.149, 0.149},
		{1.19, 0.596},
		{4.77, 4.77},
		{38.2, 19.1},
		{153, 153},
		{1220, 610},
		{4890, 4890},
		{39100, 19500},
		{156000, 156000},
		{1250000, 625000},
		{5000000, 5000000},
	}

	height := geo.BoundHeight(bounds)
	width := geo.BoundWidth(bounds)

	for index, precision := range geohashPrecisions {
		if height <= precision[0] && width <= precision[1] {
			return max(0, uint(len(geohashPrecisions)-index+1))
		}
	}

	return uint(len(geohashPrecisions))
}
