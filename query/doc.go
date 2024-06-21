package query

//go:generate ragel -e -G2 -Z parse.rl

import (
	"errors"
)

type FilterType uint

const (
	NodeFilter     FilterType = 1
	WayFilter      FilterType = 2
	RelationFilter FilterType = 3
)

func (f FilterType) String() string {
	switch f {
	case NodeFilter:
		return "1"
	case WayFilter:
		return "2"
	case RelationFilter:
		return "3"
	}

	return ""
}

type OpType uint

const (
	OpEquals OpType = iota
	OpNotEquals
	OpExists
	OpNotExists
	OpGreaterThan
	OpGreaterThanEquals
	OpLessThan
	OpLessThanEquals
	OpContains
	OpNotContains
)

func (o OpType) String() string {
	switch o {
	case OpGreaterThan:
		return ">"
	case OpGreaterThanEquals:
		return ">="
	case OpLessThan:
		return "<"
	case OpLessThanEquals:
		return "<="
	}

	return ""
}

type FilterTag struct {
	Name    string
	Lookups []string
	Op      OpType
}

type FilterDirective []string

type AST struct {
	Directives map[string]FilterDirective
	Tags       []FilterTag
	Types      []FilterType
}

var (
	ErrUndefinedFilter    = errors.New("undefined filter type")
	ErrUnbalancedBrackets = errors.New("unbalanced brackets")
	ErrUnparsableQuery    = errors.New("unparsable query")
)
