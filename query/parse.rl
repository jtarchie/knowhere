package query

import (
  "sort"

  "github.com/samber/lo"
)

%%{
  machine syslog_rfc5424;
  write data;
}%%

func Parse(data string) (*AST, error) {
  // set defaults for state machine parsing
  cs, p, pe := 0, 0, len(data)

  foundTypes := []FilterType{}
  tags := []FilterTag{}

  %%{
    action area     { foundTypes = append(foundTypes, AreaFilter) }
    action node     { foundTypes = append(foundTypes, NodeFilter) }
    action relation { foundTypes = append(foundTypes, RelationFilter) }
    action way      { foundTypes = append(foundTypes, WayFilter) }

    action all { foundTypes = append(foundTypes, NodeFilter, AreaFilter, WayFilter, RelationFilter) }

    type  = ("a" >area) | ("n" >node) | ("r" >relation) | ("w" >way);
    types = type+ | ("*" >all);
    

    main := types;
    write init;
    write exec;
  }%%

  sort.Slice(foundTypes, func(i, j int) bool {
		return foundTypes[i] < foundTypes[j]
	})

	foundTypes = lo.Uniq(foundTypes)

	return &AST{
		Types: foundTypes,
		Tags:  tags,
	}, nil
}