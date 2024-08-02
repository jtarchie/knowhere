package address

import (
	"regexp"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/pioz/countries"
	"github.com/samber/lo"
)

var states = lo.GroupBy(lo.FlatMap(lo.ToPairs(countries.Get("US").Subdivisions), func(entry lo.Entry[string, countries.Subdivision], _ int) []string {
	return []string{entry.Key, entry.Value.Name}
}), func(value string) string {
	return strings.ToLower(value[0:1])
})

var streets = map[string]string{
	"Ave":  "Avenue",
	"Blvd": "Boulevard",
	"Cir":  "Circle",
	"Ct":   "Court",
	"Dr":   "Drive",
	"Ln":   "Lane",
	"Pkwy": "Parkway",
	"Pl":   "Place",
	"Rd":   "Road",
	"Sq":   "Square",
	"St":   "Street",
	"Ter":  "Terrace",
	"Way":  "Way",
}
var streetsMatcher = regexp.MustCompile(`\b(` + strings.Join(lo.Keys(streets), "|") + `)\b`)

func Parse(fullAddress string) (map[string]string, bool) {
	fullAddress = strings.TrimSpace(fullAddress)

	for _, parser := range addressParsers {
		match := parser.FindStringSubmatch(fullAddress)
		if len(match) == 0 {
			continue
		}

		results := map[string]string{}
		subnames := parser.SubexpNames()
		for i, name := range match[1:] {
			if len(name) > 0 && name[len(name)-1] == ',' {
				name = name[0 : len(name)-1]
			}
			results[subnames[i+1]] = name
		}

		if state, ok := results["state"]; ok && state != "" {
			matches := fuzzy.RankFindNormalizedFold(state, states[strings.ToLower(state[0:1])])
			if len(matches) == 0 {
				continue
			}

			sort.Sort(matches)
			results["state"] = matches[0].Target
		}

		if road, ok := results["road"]; ok {
			results["road"] = streetsMatcher.ReplaceAllStringFunc(road, func(abbr string) string {
				return streets[abbr]
			})
		}

		return results, true
	}

	return nil, false
}
