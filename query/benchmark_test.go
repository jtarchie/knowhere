package query_test

import (
	"regexp"
	"testing"

	"github.com/jtarchie/knowhere/query"
)

//nolint: gochecknoglobals
var (
	search = `wn[name="Starbucks"][amenity=cafe]`
	ast    *query.AST
)

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ast, _ = query.Parse(search)
	}
}

var tagRegex = regexp.MustCompile(`\[(\w+)="([^"]+)"\]`)

func BenchmarkParseRegex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = tagRegex.FindAllStringSubmatch(search, -1)
	}
}

func BenchmarkToSQL(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = query.ToSQL(search)
	}
}
