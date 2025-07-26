package query

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mmcloughlin/geohash"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/samber/lo"
)

func escape(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

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
			// Parse and validate coordinate values
			minLon := toFloat(bbs[0])
			minLat := toFloat(bbs[1])
			maxLon := toFloat(bbs[2])
			maxLat := toFloat(bbs[3])

			// Validate coordinates
			if err := validateCoordinates(minLon, minLat, maxLon, maxLat); err != nil {
				return "", fmt.Errorf("invalid bounding box coordinates: %w", err)
			}

			parts = append(
				parts,
				bbs[0]+" <= s.minLon",
				bbs[1]+" <= s.minLat",
				"s.maxLon <= "+bbs[2],
				"s.maxLat <= "+bbs[3],
			)

			// find distance of lat/lng
			bounds := orb.Bound{
				Min: orb.Point{minLon, minLat},
				Max: orb.Point{maxLon, maxLat},
			}
			precision := boundsToGeohashPrecision(bounds)
			center := bounds.Center()

			// Validate center coordinates before geohash encoding
			centerLat, centerLon := center.Lat(), center.Lon()
			if math.IsNaN(centerLat) || math.IsNaN(centerLon) ||
				math.IsInf(centerLat, 0) || math.IsInf(centerLon, 0) ||
				centerLat < -90 || centerLat > 90 ||
				centerLon < -180 || centerLon > 180 {
				return "", fmt.Errorf("invalid center coordinates: lat=%f, lon=%f", centerLat, centerLon)
			}

			// get all neighboring hashes
			hash := geohash.EncodeWithPrecision(centerLat, centerLon, precision)
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
					return `"` + escape(item) + `"`
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
						return "s.tags->>'$." + tag.Name + "' = '" + escape(item) + "'"
					})
					parts = append(parts,
						"( "+strings.Join(asString, " OR ")+" )",
					)
				}
			}
		case OpNotEquals:
			for _, tag := range tags {
				asString := lo.Map(tag.Lookups, func(item string, _ int) string {
					return `"` + escape(item) + `"`
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

					asString = lo.Map(tag.Lookups, func(item string, _ int) string {
						return "s.tags->>'$." + tag.Name + "' <> '" + escape(item) + "'"
					})
					parts = append(parts,
						"( "+strings.Join(asString, " OR ")+" )",
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
						return `"` + escape(item[0:len(item)-1]) + `"*`
					}

					return `"` + escape(item) + `"`
				})

				equalParts = append(
					equalParts,
					`( "`+tag.Name+`" AND ( `+strings.Join(asString, " OR ")+" ) )",
				)

				for _, lookup := range tag.Lookups {
					parts = append(
						parts,
						"( LOWER(s.tags->>'$."+tag.Name+"') GLOB '*"+escape(strings.ToLower(lookup))+"*' )",
					)
				}
			}
		case OpNotContains:
			for _, tag := range tags {
				asString := lo.Map(tag.Lookups, func(item string, _ int) string {
					return `"` + escape(item) + `"`
				})

				notParts = append(
					notParts,
					`( "`+tag.Name+`" AND ( `+strings.Join(asString, " OR ")+" ) )",
				)

				for _, lookup := range tag.Lookups {
					parts = append(
						parts,
						"( LOWER(s.tags->>'$."+tag.Name+"') NOT GLOB '*"+escape(strings.ToLower(lookup))+"*' )",
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

func validateCoordinates(minLon, minLat, maxLon, maxLat float64) error {
	// Check for NaN or Infinity values
	coords := []float64{minLon, minLat, maxLon, maxLat}
	names := []string{"minLon", "minLat", "maxLon", "maxLat"}

	for i, coord := range coords {
		if math.IsNaN(coord) || math.IsInf(coord, 0) {
			return fmt.Errorf("invalid coordinates: %s=%f (NaN or Infinity values not allowed)", names[i], coord)
		}
	}

	// Validate latitude bounds
	if minLat < -90 || minLat > 90 {
		return fmt.Errorf("invalid latitude: %f (must be between -90 and 90)", minLat)
	}
	if maxLat < -90 || maxLat > 90 {
		return fmt.Errorf("invalid latitude: %f (must be between -90 and 90)", maxLat)
	}

	// Validate longitude bounds
	if minLon < -180 || minLon > 180 {
		return fmt.Errorf("invalid longitude: %f (must be between -180 and 180)", minLon)
	}
	if maxLon < -180 || maxLon > 180 {
		return fmt.Errorf("invalid longitude: %f (must be between -180 and 180)", maxLon)
	}

	// Validate bounds order
	if minLon > maxLon {
		return fmt.Errorf("invalid bounds: minLon (%f) must be <= maxLon (%f)", minLon, maxLon)
	}
	if minLat > maxLat {
		return fmt.Errorf("invalid bounds: minLat (%f) must be <= maxLat (%f)", minLat, maxLat)
	}

	return nil
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

	// Handle invalid bounds
	if math.IsNaN(height) || math.IsNaN(width) || math.IsInf(height, 0) || math.IsInf(width, 0) {
		return 1 // Return minimum precision for invalid bounds
	}

	for index, precision := range geohashPrecisions {
		if height <= precision[0] && width <= precision[1] {
			result := len(geohashPrecisions) - index + 1
			if result < 0 {
				return 1
			}
			return uint(result)
		}
	}

	return uint(len(geohashPrecisions))
}
