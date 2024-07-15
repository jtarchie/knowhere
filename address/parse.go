package address

import (
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

func Parse(fullAddress string) (map[string]string, bool) {

	for _, parser := range addressParsers {
		match := parser.FindStringSubmatch(fullAddress)
		if len(match) == 0 {
			continue
		}

		results := map[string]string{}
		subnames := parser.SubexpNames()
		for i, name := range match[1:] {
			if name[len(name)-1] == ',' {
				name = name[0 : len(name)-1]
			}
			results[subnames[i+1]] = name
		}

		if state, ok := results["state"]; ok {
			matches := fuzzy.RankFindNormalizedFold(state, states[strings.ToLower(state[0:1])])
			if len(matches) == 0 {
				continue
			}

			sort.Sort(matches)
			results["state"] = matches[0].Target
		}

		return results, true
	}

	return nil, false
}
