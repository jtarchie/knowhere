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

type AST struct {
	Types []FilterType
}

var ErrUndefinedFilter = errors.New("undefined filter type")

func Parse(query string) (*AST, error) {
	setTypes := map[FilterType]struct{}{}

	for _, char := range query {
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
			setTypes[NodeFilter] = struct{}{}
			setTypes[AreaFilter] = struct{}{}
			setTypes[WayFilter] = struct{}{}
			setTypes[RelationFilter] = struct{}{}
		default:
			return nil, fmt.Errorf("an undefined type was specified %b: %w", char, ErrUndefinedFilter)
		}
	}

	foundTypes := []FilterType{}

	for key := range setTypes {
		foundTypes = append(foundTypes, key)
	}

	sort.Slice(foundTypes, func(i, j int) bool {
		return foundTypes[i] < foundTypes[j]
	})

	return &AST{
		Types: foundTypes,
	}, nil
}
