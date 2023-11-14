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
	// set defaults for state machine parsing
	cs, p, pe, eof := 0, 0, len(data), len(data)
	foundTypes := []FilterType{}
	tags := []FilterTag{}
	var tag FilterTag
	mark := 0

//line parse.go:34
	{
		cs = syslog_rfc5424_start
	}

//line parse.go:39
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
		case 11:
			goto st_case_11
		case 12:
			goto st_case_12
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
//line parse.rl:31

		return nil, fmt.Errorf("an undefined type was specified %c: %w", data[p], ErrUndefinedFilter)

		goto st0
//line parse.go:93
	st_case_0:
	st0:
		cs = 0
		goto _out
	tr0:
//line parse.rl:30
		foundTypes = append(foundTypes, NodeFilter, AreaFilter, WayFilter, RelationFilter)
		goto st8
	st8:
		if p++; p == pe {
			goto _test_eof8
		}
	st_case_8:
//line parse.go:107
		if data[p] == 91 {
			goto tr19
		}
		goto tr18
	tr19:
//line parse.rl:34
		tag = FilterTag{Lookups: []string{}}
		goto st2
	tr20:
//line parse.rl:39
		tag.Op = OpNotExist
//line parse.rl:35
		tags = append(tags, tag)
//line parse.rl:34
		tag = FilterTag{Lookups: []string{}}
		goto st2
	tr21:
//line parse.rl:37
		tag.Op = OpEquals
//line parse.rl:35
		tags = append(tags, tag)
//line parse.rl:34
		tag = FilterTag{Lookups: []string{}}
		goto st2
	tr22:
//line parse.rl:38
		tag.Op = OpExists
//line parse.rl:35
		tags = append(tags, tag)
//line parse.rl:34
		tag = FilterTag{Lookups: []string{}}
		goto st2
	st2:
		if p++; p == pe {
			goto _test_eof2
		}
	st_case_2:
//line parse.go:145
		if data[p] == 33 {
			goto st3
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
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
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
//line parse.rl:24
		mark = p
		goto st4
	st4:
		if p++; p == pe {
			goto _test_eof4
		}
	st_case_4:
//line parse.go:189
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
//line parse.rl:41
		tag.Name = data[mark:p]
		goto st9
	st9:
		if p++; p == pe {
			goto _test_eof9
		}
	st_case_9:
//line parse.go:215
		if data[p] == 91 {
			goto tr20
		}
		goto st0
	tr7:
//line parse.rl:24
		mark = p
		goto st5
	st5:
		if p++; p == pe {
			goto _test_eof5
		}
	st_case_5:
//line parse.go:229
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
//line parse.rl:41
		tag.Name = data[mark:p]
		goto st6
	tr15:
//line parse.rl:42
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st6
	st6:
		if p++; p == pe {
			goto _test_eof6
		}
	st_case_6:
//line parse.go:262
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
//line parse.rl:24
		mark = p
		goto st7
	st7:
		if p++; p == pe {
			goto _test_eof7
		}
	st_case_7:
//line parse.go:285
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
	tr17:
//line parse.rl:42
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st10
	st10:
		if p++; p == pe {
			goto _test_eof10
		}
	st_case_10:
//line parse.go:314
		if data[p] == 91 {
			goto tr21
		}
		goto st0
	tr13:
//line parse.rl:41
		tag.Name = data[mark:p]
		goto st11
	st11:
		if p++; p == pe {
			goto _test_eof11
		}
	st_case_11:
//line parse.go:328
		if data[p] == 91 {
			goto tr22
		}
		goto st0
	tr2:
//line parse.rl:26
		foundTypes = append(foundTypes, AreaFilter)
		goto st12
	tr3:
//line parse.rl:27
		foundTypes = append(foundTypes, NodeFilter)
		goto st12
	tr4:
//line parse.rl:28
		foundTypes = append(foundTypes, RelationFilter)
		goto st12
	tr5:
//line parse.rl:29
		foundTypes = append(foundTypes, WayFilter)
		goto st12
	st12:
		if p++; p == pe {
			goto _test_eof12
		}
	st_case_12:
//line parse.go:354
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
	_test_eof11:
		cs = 11
		goto _test_eof
	_test_eof12:
		cs = 12
		goto _test_eof

	_test_eof:
		{
		}
		if p == eof {
			switch cs {
			case 10:
//line parse.rl:37
				tag.Op = OpEquals
//line parse.rl:35
				tags = append(tags, tag)
			case 11:
//line parse.rl:38
				tag.Op = OpExists
//line parse.rl:35
				tags = append(tags, tag)
			case 9:
//line parse.rl:39
				tag.Op = OpNotExist
//line parse.rl:35
				tags = append(tags, tag)
//line parse.go:399
			}
		}

	_out:
		{
		}
	}

//line parse.rl:60

	sort.Slice(foundTypes, func(i, j int) bool {
		return foundTypes[i] < foundTypes[j]
	})

	foundTypes = lo.Uniq(foundTypes)

	return &AST{
		Types: foundTypes,
		Tags:  tags,
	}, nil
}
