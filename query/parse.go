
//line parse.rl:1
package query

import (
  "sort"
  "fmt"

  "github.com/samber/lo"
)


//line parse.go:14
const query_start int = 1
const query_first_final int = 19
const query_error int = 0

const query_en_main int = 1


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
	cs = query_start
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
	case 19:
		goto st_case_19
	case 2:
		goto st_case_2
	case 3:
		goto st_case_3
	case 4:
		goto st_case_4
	case 20:
		goto st_case_20
	case 5:
		goto st_case_5
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
	case 17:
		goto st_case_17
	case 18:
		goto st_case_18
	case 21:
		goto st_case_21
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
tr40:
//line parse.rl:39

      return nil, fmt.Errorf("an undefined type was specified %c: %w", data[p], ErrUndefinedFilter)
    
	goto st0
//line parse.go:119
st_case_0:
	st0:
		cs = 0
		goto _out
tr0:
//line parse.rl:38
 foundTypes = append(foundTypes, NodeFilter, AreaFilter, WayFilter, RelationFilter) 
	goto st19
	st19:
		if p++; p == pe {
			goto _test_eof19
		}
	st_case_19:
//line parse.go:133
		if data[p] == 91 {
			goto tr41
		}
		goto tr40
tr41:
//line parse.rl:42
 tag = FilterTag{Lookups: []string{}} 
	goto st2
tr42:
//line parse.rl:54
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
//line parse.go:155
		switch data[p] {
		case 33:
			goto tr6
		case 42:
			goto tr7
		}
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
tr6:
//line parse.rl:53
 brackets++ 
	goto st3
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
//line parse.go:184
		if data[p] == 42 {
			goto st4
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr10
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr10
			}
		default:
			goto tr10
		}
		goto st0
	st4:
		if p++; p == pe {
			goto _test_eof4
		}
	st_case_4:
		if data[p] == 93 {
			goto tr11
		}
		goto st0
tr11:
//line parse.rl:48
 tag.Op = OpNotExists 
	goto st20
tr13:
//line parse.rl:50
 tag.Name    = data[mark:p] 
//line parse.rl:48
 tag.Op = OpNotExists 
	goto st20
tr16:
//line parse.rl:47
 tag.Op = OpExists 
	goto st20
tr22:
//line parse.rl:51
 tag.Lookups = append(tag.Lookups, data[mark:p]) 
//line parse.rl:46
 tag.Op = OpNotEquals 
	goto st20
tr26:
//line parse.rl:46
 tag.Op = OpNotEquals 
	goto st20
tr31:
//line parse.rl:51
 tag.Lookups = append(tag.Lookups, data[mark:p]) 
//line parse.rl:45
 tag.Op = OpEquals 
	goto st20
tr35:
//line parse.rl:45
 tag.Op = OpEquals 
	goto st20
tr39:
//line parse.rl:50
 tag.Name    = data[mark:p] 
//line parse.rl:47
 tag.Op = OpExists 
	goto st20
	st20:
		if p++; p == pe {
			goto _test_eof20
		}
	st_case_20:
//line parse.go:255
		if data[p] == 91 {
			goto tr42
		}
		goto st0
tr10:
//line parse.rl:32
 mark = p
	goto st5
	st5:
		if p++; p == pe {
			goto _test_eof5
		}
	st_case_5:
//line parse.go:269
		if data[p] == 93 {
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
tr7:
//line parse.rl:53
 brackets++ 
	goto st6
	st6:
		if p++; p == pe {
			goto _test_eof6
		}
	st_case_6:
//line parse.go:295
		switch data[p] {
		case 33:
			goto st7
		case 61:
			goto st13
		case 93:
			goto tr16
		}
		goto st0
tr36:
//line parse.rl:50
 tag.Name    = data[mark:p] 
	goto st7
	st7:
		if p++; p == pe {
			goto _test_eof7
		}
	st_case_7:
//line parse.go:314
		if data[p] == 61 {
			goto st8
		}
		goto st0
tr21:
//line parse.rl:51
 tag.Lookups = append(tag.Lookups, data[mark:p]) 
	goto st8
	st8:
		if p++; p == pe {
			goto _test_eof8
		}
	st_case_8:
//line parse.go:328
		switch data[p] {
		case 34:
			goto st10
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr18
tr18:
//line parse.rl:32
 mark = p
	goto st9
	st9:
		if p++; p == pe {
			goto _test_eof9
		}
	st_case_9:
//line parse.go:347
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr21
		case 93:
			goto tr22
		}
		goto st9
	st10:
		if p++; p == pe {
			goto _test_eof10
		}
	st_case_10:
		if data[p] == 34 {
			goto st0
		}
		goto tr23
tr23:
//line parse.rl:32
 mark = p
	goto st11
	st11:
		if p++; p == pe {
			goto _test_eof11
		}
	st_case_11:
//line parse.go:375
		if data[p] == 34 {
			goto tr25
		}
		goto st11
tr25:
//line parse.rl:51
 tag.Lookups = append(tag.Lookups, data[mark:p]) 
	goto st12
	st12:
		if p++; p == pe {
			goto _test_eof12
		}
	st_case_12:
//line parse.go:389
		switch data[p] {
		case 44:
			goto st8
		case 93:
			goto tr26
		}
		goto st0
tr30:
//line parse.rl:51
 tag.Lookups = append(tag.Lookups, data[mark:p]) 
	goto st13
tr38:
//line parse.rl:50
 tag.Name    = data[mark:p] 
	goto st13
	st13:
		if p++; p == pe {
			goto _test_eof13
		}
	st_case_13:
//line parse.go:410
		switch data[p] {
		case 34:
			goto st15
		case 44:
			goto st0
		case 93:
			goto st0
		}
		goto tr27
tr27:
//line parse.rl:32
 mark = p
	goto st14
	st14:
		if p++; p == pe {
			goto _test_eof14
		}
	st_case_14:
//line parse.go:429
		switch data[p] {
		case 34:
			goto st0
		case 44:
			goto tr30
		case 93:
			goto tr31
		}
		goto st14
	st15:
		if p++; p == pe {
			goto _test_eof15
		}
	st_case_15:
		if data[p] == 34 {
			goto st0
		}
		goto tr32
tr32:
//line parse.rl:32
 mark = p
	goto st16
	st16:
		if p++; p == pe {
			goto _test_eof16
		}
	st_case_16:
//line parse.go:457
		if data[p] == 34 {
			goto tr34
		}
		goto st16
tr34:
//line parse.rl:51
 tag.Lookups = append(tag.Lookups, data[mark:p]) 
	goto st17
	st17:
		if p++; p == pe {
			goto _test_eof17
		}
	st_case_17:
//line parse.go:471
		switch data[p] {
		case 44:
			goto st13
		case 93:
			goto tr35
		}
		goto st0
tr8:
//line parse.rl:53
 brackets++ 
//line parse.rl:32
 mark = p
	goto st18
	st18:
		if p++; p == pe {
			goto _test_eof18
		}
	st_case_18:
//line parse.go:490
		switch data[p] {
		case 33:
			goto tr36
		case 61:
			goto tr38
		case 93:
			goto tr39
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st18
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st18
			}
		default:
			goto st18
		}
		goto st0
tr2:
//line parse.rl:34
 foundTypes = append(foundTypes, AreaFilter) 
	goto st21
tr3:
//line parse.rl:35
 foundTypes = append(foundTypes, NodeFilter) 
	goto st21
tr4:
//line parse.rl:36
 foundTypes = append(foundTypes, RelationFilter) 
	goto st21
tr5:
//line parse.rl:37
 foundTypes = append(foundTypes, WayFilter) 
	goto st21
	st21:
		if p++; p == pe {
			goto _test_eof21
		}
	st_case_21:
//line parse.go:533
		switch data[p] {
		case 91:
			goto tr41
		case 97:
			goto tr2
		case 110:
			goto tr3
		case 114:
			goto tr4
		case 119:
			goto tr5
		}
		goto tr40
	st_out:
	_test_eof19: cs = 19; goto _test_eof
	_test_eof2: cs = 2; goto _test_eof
	_test_eof3: cs = 3; goto _test_eof
	_test_eof4: cs = 4; goto _test_eof
	_test_eof20: cs = 20; goto _test_eof
	_test_eof5: cs = 5; goto _test_eof
	_test_eof6: cs = 6; goto _test_eof
	_test_eof7: cs = 7; goto _test_eof
	_test_eof8: cs = 8; goto _test_eof
	_test_eof9: cs = 9; goto _test_eof
	_test_eof10: cs = 10; goto _test_eof
	_test_eof11: cs = 11; goto _test_eof
	_test_eof12: cs = 12; goto _test_eof
	_test_eof13: cs = 13; goto _test_eof
	_test_eof14: cs = 14; goto _test_eof
	_test_eof15: cs = 15; goto _test_eof
	_test_eof16: cs = 16; goto _test_eof
	_test_eof17: cs = 17; goto _test_eof
	_test_eof18: cs = 18; goto _test_eof
	_test_eof21: cs = 21; goto _test_eof

	_test_eof: {}
	if p == eof {
		switch cs {
		case 20:
//line parse.rl:54
 brackets-- 
//line parse.rl:43
 tags = append(tags, tag) 
//line parse.go:577
		}
	}

	_out: {}
	}

//line parse.rl:78


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
		Types: foundTypes,
		Tags:  tags,
	}, nil
}