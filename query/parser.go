package query

import (
	"errors"
	"fmt"
	"sort"
)

type FilterType uint

const (
	NodeFilter FilterType = iota
	AreaFilter
	WayFilter
	RelationFilter
)

type OpType uint

const (
	OpEquals OpType = iota
)

type FilterTag struct {
	Name   string
	Lookup string
	Op     OpType
}

type AST struct {
	Types []FilterType
	Tags  []FilterTag
}

var (
	ErrUndefinedFilter    = errors.New("undefined filter type")
	ErrUnbalancedBrackets = errors.New("unbalanced brackets")
)

func Parse(query string) (*AST, error) {
	setTypes := map[FilterType]struct{}{}

	startTags := 0

	for index, char := range query {
		switch char {
		case 'n':
			setTypes[NodeFilter] = struct{}{}
		case 'a':
			setTypes[AreaFilter] = struct{}{}
		case 'w':
			setTypes[WayFilter] = struct{}{}
		case 'r':
			setTypes[RelationFilter] = struct{}{}
		case '*':
			setTypes[AreaFilter] = struct{}{}
			setTypes[NodeFilter] = struct{}{}
			setTypes[RelationFilter] = struct{}{}
			setTypes[WayFilter] = struct{}{}
		case '[':
			startTags = index

			goto outOfLoop
		default:
			return nil, fmt.Errorf("an undefined type was specified %c: %w", char, ErrUndefinedFilter)
		}
	}

outOfLoop:

	foundTypes := []FilterType{}

	for key := range setTypes {
		foundTypes = append(foundTypes, key)
	}

	sort.Slice(foundTypes, func(i, j int) bool {
		return foundTypes[i] < foundTypes[j]
	})

	tags := []FilterTag{}

	if 0 < startTags {
		brackets := 0

		for index := startTags; index < len(query); index++ {
			switch query[index] {
			case '[':
				brackets++

				tag := FilterTag{}

				var start int

				index++

				for start = index; query[index] != '='; index++ {
				}

				tag.Name = query[start:index]

				index++

				for start = index; query[index] != ']'; index++ {
				}

				tag.Lookup = query[start:index]

				tags = append(tags, tag)
				index--

			case ']':
				brackets--
			}
		}

		if brackets != 0 {
			return nil, fmt.Errorf("could not parse tags: %w", ErrUnbalancedBrackets)
		}
	}

	return &AST{
		Types: foundTypes,
		Tags:  tags,
	}, nil
}
