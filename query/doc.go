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
)

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
