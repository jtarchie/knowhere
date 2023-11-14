
//line parse.rl:1
package query

import (
  "sort"

  "github.com/samber/lo"
)


//line parse.go:13
const syslog_rfc5424_start int = 1
const syslog_rfc5424_first_final int = 2
const syslog_rfc5424_error int = 0

const syslog_rfc5424_en_main int = 1


//line parse.rl:12


func Parse(data string) (*AST, error) {
  // set defaults for state machine parsing
  cs, p, pe := 0, 0, len(data)

  foundTypes := []FilterType{}
  tags := []FilterTag{}

  
//line parse.go:32
	{
	cs = syslog_rfc5424_start
	}

//line parse.go:37
	{
	if p == pe {
		goto _test_eof
	}
	switch cs {
	case 1:
		goto st_case_1
	case 0:
		goto st_case_0
	case 2:
		goto st_case_2
	case 3:
		goto st_case_3
	}
	goto st_out
	st_case_1:
		switch data[p] {
		case 42:
			goto tr0
		case 97:
			goto tr2
		case 110:
			goto tr3
		case 114:
			goto tr4
		case 119:
			goto tr5
		}
		goto st0
st_case_0:
	st0:
		cs = 0
		goto _out
tr0:
//line parse.rl:27
 foundTypes = append(foundTypes, NodeFilter, AreaFilter, WayFilter, RelationFilter) 
	goto st2
	st2:
		if p++; p == pe {
			goto _test_eof2
		}
	st_case_2:
//line parse.go:80
		goto st0
tr2:
//line parse.rl:22
 foundTypes = append(foundTypes, AreaFilter) 
	goto st3
tr3:
//line parse.rl:23
 foundTypes = append(foundTypes, NodeFilter) 
	goto st3
tr4:
//line parse.rl:24
 foundTypes = append(foundTypes, RelationFilter) 
	goto st3
tr5:
//line parse.rl:25
 foundTypes = append(foundTypes, WayFilter) 
	goto st3
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
//line parse.go:103
		switch data[p] {
		case 97:
			goto tr2
		case 110:
			goto tr3
		case 114:
			goto tr4
		case 119:
			goto tr5
		}
		goto st0
	st_out:
	_test_eof2: cs = 2; goto _test_eof
	_test_eof3: cs = 3; goto _test_eof

	_test_eof: {}
	_out: {}
	}

//line parse.rl:36


  sort.Slice(foundTypes, func(i, j int) bool {
		return foundTypes[i] < foundTypes[j]
	})

	foundTypes = lo.Uniq(foundTypes)

	return &AST{
		Types: foundTypes,
		Tags:  tags,
	}, nil
}