//line parse.rl:1
package query

import (
	"fmt"
	"sort"

	"github.com/samber/lo"
)

//line parse.go:14
const syslog_rfc5424_start int = 1
const syslog_rfc5424_first_final int = 8
const syslog_rfc5424_error int = 0

const syslog_rfc5424_en_main int = 1

//line parse.rl:13

func Parse(data string) (*AST, error) {
	// types used for the AST
	foundTypes := []FilterType{}
	tags := []FilterTag{}
	var tag FilterTag

	// set defaults for state machine parsing
	cs, p, pe, eof := 0, 0, len(data), len(data)

	// tracks where the beginning of a word starts
	mark := 0

	// keep track of opening and closing brackets
	brackets := 0

//line parse.go:42
	{
		cs = syslog_rfc5424_start
	}

//line parse.go:47
	{
		if p == pe {
			goto _test_eof
		}
		switch cs {
		case 1:
			goto st_case_1
		case 0:
			goto st_case_0
		case 8:
			goto st_case_8
		case 2:
			goto st_case_2
		case 3:
			goto st_case_3
		case 4:
			goto st_case_4
		case 9:
			goto st_case_9
		case 5:
			goto st_case_5
		case 6:
			goto st_case_6
		case 7:
			goto st_case_7
		case 10:
			goto st_case_10
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
	tr18:
//line parse.rl:39

		return nil, fmt.Errorf("an undefined type was specified %c: %w", data[p], ErrUndefinedFilter)

		goto st0
//line parse.go:97
	st_case_0:
	st0:
		cs = 0
		goto _out
	tr0:
//line parse.rl:38
		foundTypes = append(foundTypes, NodeFilter, AreaFilter, WayFilter, RelationFilter)
		goto st8
	st8:
		if p++; p == pe {
			goto _test_eof8
		}
	st_case_8:
//line parse.go:111
		if data[p] == 91 {
			goto tr19
		}
		goto tr18
	tr19:
//line parse.rl:42
		tag = FilterTag{Lookups: []string{}}
		goto st2
	tr20:
//line parse.rl:53
		brackets--
//line parse.rl:43
		tags = append(tags, tag)
//line parse.rl:42
		tag = FilterTag{Lookups: []string{}}
		goto st2
	st2:
		if p++; p == pe {
			goto _test_eof2
		}
	st_case_2:
//line parse.go:133
		if data[p] == 33 {
			goto tr6
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr7
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr7
			}
		default:
			goto tr7
		}
		goto st0
	tr6:
//line parse.rl:52
		brackets++
		goto st3
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
//line parse.go:159
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr8
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr8
			}
		default:
			goto tr8
		}
		goto st0
	tr8:
//line parse.rl:32
		mark = p
		goto st4
	st4:
		if p++; p == pe {
			goto _test_eof4
		}
	st_case_4:
//line parse.go:182
		if data[p] == 93 {
			goto tr10
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st4
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st4
			}
		default:
			goto st4
		}
		goto st0
	tr10:
//line parse.rl:49
		tag.Name = data[mark:p]
//line parse.rl:47
		tag.Op = OpNotExist
		goto st9
	tr13:
//line parse.rl:49
		tag.Name = data[mark:p]
//line parse.rl:46
		tag.Op = OpExists
		goto st9
	tr17:
//line parse.rl:50
		tag.Lookups = append(tag.Lookups, data[mark:p])
//line parse.rl:45
		tag.Op = OpEquals
		goto st9
	st9:
		if p++; p == pe {
			goto _test_eof9
		}
	st_case_9:
//line parse.go:222
		if data[p] == 91 {
			goto tr20
		}
		goto st0
	tr7:
//line parse.rl:52
		brackets++
//line parse.rl:32
		mark = p
		goto st5
	st5:
		if p++; p == pe {
			goto _test_eof5
		}
	st_case_5:
//line parse.go:238
		switch data[p] {
		case 61:
			goto tr12
		case 93:
			goto tr13
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st5
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st5
			}
		default:
			goto st5
		}
		goto st0
	tr12:
//line parse.rl:49
		tag.Name = data[mark:p]
		goto st6
	tr15:
//line parse.rl:50
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st6
	st6:
		if p++; p == pe {
			goto _test_eof6
		}
	st_case_6:
//line parse.go:271
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr14
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr14
			}
		default:
			goto tr14
		}
		goto st0
	tr14:
//line parse.rl:32
		mark = p
		goto st7
	st7:
		if p++; p == pe {
			goto _test_eof7
		}
	st_case_7:
//line parse.go:294
		switch data[p] {
		case 44:
			goto tr15
		case 93:
			goto tr17
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st7
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st7
			}
		default:
			goto st7
		}
		goto st0
	tr2:
//line parse.rl:34
		foundTypes = append(foundTypes, AreaFilter)
		goto st10
	tr3:
//line parse.rl:35
		foundTypes = append(foundTypes, NodeFilter)
		goto st10
	tr4:
//line parse.rl:36
		foundTypes = append(foundTypes, RelationFilter)
		goto st10
	tr5:
//line parse.rl:37
		foundTypes = append(foundTypes, WayFilter)
		goto st10
	st10:
		if p++; p == pe {
			goto _test_eof10
		}
	st_case_10:
//line parse.go:335
		switch data[p] {
		case 91:
			goto tr19
		case 97:
			goto tr2
		case 110:
			goto tr3
		case 114:
			goto tr4
		case 119:
			goto tr5
		}
		goto tr18
	st_out:
	_test_eof8:
		cs = 8
		goto _test_eof
	_test_eof2:
		cs = 2
		goto _test_eof
	_test_eof3:
		cs = 3
		goto _test_eof
	_test_eof4:
		cs = 4
		goto _test_eof
	_test_eof9:
		cs = 9
		goto _test_eof
	_test_eof5:
		cs = 5
		goto _test_eof
	_test_eof6:
		cs = 6
		goto _test_eof
	_test_eof7:
		cs = 7
		goto _test_eof
	_test_eof10:
		cs = 10
		goto _test_eof

	_test_eof:
		{
		}
		if p == eof {
			switch cs {
			case 9:
//line parse.rl:53
				brackets--
//line parse.rl:43
				tags = append(tags, tag)
//line parse.go:368
			}
		}

	_out:
		{
		}
	}

//line parse.rl:71

	if brackets != 0 {
		return nil, fmt.Errorf("tags not enclosed properly (%d): %w", brackets, ErrUnbalancedBrackets)
	}

	sort.Slice(foundTypes, func(i, j int) bool {
		return foundTypes[i] < foundTypes[j]
	})

	foundTypes = lo.Uniq(foundTypes)

	return &AST{
		Types: foundTypes,
		Tags:  tags,
	}, nil
}
