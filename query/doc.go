package query

//go:generate ragel -e -G2 -Z parse.rl

import (
	"errors"
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
	OpNotEquals
	OpExists
	OpNotExists
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
