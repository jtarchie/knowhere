package query

import (
  "sort"
  "fmt"

  "github.com/samber/lo"
)

%%{
  machine query;
  write data;
}%%

func Parse(data string) (*AST, error) {
  // types used for the AST
  foundTypes := []FilterType{}
  tags := []FilterTag{}
  directives := map[string]FilterDirective{}
  
  var (
    tag FilterTag
    directiveName string
    directive FilterDirective
  )

  // set defaults for state machine parsing
  cs, p, pe, eof := 0, 0, len(data), len(data)
  
  // tracks where the beginning of a word starts
  mark := 0
  
  // keep track of opening and closing brackets
  brackets := 0


  %%{
    action mark { mark = p}

    action node     { foundTypes = append(foundTypes, NodeFilter) }
    action relation { foundTypes = append(foundTypes, RelationFilter) }
    action way      { foundTypes = append(foundTypes, WayFilter) }
    action all      { foundTypes = append(foundTypes, NodeFilter, WayFilter, RelationFilter) }
    action type_error {
      return nil, fmt.Errorf("an undefined type was specified %c: %w", data[p], ErrUndefinedFilter)
    }
    action create_tag { tag = FilterTag{Lookups: []string{}} }
    action append_tag { tags = append(tags, tag) }

    action tag_eq     { tag.Op = OpEquals }
    action tag_ne     { tag.Op = OpNotEquals }
    action tag_exists { tag.Op = OpExists }
    action tag_not    { tag.Op = OpNotExists }

    action tag_name  { tag.Name    = data[mark:p] }
    action tag_value { tag.Lookups = append(tag.Lookups, data[mark:p]) }

    action directive_name  { directiveName  = data[mark:p] }
    action directive_value { directive = append(directive, data[mark:p]) }

    action create_directive { directive = FilterDirective{} }
    action append_directive { directives[directiveName] = directive }

    action inc_bracket { brackets++ }
    action dec_bracket { brackets-- }

    type  = ("n" >node) | ("r" >relation) | ("w" >way);
    types = (type+ | ("*" >all)) %!type_error;
    
    tag_name = (alnum+ >mark %tag_name) | "*";
    tag_value_unquoted = [^,\"\]]+ >mark %tag_value;
    tag_value_quoted   = '"' ([^"]+ >mark %tag_value) '"';
    tag_value = tag_value_quoted | tag_value_unquoted;
    tag_values = tag_value ( "," tag_value )*;
    tag_eq = (
      tag_name "=" tag_values
    ) %tag_eq;
    tag_ne = (
      tag_name "!=" tag_values
    ) %tag_ne;
    tag_exists = (tag_name) %tag_exists;
    tag_not    = ("!" tag_name) %tag_not;
    tag    = ("[" %inc_bracket) (tag_eq | tag_ne | tag_exists | tag_not) ("]" %dec_bracket);
    tags   = (tag >create_tag %append_tag)*;

    directive_value_unquoted = [^,\"\]]+ >mark %directive_value;
    directive_value_quoted   = '"' ([^"]+ >mark %directive_value) '"';
    directive_value = directive_value_quoted | directive_value_unquoted;
    directive_values = directive_value ( "," directive_value )*;
    directive_name = (alnum+ >mark %directive_name) | "*";
    directive = ("(" %inc_bracket) (directive_name "=" directive_values ) (")" %dec_bracket);
    directives = (directive >create_directive %append_directive)*;

    main := types tags directives;
    write init;
    write exec;
  }%%

  if cs < query_first_final {
    return nil, ErrUnparsableQuery
  }

  if brackets != 0 {
    return nil, fmt.Errorf("tags not enclosed properly (%d): %w", brackets, ErrUnbalancedBrackets)
  }

  sort.Slice(foundTypes, func(i, j int) bool {
		return foundTypes[i] < foundTypes[j]
	})

	foundTypes = lo.Uniq(foundTypes)

	return &AST{
		Tags:  tags,
		Types: foundTypes,
    Directives: directives,
	}, nil
}