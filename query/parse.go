//line parse.rl:1
package query

import (
	"fmt"
	"sort"

	"github.com/samber/lo"
)

//line parse.go:14
const query_start int = 1
const query_first_final int = 31
const query_error int = 0

const query_en_main int = 1

//line parse.rl:13

func Parse(data string) (*AST, error) {
	// types used for the AST
	foundTypes := []FilterType{}
	tags := []FilterTag{}
	directives := map[string]FilterDirective{}

	var (
		tag           FilterTag
		directiveName string
		directive     FilterDirective
	)

	// set defaults for state machine parsing
	cs, p, pe, eof := 0, 0, len(data), len(data)

	// tracks where the beginning of a word starts
	mark := 0

	// keep track of opening and closing brackets
	brackets := 0

//line parse.go:48
	{
		cs = query_start
	}

//line parse.go:53
	{
		if p == pe {
			goto _test_eof
		}
		switch cs {
		case 1:
			goto st_case_1
		case 0:
			goto st_case_0
		case 31:
			goto st_case_31
		case 2:
			goto st_case_2
		case 3:
			goto st_case_3
		case 4:
			goto st_case_4
		case 5:
			goto st_case_5
		case 32:
			goto st_case_32
		case 6:
			goto st_case_6
		case 7:
			goto st_case_7
		case 8:
			goto st_case_8
		case 9:
			goto st_case_9
		case 10:
			goto st_case_10
		case 11:
			goto st_case_11
		case 33:
			goto st_case_33
		case 12:
			goto st_case_12
		case 13:
			goto st_case_13
		case 14:
			goto st_case_14
		case 15:
			goto st_case_15
		case 16:
			goto st_case_16
		case 34:
			goto st_case_34
		case 17:
			goto st_case_17
		case 18:
			goto st_case_18
		case 19:
			goto st_case_19
		case 20:
			goto st_case_20
		case 21:
			goto st_case_21
		case 22:
			goto st_case_22
		case 23:
			goto st_case_23
		case 24:
			goto st_case_24
		case 25:
			goto st_case_25
		case 26:
			goto st_case_26
		case 27:
			goto st_case_27
		case 28:
			goto st_case_28
		case 29:
			goto st_case_29
		case 30:
			goto st_case_30
		case 35:
			goto st_case_35
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
	tr60:
//line parse.rl:45

		return nil, fmt.Errorf("an undefined type was specified %c: %w", data[p], ErrUndefinedFilter)

		goto st0
//line parse.go:153
	st_case_0:
	st0:
		cs = 0
		goto _out
	tr0:
//line parse.rl:44
		foundTypes = append(foundTypes, NodeFilter, AreaFilter, WayFilter, RelationFilter)
		goto st31
	st31:
		if p++; p == pe {
			goto _test_eof31
		}
	st_case_31:
//line parse.go:167
		switch data[p] {
		case 40:
			goto tr61
		case 91:
			goto tr62
		}
		goto tr60
	tr61:
//line parse.rl:62
		directive = FilterDirective{}
		goto st2
	tr64:
//line parse.rl:66
		brackets--
//line parse.rl:63
		directives[directiveName] = directive
//line parse.rl:62
		directive = FilterDirective{}
		goto st2
	tr65:
//line parse.rl:66
		brackets--
//line parse.rl:49
		tags = append(tags, tag)
//line parse.rl:62
		directive = FilterDirective{}
		goto st2
	st2:
		if p++; p == pe {
			goto _test_eof2
		}
	st_case_2:
//line parse.go:200
		if data[p] == 42 {
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
//line parse.rl:65
		brackets++
		goto st3
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
//line parse.go:226
		if data[p] == 61 {
			goto st4
		}
		goto st0
	tr13:
//line parse.rl:60
		directive = append(directive, data[mark:p])
		goto st4
	tr25:
//line parse.rl:59
		directiveName = data[mark:p]
		goto st4
	st4:
		if p++; p == pe {
			goto _test_eof4
		}
	st_case_4:
//line parse.go:244
		switch data[p] {
		case 34:
			goto st9
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr9
	tr9:
//line parse.rl:38
		mark = p
		goto st5
	st5:
		if p++; p == pe {
			goto _test_eof5
		}
	st_case_5:
//line parse.go:263
		switch data[p] {
		case 34:
			goto st0
		case 41:
			goto tr12
		case 44:
			goto tr13
		case 93:
			goto st0
		}
		goto st5
	tr12:
//line parse.rl:60
		directive = append(directive, data[mark:p])
		goto st32
	tr17:
//line parse.rl:38
		mark = p
//line parse.rl:60
		directive = append(directive, data[mark:p])
		goto st32
	st32:
		if p++; p == pe {
			goto _test_eof32
		}
	st_case_32:
//line parse.go:290
		switch data[p] {
		case 34:
			goto st0
		case 40:
			goto tr63
		case 41:
			goto tr12
		case 44:
			goto tr13
		case 93:
			goto st0
		}
		goto st5
	tr63:
//line parse.rl:66
		brackets--
//line parse.rl:63
		directives[directiveName] = directive
//line parse.rl:62
		directive = FilterDirective{}
		goto st6
	st6:
		if p++; p == pe {
			goto _test_eof6
		}
	st_case_6:
//line parse.go:317
		switch data[p] {
		case 34:
			goto st0
		case 41:
			goto tr12
		case 42:
			goto tr14
		case 44:
			goto tr13
		case 93:
			goto st0
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr15
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr15
			}
		default:
			goto tr15
		}
		goto st5
	tr14:
//line parse.rl:65
		brackets++
		goto st7
	st7:
		if p++; p == pe {
			goto _test_eof7
		}
	st_case_7:
//line parse.go:352
		switch data[p] {
		case 34:
			goto st0
		case 41:
			goto tr12
		case 44:
			goto tr13
		case 61:
			goto st8
		case 93:
			goto st0
		}
		goto st5
	tr23:
//line parse.rl:59
		directiveName = data[mark:p]
		goto st8
	st8:
		if p++; p == pe {
			goto _test_eof8
		}
	st_case_8:
//line parse.go:375
		switch data[p] {
		case 34:
			goto st9
		case 41:
			goto tr17
		case 44:
			goto tr13
		case 93:
			goto st0
		}
		goto tr9
	st9:
		if p++; p == pe {
			goto _test_eof9
		}
	st_case_9:
		if data[p] == 34 {
			goto st0
		}
		goto tr18
	tr18:
//line parse.rl:38
		mark = p
		goto st10
	st10:
		if p++; p == pe {
			goto _test_eof10
		}
	st_case_10:
//line parse.go:405
		if data[p] == 34 {
			goto tr20
		}
		goto st10
	tr20:
//line parse.rl:60
		directive = append(directive, data[mark:p])
		goto st11
	st11:
		if p++; p == pe {
			goto _test_eof11
		}
	st_case_11:
//line parse.go:419
		switch data[p] {
		case 41:
			goto st33
		case 44:
			goto st4
		}
		goto st0
	st33:
		if p++; p == pe {
			goto _test_eof33
		}
	st_case_33:
		if data[p] == 40 {
			goto tr64
		}
		goto st0
	tr15:
//line parse.rl:65
		brackets++
//line parse.rl:38
		mark = p
		goto st12
	st12:
		if p++; p == pe {
			goto _test_eof12
		}
	st_case_12:
//line parse.go:447
		switch data[p] {
		case 34:
			goto st0
		case 41:
			goto tr12
		case 44:
			goto tr13
		case 61:
			goto tr23
		case 93:
			goto st0
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st12
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st12
			}
		default:
			goto st12
		}
		goto st5
	tr7:
//line parse.rl:65
		brackets++
//line parse.rl:38
		mark = p
		goto st13
	st13:
		if p++; p == pe {
			goto _test_eof13
		}
	st_case_13:
//line parse.go:484
		if data[p] == 61 {
			goto tr25
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st13
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st13
			}
		default:
			goto st13
		}
		goto st0
	tr62:
//line parse.rl:48
		tag = FilterTag{Lookups: []string{}}
		goto st14
	tr66:
//line parse.rl:66
		brackets--
//line parse.rl:49
		tags = append(tags, tag)
//line parse.rl:48
		tag = FilterTag{Lookups: []string{}}
		goto st14
	st14:
		if p++; p == pe {
			goto _test_eof14
		}
	st_case_14:
//line parse.go:518
		switch data[p] {
		case 33:
			goto tr26
		case 42:
			goto tr27
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr28
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr28
			}
		default:
			goto tr28
		}
		goto st0
	tr26:
//line parse.rl:65
		brackets++
		goto st15
	st15:
		if p++; p == pe {
			goto _test_eof15
		}
	st_case_15:
//line parse.go:547
		if data[p] == 42 {
			goto st16
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr30
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr30
			}
		default:
			goto tr30
		}
		goto st0
	st16:
		if p++; p == pe {
			goto _test_eof16
		}
	st_case_16:
		if data[p] == 93 {
			goto tr31
		}
		goto st0
	tr31:
//line parse.rl:54
		tag.Op = OpNotExists
		goto st34
	tr33:
//line parse.rl:56
		tag.Name = data[mark:p]
//line parse.rl:54
		tag.Op = OpNotExists
		goto st34
	tr36:
//line parse.rl:53
		tag.Op = OpExists
		goto st34
	tr42:
//line parse.rl:57
		tag.Lookups = append(tag.Lookups, data[mark:p])
//line parse.rl:52
		tag.Op = OpNotEquals
		goto st34
	tr46:
//line parse.rl:52
		tag.Op = OpNotEquals
		goto st34
	tr51:
//line parse.rl:57
		tag.Lookups = append(tag.Lookups, data[mark:p])
//line parse.rl:51
		tag.Op = OpEquals
		goto st34
	tr55:
//line parse.rl:51
		tag.Op = OpEquals
		goto st34
	tr59:
//line parse.rl:56
		tag.Name = data[mark:p]
//line parse.rl:53
		tag.Op = OpExists
		goto st34
	st34:
		if p++; p == pe {
			goto _test_eof34
		}
	st_case_34:
//line parse.go:618
		switch data[p] {
		case 40:
			goto tr65
		case 91:
			goto tr66
		}
		goto st0
	tr30:
//line parse.rl:38
		mark = p
		goto st17
	st17:
		if p++; p == pe {
			goto _test_eof17
		}
	st_case_17:
//line parse.go:635
		if data[p] == 93 {
			goto tr33
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st17
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st17
			}
		default:
			goto st17
		}
		goto st0
	tr27:
//line parse.rl:65
		brackets++
		goto st18
	st18:
		if p++; p == pe {
			goto _test_eof18
		}
	st_case_18:
//line parse.go:661
		switch data[p] {
		case 33:
			goto st19
		case 61:
			goto st25
		case 93:
			goto tr36
		}
		goto st0
	tr56:
//line parse.rl:56
		tag.Name = data[mark:p]
		goto st19
	st19:
		if p++; p == pe {
			goto _test_eof19
		}
	st_case_19:
//line parse.go:680
		if data[p] == 61 {
			goto st20
		}
		goto st0
	tr41:
//line parse.rl:57
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st20
	st20:
		if p++; p == pe {
			goto _test_eof20
		}
	st_case_20:
//line parse.go:694
		switch data[p] {
		case 34:
			goto st22
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr38
	tr38:
//line parse.rl:38
		mark = p
		goto st21
	st21:
		if p++; p == pe {
			goto _test_eof21
		}
	st_case_21:
//line parse.go:713
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr41
		case 93:
			goto tr42
		}
		goto st21
	st22:
		if p++; p == pe {
			goto _test_eof22
		}
	st_case_22:
		if data[p] == 34 {
			goto st0
		}
		goto tr43
	tr43:
//line parse.rl:38
		mark = p
		goto st23
	st23:
		if p++; p == pe {
			goto _test_eof23
		}
	st_case_23:
//line parse.go:741
		if data[p] == 34 {
			goto tr45
		}
		goto st23
	tr45:
//line parse.rl:57
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st24
	st24:
		if p++; p == pe {
			goto _test_eof24
		}
	st_case_24:
//line parse.go:755
		switch data[p] {
		case 44:
			goto st20
		case 93:
			goto tr46
		}
		goto st0
	tr50:
//line parse.rl:57
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st25
	tr58:
//line parse.rl:56
		tag.Name = data[mark:p]
		goto st25
	st25:
		if p++; p == pe {
			goto _test_eof25
		}
	st_case_25:
//line parse.go:776
		switch data[p] {
		case 34:
			goto st27
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr47
	tr47:
//line parse.rl:38
		mark = p
		goto st26
	st26:
		if p++; p == pe {
			goto _test_eof26
		}
	st_case_26:
//line parse.go:795
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr50
		case 93:
			goto tr51
		}
		goto st26
	st27:
		if p++; p == pe {
			goto _test_eof27
		}
	st_case_27:
		if data[p] == 34 {
			goto st0
		}
		goto tr52
	tr52:
//line parse.rl:38
		mark = p
		goto st28
	st28:
		if p++; p == pe {
			goto _test_eof28
		}
	st_case_28:
//line parse.go:823
		if data[p] == 34 {
			goto tr54
		}
		goto st28
	tr54:
//line parse.rl:57
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st29
	st29:
		if p++; p == pe {
			goto _test_eof29
		}
	st_case_29:
//line parse.go:837
		switch data[p] {
		case 44:
			goto st25
		case 93:
			goto tr55
		}
		goto st0
	tr28:
//line parse.rl:65
		brackets++
//line parse.rl:38
		mark = p
		goto st30
	st30:
		if p++; p == pe {
			goto _test_eof30
		}
	st_case_30:
//line parse.go:856
		switch data[p] {
		case 33:
			goto tr56
		case 61:
			goto tr58
		case 93:
			goto tr59
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st30
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st30
			}
		default:
			goto st30
		}
		goto st0
	tr2:
//line parse.rl:40
		foundTypes = append(foundTypes, AreaFilter)
		goto st35
	tr3:
//line parse.rl:41
		foundTypes = append(foundTypes, NodeFilter)
		goto st35
	tr4:
//line parse.rl:42
		foundTypes = append(foundTypes, RelationFilter)
		goto st35
	tr5:
//line parse.rl:43
		foundTypes = append(foundTypes, WayFilter)
		goto st35
	st35:
		if p++; p == pe {
			goto _test_eof35
		}
	st_case_35:
//line parse.go:899
		switch data[p] {
		case 40:
			goto tr61
		case 91:
			goto tr62
		case 97:
			goto tr2
		case 110:
			goto tr3
		case 114:
			goto tr4
		case 119:
			goto tr5
		}
		goto tr60
	st_out:
	_test_eof31:
		cs = 31
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
	_test_eof5:
		cs = 5
		goto _test_eof
	_test_eof32:
		cs = 32
		goto _test_eof
	_test_eof6:
		cs = 6
		goto _test_eof
	_test_eof7:
		cs = 7
		goto _test_eof
	_test_eof8:
		cs = 8
		goto _test_eof
	_test_eof9:
		cs = 9
		goto _test_eof
	_test_eof10:
		cs = 10
		goto _test_eof
	_test_eof11:
		cs = 11
		goto _test_eof
	_test_eof33:
		cs = 33
		goto _test_eof
	_test_eof12:
		cs = 12
		goto _test_eof
	_test_eof13:
		cs = 13
		goto _test_eof
	_test_eof14:
		cs = 14
		goto _test_eof
	_test_eof15:
		cs = 15
		goto _test_eof
	_test_eof16:
		cs = 16
		goto _test_eof
	_test_eof34:
		cs = 34
		goto _test_eof
	_test_eof17:
		cs = 17
		goto _test_eof
	_test_eof18:
		cs = 18
		goto _test_eof
	_test_eof19:
		cs = 19
		goto _test_eof
	_test_eof20:
		cs = 20
		goto _test_eof
	_test_eof21:
		cs = 21
		goto _test_eof
	_test_eof22:
		cs = 22
		goto _test_eof
	_test_eof23:
		cs = 23
		goto _test_eof
	_test_eof24:
		cs = 24
		goto _test_eof
	_test_eof25:
		cs = 25
		goto _test_eof
	_test_eof26:
		cs = 26
		goto _test_eof
	_test_eof27:
		cs = 27
		goto _test_eof
	_test_eof28:
		cs = 28
		goto _test_eof
	_test_eof29:
		cs = 29
		goto _test_eof
	_test_eof30:
		cs = 30
		goto _test_eof
	_test_eof35:
		cs = 35
		goto _test_eof

	_test_eof:
		{
		}
		if p == eof {
			switch cs {
			case 34:
//line parse.rl:66
				brackets--
//line parse.rl:49
				tags = append(tags, tag)
			case 32, 33:
//line parse.rl:66
				brackets--
//line parse.rl:63
				directives[directiveName] = directive
//line parse.go:964
			}
		}

	_out:
		{
		}
	}

//line parse.rl:98

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
		Tags:       tags,
		Types:      foundTypes,
		Directives: directives,
	}, nil
}
