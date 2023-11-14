package query

import (
  "sort"
  "fmt"

  "github.com/samber/lo"
)

%%{
  machine syslog_rfc5424;
  write data;
}%%

func Parse(data string) (*AST, error) {
  // set defaults for state machine parsing
  cs, p, pe, eof := 0, 0, len(data), len(data)
  foundTypes := []FilterType{}
  tags := []FilterTag{}
  var tag FilterTag
  mark := 0

  %%{
    action mark { mark = p}

    action area     { foundTypes = append(foundTypes, AreaFilter) }
    action node     { foundTypes = append(foundTypes, NodeFilter) }
    action relation { foundTypes = append(foundTypes, RelationFilter) }
    action way      { foundTypes = append(foundTypes, WayFilter) }
    action all      { foundTypes = append(foundTypes, NodeFilter, AreaFilter, WayFilter, RelationFilter) }
    action type_error {
      return nil, fmt.Errorf("an undefined type was specified %c: %w", data[p], ErrUndefinedFilter)
    }
    action create_tag { tag = FilterTag{Lookups: []string{}} }
    action append_tag { tags = append(tags, tag) }

    action tag_equals { tag.Op = OpEquals }
    action tag_exists { tag.Op = OpExists }
    action tag_not    { tag.Op = OpNotExist }

    action tag_name  { tag.Name    = data[mark:p] }
    action tag_value { tag.Lookups = append(tag.Lookups, data[mark:p]) }

    type  = ("a" >area) | ("n" >node) | ("r" >relation) | ("w" >way);
    types = (type+ | ("*" >all)) %!type_error;
    
    tag_name = alnum+ >mark %tag_name;
    tag_value = alnum+ >mark %tag_value;
    tag_eq = (
      "[" tag_name "=" tag_value ( "," tag_value )* "]"
    ) %tag_equals;
    tag_exists = ("[" tag_name "]") %tag_exists;
    tag_not    = ("[!" tag_name "]") %tag_not;
    tag    = tag_eq | tag_exists | tag_not;
    tags   = (tag >create_tag %append_tag)*;

    main := types tags;
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