//line parse.rl:1
package query

import (
	"fmt"
	"sort"

	"github.com/samber/lo"
)

//line parse.go:14
const query_start int = 1
const query_first_final int = 65
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
		case 65:
			goto st_case_65
		case 2:
			goto st_case_2
		case 3:
			goto st_case_3
		case 66:
			goto st_case_66
		case 4:
			goto st_case_4
		case 5:
			goto st_case_5
		case 67:
			goto st_case_67
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
		case 68:
			goto st_case_68
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
		case 31:
			goto st_case_31
		case 32:
			goto st_case_32
		case 33:
			goto st_case_33
		case 34:
			goto st_case_34
		case 35:
			goto st_case_35
		case 36:
			goto st_case_36
		case 37:
			goto st_case_37
		case 38:
			goto st_case_38
		case 39:
			goto st_case_39
		case 40:
			goto st_case_40
		case 41:
			goto st_case_41
		case 42:
			goto st_case_42
		case 43:
			goto st_case_43
		case 44:
			goto st_case_44
		case 45:
			goto st_case_45
		case 46:
			goto st_case_46
		case 47:
			goto st_case_47
		case 48:
			goto st_case_48
		case 49:
			goto st_case_49
		case 50:
			goto st_case_50
		case 51:
			goto st_case_51
		case 52:
			goto st_case_52
		case 53:
			goto st_case_53
		case 54:
			goto st_case_54
		case 55:
			goto st_case_55
		case 56:
			goto st_case_56
		case 57:
			goto st_case_57
		case 58:
			goto st_case_58
		case 59:
			goto st_case_59
		case 60:
			goto st_case_60
		case 61:
			goto st_case_61
		case 62:
			goto st_case_62
		case 63:
			goto st_case_63
		case 64:
			goto st_case_64
		case 69:
			goto st_case_69
		}
		goto st_out
	st_case_1:
		switch data[p] {
		case 42:
			goto tr0
		case 110:
			goto tr2
		case 114:
			goto tr3
		case 119:
			goto tr4
		}
		goto st0
	tr127:
//line parse.rl:44

		return nil, fmt.Errorf("an undefined type was specified %c: %w", data[p], ErrUndefinedFilter)

		goto st0
//line parse.go:219
	st_case_0:
	st0:
		cs = 0
		goto _out
	tr0:
//line parse.rl:43
		foundTypes = append(foundTypes, NodeFilter, WayFilter, RelationFilter)
		goto st65
	st65:
		if p++; p == pe {
			goto _test_eof65
		}
	st_case_65:
//line parse.go:233
		switch data[p] {
		case 40:
			goto tr128
		case 91:
			goto tr129
		}
		goto tr127
	tr128:
//line parse.rl:65
		directive = FilterDirective{}
		goto st2
	tr130:
//line parse.rl:69
		brackets--
//line parse.rl:66
		directives[directiveName] = directive
//line parse.rl:65
		directive = FilterDirective{}
		goto st2
	tr132:
//line parse.rl:69
		brackets--
//line parse.rl:48
		tags = append(tags, tag)
//line parse.rl:65
		directive = FilterDirective{}
		goto st2
	st2:
		if p++; p == pe {
			goto _test_eof2
		}
	st_case_2:
//line parse.go:266
		if data[p] == 42 {
			goto tr5
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr6
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr6
			}
		default:
			goto tr6
		}
		goto st0
	tr5:
//line parse.rl:68
		brackets++
		goto st3
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
//line parse.go:292
		switch data[p] {
		case 41:
			goto st66
		case 61:
			goto st4
		}
		goto st0
	tr24:
//line parse.rl:62
		directiveName = data[mark:p]
		goto st66
	st66:
		if p++; p == pe {
			goto _test_eof66
		}
	st_case_66:
//line parse.go:309
		if data[p] == 40 {
			goto tr130
		}
		goto st0
	tr13:
//line parse.rl:63
		directive = append(directive, data[mark:p])
		goto st4
	tr26:
//line parse.rl:62
		directiveName = data[mark:p]
		goto st4
	st4:
		if p++; p == pe {
			goto _test_eof4
		}
	st_case_4:
//line parse.go:327
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
//line parse.go:346
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
//line parse.rl:63
		directive = append(directive, data[mark:p])
		goto st67
	tr17:
//line parse.rl:38
		mark = p
//line parse.rl:63
		directive = append(directive, data[mark:p])
		goto st67
	tr21:
//line parse.rl:62
		directiveName = data[mark:p]
//line parse.rl:63
		directive = append(directive, data[mark:p])
		goto st67
	st67:
		if p++; p == pe {
			goto _test_eof67
		}
	st_case_67:
//line parse.go:379
		switch data[p] {
		case 34:
			goto st0
		case 40:
			goto tr131
		case 41:
			goto tr12
		case 44:
			goto tr13
		case 93:
			goto st0
		}
		goto st5
	tr131:
//line parse.rl:69
		brackets--
//line parse.rl:66
		directives[directiveName] = directive
//line parse.rl:65
		directive = FilterDirective{}
		goto st6
	st6:
		if p++; p == pe {
			goto _test_eof6
		}
	st_case_6:
//line parse.go:406
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
//line parse.rl:68
		brackets++
		goto st7
	st7:
		if p++; p == pe {
			goto _test_eof7
		}
	st_case_7:
//line parse.go:441
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
//line parse.rl:62
		directiveName = data[mark:p]
		goto st8
	st8:
		if p++; p == pe {
			goto _test_eof8
		}
	st_case_8:
//line parse.go:464
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
//line parse.go:494
		if data[p] == 34 {
			goto tr20
		}
		goto st10
	tr20:
//line parse.rl:63
		directive = append(directive, data[mark:p])
		goto st11
	st11:
		if p++; p == pe {
			goto _test_eof11
		}
	st_case_11:
//line parse.go:508
		switch data[p] {
		case 41:
			goto st66
		case 44:
			goto st4
		}
		goto st0
	tr15:
//line parse.rl:68
		brackets++
//line parse.rl:38
		mark = p
		goto st12
	st12:
		if p++; p == pe {
			goto _test_eof12
		}
	st_case_12:
//line parse.go:527
		switch data[p] {
		case 34:
			goto st0
		case 41:
			goto tr21
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
	tr6:
//line parse.rl:68
		brackets++
//line parse.rl:38
		mark = p
		goto st13
	st13:
		if p++; p == pe {
			goto _test_eof13
		}
	st_case_13:
//line parse.go:564
		switch data[p] {
		case 41:
			goto tr24
		case 61:
			goto tr26
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
	tr129:
//line parse.rl:47
		tag = FilterTag{Lookups: []string{}}
		goto st14
	tr133:
//line parse.rl:69
		brackets--
//line parse.rl:48
		tags = append(tags, tag)
//line parse.rl:47
		tag = FilterTag{Lookups: []string{}}
		goto st14
	st14:
		if p++; p == pe {
			goto _test_eof14
		}
	st_case_14:
//line parse.go:601
		switch data[p] {
		case 33:
			goto tr27
		case 42:
			goto tr28
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr29
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr29
			}
		default:
			goto tr29
		}
		goto st0
	tr27:
//line parse.rl:68
		brackets++
		goto st15
	st15:
		if p++; p == pe {
			goto _test_eof15
		}
	st_case_15:
//line parse.go:630
		if data[p] == 42 {
			goto st16
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr31
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr31
			}
		default:
			goto tr31
		}
		goto st0
	st16:
		if p++; p == pe {
			goto _test_eof16
		}
	st_case_16:
		if data[p] == 93 {
			goto tr32
		}
		goto st0
	tr32:
//line parse.rl:53
		tag.Op = OpNotExists
		goto st68
	tr34:
//line parse.rl:59
		tag.Name = data[mark:p]
//line parse.rl:53
		tag.Op = OpNotExists
		goto st68
	tr39:
//line parse.rl:52
		tag.Op = OpExists
		goto st68
	tr45:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
//line parse.rl:51
		tag.Op = OpNotEquals
		goto st68
	tr49:
//line parse.rl:51
		tag.Op = OpNotEquals
		goto st68
	tr55:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
//line parse.rl:56
		tag.Op = OpLessThan
		goto st68
	tr60:
//line parse.rl:56
		tag.Op = OpLessThan
		goto st68
	tr65:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
//line parse.rl:56
		tag.Op = OpLessThan
//line parse.rl:57
		tag.Op = OpLessThanEquals
		goto st68
	tr71:
//line parse.rl:56
		tag.Op = OpLessThan
//line parse.rl:57
		tag.Op = OpLessThanEquals
		goto st68
	tr76:
//line parse.rl:57
		tag.Op = OpLessThanEquals
		goto st68
	tr80:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
//line parse.rl:57
		tag.Op = OpLessThanEquals
		goto st68
	tr85:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
//line parse.rl:50
		tag.Op = OpEquals
		goto st68
	tr89:
//line parse.rl:50
		tag.Op = OpEquals
		goto st68
	tr95:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
//line parse.rl:54
		tag.Op = OpGreaterThan
		goto st68
	tr100:
//line parse.rl:54
		tag.Op = OpGreaterThan
		goto st68
	tr105:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
//line parse.rl:54
		tag.Op = OpGreaterThan
//line parse.rl:55
		tag.Op = OpGreaterThanEquals
		goto st68
	tr111:
//line parse.rl:54
		tag.Op = OpGreaterThan
//line parse.rl:55
		tag.Op = OpGreaterThanEquals
		goto st68
	tr116:
//line parse.rl:55
		tag.Op = OpGreaterThanEquals
		goto st68
	tr120:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
//line parse.rl:55
		tag.Op = OpGreaterThanEquals
		goto st68
	tr126:
//line parse.rl:59
		tag.Name = data[mark:p]
//line parse.rl:52
		tag.Op = OpExists
		goto st68
	st68:
		if p++; p == pe {
			goto _test_eof68
		}
	st_case_68:
//line parse.go:769
		switch data[p] {
		case 40:
			goto tr132
		case 91:
			goto tr133
		}
		goto st0
	tr31:
//line parse.rl:38
		mark = p
		goto st17
	st17:
		if p++; p == pe {
			goto _test_eof17
		}
	st_case_17:
//line parse.go:786
		if data[p] == 93 {
			goto tr34
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
	tr28:
//line parse.rl:68
		brackets++
		goto st18
	st18:
		if p++; p == pe {
			goto _test_eof18
		}
	st_case_18:
//line parse.go:812
		switch data[p] {
		case 33:
			goto st19
		case 60:
			goto st25
		case 61:
			goto st42
		case 62:
			goto st47
		case 93:
			goto tr39
		}
		goto st0
	tr121:
//line parse.rl:59
		tag.Name = data[mark:p]
		goto st19
	st19:
		if p++; p == pe {
			goto _test_eof19
		}
	st_case_19:
//line parse.go:835
		if data[p] == 61 {
			goto st20
		}
		goto st0
	tr44:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st20
	st20:
		if p++; p == pe {
			goto _test_eof20
		}
	st_case_20:
//line parse.go:849
		switch data[p] {
		case 34:
			goto st22
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr41
	tr41:
//line parse.rl:38
		mark = p
		goto st21
	st21:
		if p++; p == pe {
			goto _test_eof21
		}
	st_case_21:
//line parse.go:868
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr44
		case 93:
			goto tr45
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
		goto tr46
	tr46:
//line parse.rl:38
		mark = p
		goto st23
	st23:
		if p++; p == pe {
			goto _test_eof23
		}
	st_case_23:
//line parse.go:896
		if data[p] == 34 {
			goto tr48
		}
		goto st23
	tr48:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st24
	st24:
		if p++; p == pe {
			goto _test_eof24
		}
	st_case_24:
//line parse.go:910
		switch data[p] {
		case 44:
			goto st20
		case 93:
			goto tr49
		}
		goto st0
	tr123:
//line parse.rl:59
		tag.Name = data[mark:p]
		goto st25
	st25:
		if p++; p == pe {
			goto _test_eof25
		}
	st_case_25:
//line parse.go:927
		switch data[p] {
		case 34:
			goto st28
		case 44:
			goto st0
		case 61:
			goto tr52
		case 93:
			goto st0
		}
		goto tr50
	tr50:
//line parse.rl:38
		mark = p
		goto st26
	st26:
		if p++; p == pe {
			goto _test_eof26
		}
	st_case_26:
//line parse.go:948
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr54
		case 93:
			goto tr55
		}
		goto st26
	tr54:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st27
	st27:
		if p++; p == pe {
			goto _test_eof27
		}
	st_case_27:
//line parse.go:967
		switch data[p] {
		case 34:
			goto st28
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr50
	st28:
		if p++; p == pe {
			goto _test_eof28
		}
	st_case_28:
		if data[p] == 34 {
			goto st0
		}
		goto tr56
	tr56:
//line parse.rl:38
		mark = p
		goto st29
	st29:
		if p++; p == pe {
			goto _test_eof29
		}
	st_case_29:
//line parse.go:995
		if data[p] == 34 {
			goto tr58
		}
		goto st29
	tr58:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st30
	st30:
		if p++; p == pe {
			goto _test_eof30
		}
	st_case_30:
//line parse.go:1009
		switch data[p] {
		case 44:
			goto st27
		case 93:
			goto tr60
		}
		goto st0
	tr52:
//line parse.rl:38
		mark = p
		goto st31
	st31:
		if p++; p == pe {
			goto _test_eof31
		}
	st_case_31:
//line parse.go:1026
		switch data[p] {
		case 34:
			goto st37
		case 44:
			goto tr54
		case 93:
			goto tr55
		}
		goto tr61
	tr61:
//line parse.rl:38
		mark = p
		goto st32
	st32:
		if p++; p == pe {
			goto _test_eof32
		}
	st_case_32:
//line parse.go:1045
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr64
		case 93:
			goto tr65
		}
		goto st32
	tr64:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st33
	st33:
		if p++; p == pe {
			goto _test_eof33
		}
	st_case_33:
//line parse.go:1064
		switch data[p] {
		case 34:
			goto st34
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr61
	st34:
		if p++; p == pe {
			goto _test_eof34
		}
	st_case_34:
		if data[p] == 34 {
			goto st0
		}
		goto tr67
	tr67:
//line parse.rl:38
		mark = p
		goto st35
	st35:
		if p++; p == pe {
			goto _test_eof35
		}
	st_case_35:
//line parse.go:1092
		if data[p] == 34 {
			goto tr69
		}
		goto st35
	tr69:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st36
	st36:
		if p++; p == pe {
			goto _test_eof36
		}
	st_case_36:
//line parse.go:1106
		switch data[p] {
		case 44:
			goto st33
		case 93:
			goto tr71
		}
		goto st0
	st37:
		if p++; p == pe {
			goto _test_eof37
		}
	st_case_37:
		if data[p] == 34 {
			goto st0
		}
		goto tr72
	tr72:
//line parse.rl:38
		mark = p
		goto st38
	st38:
		if p++; p == pe {
			goto _test_eof38
		}
	st_case_38:
//line parse.go:1132
		if data[p] == 34 {
			goto tr74
		}
		goto st38
	tr74:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st39
	st39:
		if p++; p == pe {
			goto _test_eof39
		}
	st_case_39:
//line parse.go:1146
		switch data[p] {
		case 44:
			goto st40
		case 93:
			goto tr76
		}
		goto st0
	tr79:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st40
	st40:
		if p++; p == pe {
			goto _test_eof40
		}
	st_case_40:
//line parse.go:1163
		switch data[p] {
		case 34:
			goto st37
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr77
	tr77:
//line parse.rl:38
		mark = p
		goto st41
	st41:
		if p++; p == pe {
			goto _test_eof41
		}
	st_case_41:
//line parse.go:1182
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr79
		case 93:
			goto tr80
		}
		goto st41
	tr84:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st42
	tr124:
//line parse.rl:59
		tag.Name = data[mark:p]
		goto st42
	st42:
		if p++; p == pe {
			goto _test_eof42
		}
	st_case_42:
//line parse.go:1205
		switch data[p] {
		case 34:
			goto st44
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr81
	tr81:
//line parse.rl:38
		mark = p
		goto st43
	st43:
		if p++; p == pe {
			goto _test_eof43
		}
	st_case_43:
//line parse.go:1224
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr84
		case 93:
			goto tr85
		}
		goto st43
	st44:
		if p++; p == pe {
			goto _test_eof44
		}
	st_case_44:
		if data[p] == 34 {
			goto st0
		}
		goto tr86
	tr86:
//line parse.rl:38
		mark = p
		goto st45
	st45:
		if p++; p == pe {
			goto _test_eof45
		}
	st_case_45:
//line parse.go:1252
		if data[p] == 34 {
			goto tr88
		}
		goto st45
	tr88:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st46
	st46:
		if p++; p == pe {
			goto _test_eof46
		}
	st_case_46:
//line parse.go:1266
		switch data[p] {
		case 44:
			goto st42
		case 93:
			goto tr89
		}
		goto st0
	tr125:
//line parse.rl:59
		tag.Name = data[mark:p]
		goto st47
	st47:
		if p++; p == pe {
			goto _test_eof47
		}
	st_case_47:
//line parse.go:1283
		switch data[p] {
		case 34:
			goto st50
		case 44:
			goto st0
		case 61:
			goto tr92
		case 93:
			goto st0
		}
		goto tr90
	tr90:
//line parse.rl:38
		mark = p
		goto st48
	st48:
		if p++; p == pe {
			goto _test_eof48
		}
	st_case_48:
//line parse.go:1304
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr94
		case 93:
			goto tr95
		}
		goto st48
	tr94:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st49
	st49:
		if p++; p == pe {
			goto _test_eof49
		}
	st_case_49:
//line parse.go:1323
		switch data[p] {
		case 34:
			goto st50
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr90
	st50:
		if p++; p == pe {
			goto _test_eof50
		}
	st_case_50:
		if data[p] == 34 {
			goto st0
		}
		goto tr96
	tr96:
//line parse.rl:38
		mark = p
		goto st51
	st51:
		if p++; p == pe {
			goto _test_eof51
		}
	st_case_51:
//line parse.go:1351
		if data[p] == 34 {
			goto tr98
		}
		goto st51
	tr98:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st52
	st52:
		if p++; p == pe {
			goto _test_eof52
		}
	st_case_52:
//line parse.go:1365
		switch data[p] {
		case 44:
			goto st49
		case 93:
			goto tr100
		}
		goto st0
	tr92:
//line parse.rl:38
		mark = p
		goto st53
	st53:
		if p++; p == pe {
			goto _test_eof53
		}
	st_case_53:
//line parse.go:1382
		switch data[p] {
		case 34:
			goto st59
		case 44:
			goto tr94
		case 93:
			goto tr95
		}
		goto tr101
	tr101:
//line parse.rl:38
		mark = p
		goto st54
	st54:
		if p++; p == pe {
			goto _test_eof54
		}
	st_case_54:
//line parse.go:1401
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr104
		case 93:
			goto tr105
		}
		goto st54
	tr104:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st55
	st55:
		if p++; p == pe {
			goto _test_eof55
		}
	st_case_55:
//line parse.go:1420
		switch data[p] {
		case 34:
			goto st56
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr101
	st56:
		if p++; p == pe {
			goto _test_eof56
		}
	st_case_56:
		if data[p] == 34 {
			goto st0
		}
		goto tr107
	tr107:
//line parse.rl:38
		mark = p
		goto st57
	st57:
		if p++; p == pe {
			goto _test_eof57
		}
	st_case_57:
//line parse.go:1448
		if data[p] == 34 {
			goto tr109
		}
		goto st57
	tr109:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st58
	st58:
		if p++; p == pe {
			goto _test_eof58
		}
	st_case_58:
//line parse.go:1462
		switch data[p] {
		case 44:
			goto st55
		case 93:
			goto tr111
		}
		goto st0
	st59:
		if p++; p == pe {
			goto _test_eof59
		}
	st_case_59:
		if data[p] == 34 {
			goto st0
		}
		goto tr112
	tr112:
//line parse.rl:38
		mark = p
		goto st60
	st60:
		if p++; p == pe {
			goto _test_eof60
		}
	st_case_60:
//line parse.go:1488
		if data[p] == 34 {
			goto tr114
		}
		goto st60
	tr114:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st61
	st61:
		if p++; p == pe {
			goto _test_eof61
		}
	st_case_61:
//line parse.go:1502
		switch data[p] {
		case 44:
			goto st62
		case 93:
			goto tr116
		}
		goto st0
	tr119:
//line parse.rl:60
		tag.Lookups = append(tag.Lookups, data[mark:p])
		goto st62
	st62:
		if p++; p == pe {
			goto _test_eof62
		}
	st_case_62:
//line parse.go:1519
		switch data[p] {
		case 34:
			goto st59
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr117
	tr117:
//line parse.rl:38
		mark = p
		goto st63
	st63:
		if p++; p == pe {
			goto _test_eof63
		}
	st_case_63:
//line parse.go:1538
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr119
		case 93:
			goto tr120
		}
		goto st63
	tr29:
//line parse.rl:68
		brackets++
//line parse.rl:38
		mark = p
		goto st64
	st64:
		if p++; p == pe {
			goto _test_eof64
		}
	st_case_64:
//line parse.go:1559
		switch data[p] {
		case 33:
			goto tr121
		case 60:
			goto tr123
		case 61:
			goto tr124
		case 62:
			goto tr125
		case 93:
			goto tr126
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st64
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st64
			}
		default:
			goto st64
		}
		goto st0
	tr2:
//line parse.rl:40
		foundTypes = append(foundTypes, NodeFilter)
		goto st69
	tr3:
//line parse.rl:41
		foundTypes = append(foundTypes, RelationFilter)
		goto st69
	tr4:
//line parse.rl:42
		foundTypes = append(foundTypes, WayFilter)
		goto st69
	st69:
		if p++; p == pe {
			goto _test_eof69
		}
	st_case_69:
//line parse.go:1602
		switch data[p] {
		case 40:
			goto tr128
		case 91:
			goto tr129
		case 110:
			goto tr2
		case 114:
			goto tr3
		case 119:
			goto tr4
		}
		goto tr127
	st_out:
	_test_eof65:
		cs = 65
		goto _test_eof
	_test_eof2:
		cs = 2
		goto _test_eof
	_test_eof3:
		cs = 3
		goto _test_eof
	_test_eof66:
		cs = 66
		goto _test_eof
	_test_eof4:
		cs = 4
		goto _test_eof
	_test_eof5:
		cs = 5
		goto _test_eof
	_test_eof67:
		cs = 67
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
	_test_eof68:
		cs = 68
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
	_test_eof31:
		cs = 31
		goto _test_eof
	_test_eof32:
		cs = 32
		goto _test_eof
	_test_eof33:
		cs = 33
		goto _test_eof
	_test_eof34:
		cs = 34
		goto _test_eof
	_test_eof35:
		cs = 35
		goto _test_eof
	_test_eof36:
		cs = 36
		goto _test_eof
	_test_eof37:
		cs = 37
		goto _test_eof
	_test_eof38:
		cs = 38
		goto _test_eof
	_test_eof39:
		cs = 39
		goto _test_eof
	_test_eof40:
		cs = 40
		goto _test_eof
	_test_eof41:
		cs = 41
		goto _test_eof
	_test_eof42:
		cs = 42
		goto _test_eof
	_test_eof43:
		cs = 43
		goto _test_eof
	_test_eof44:
		cs = 44
		goto _test_eof
	_test_eof45:
		cs = 45
		goto _test_eof
	_test_eof46:
		cs = 46
		goto _test_eof
	_test_eof47:
		cs = 47
		goto _test_eof
	_test_eof48:
		cs = 48
		goto _test_eof
	_test_eof49:
		cs = 49
		goto _test_eof
	_test_eof50:
		cs = 50
		goto _test_eof
	_test_eof51:
		cs = 51
		goto _test_eof
	_test_eof52:
		cs = 52
		goto _test_eof
	_test_eof53:
		cs = 53
		goto _test_eof
	_test_eof54:
		cs = 54
		goto _test_eof
	_test_eof55:
		cs = 55
		goto _test_eof
	_test_eof56:
		cs = 56
		goto _test_eof
	_test_eof57:
		cs = 57
		goto _test_eof
	_test_eof58:
		cs = 58
		goto _test_eof
	_test_eof59:
		cs = 59
		goto _test_eof
	_test_eof60:
		cs = 60
		goto _test_eof
	_test_eof61:
		cs = 61
		goto _test_eof
	_test_eof62:
		cs = 62
		goto _test_eof
	_test_eof63:
		cs = 63
		goto _test_eof
	_test_eof64:
		cs = 64
		goto _test_eof
	_test_eof69:
		cs = 69
		goto _test_eof

	_test_eof:
		{
		}
		if p == eof {
			switch cs {
			case 68:
//line parse.rl:69
				brackets--
//line parse.rl:48
				tags = append(tags, tag)
			case 66, 67:
//line parse.rl:69
				brackets--
//line parse.rl:66
				directives[directiveName] = directive
//line parse.go:1699
			}
		}

	_out:
		{
		}
	}

//line parse.rl:101

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
