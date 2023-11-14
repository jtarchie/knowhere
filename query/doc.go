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
	OpEquals    OpType = iota
	OpNotEquals OpType = iota
	OpExists
	OpNotExists
)

type FilterTag struct {
	Name    string
	Lookups []string
	Op      OpType
}

type AST struct {
	Types []FilterType
	Tags  []FilterTag
}

var (
	ErrUndefinedFilter    = errors.New("undefined filter type")
	ErrUnbalancedBrackets = errors.New("unbalanced brackets")
	ErrUnparsableQuery    = errors.New("unparsable query")
)
