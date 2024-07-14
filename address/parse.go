package address

func Parse(fullAddress string) (map[string]string, bool) {
	for _, parser := range addressParsers {
		match := parser.FindStringSubmatch(fullAddress)
		if len(match) == 0 {
			continue
		}

		results := map[string]string{}
		subnames := parser.SubexpNames()
		for i, name := range match[1:] {
			results[subnames[i+1]] = name
		}

		return results, true
	}

	return nil, false
}
