package query

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/samber/lo"
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
	OpExists
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
	foundTypes := []FilterType{}
	tags := []FilterTag{}
	scanner := bufio.NewReader(strings.NewReader(query))

	for {
		char, err := scanner.ReadByte()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("could not read character while scanning for types: %w", err)
		}

		switch char {
		case 'n':
			foundTypes = append(foundTypes, NodeFilter)
		case 'a':
			foundTypes = append(foundTypes, AreaFilter)
		case 'w':
			foundTypes = append(foundTypes, WayFilter)
		case 'r':
			foundTypes = append(foundTypes, RelationFilter)
		case '*':
			foundTypes = append(foundTypes, NodeFilter, AreaFilter, WayFilter, RelationFilter)
		default:
			return nil, fmt.Errorf("an undefined type was specified %c: %w", char, ErrUndefinedFilter)
		}

		peek, _ := scanner.Peek(1)
		if string(peek) == "[" {
			break
		}
	}

	brackets := 0

	peek, _ := scanner.Peek(1)
	if string(peek) == "[" {
		for {
			char, err := scanner.ReadByte()

			//nolint: errorlint
			if err == io.EOF {
				break
			}

			if err != nil {
				return nil, fmt.Errorf("could not read character while scanning for tags: %w", err)
			}

			switch char {
			case '[':
				brackets++

				tag := FilterTag{}
				tag.Name, err = readWord(scanner)

				if err != nil {
					return nil, fmt.Errorf("could not read tag name: %w", err)
				}

				op, _ := scanner.ReadByte()
				switch op {
				case '=':
					tag.Op = OpEquals
					tag.Lookup, err = readWord(scanner)

					if err != nil {
						return nil, fmt.Errorf("could not read tag assignment: %w", err)
					}
				case ']':
					tag.Op = OpExists
					_ = scanner.UnreadByte()
				}

				tags = append(tags, tag)

			case ']':
				brackets--
			}
		}

		if brackets != 0 {
			return nil, fmt.Errorf("could not parse tags: %w", ErrUnbalancedBrackets)
		}
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

func readWord(scanner *bufio.Reader) (string, error) {
	var builder strings.Builder

	for {
		char, err := scanner.ReadByte()

		if err == io.EOF {
			return builder.String(), nil
		}

		if err != nil {
			return "", fmt.Errorf("could not read byte for word: %w", err)
		}

		if 'a' <= char && char <= 'z' {
			_ = builder.WriteByte(char)
		} else {
			_ = scanner.UnreadByte()

			break
		}
	}

	return builder.String(), nil
}
