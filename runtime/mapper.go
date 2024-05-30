package runtime

import (
	"reflect"
	"strings"

	"github.com/dop251/goja"
)

type tagFieldNameMapper struct {
	cache map[string]string
}

func (tfm *tagFieldNameMapper) FieldName(t reflect.Type, f reflect.StructField) string {
	name := t.PkgPath() + "." + t.Name() + "." + f.Name

	if tag, ok := tfm.cache[name]; ok {
		return tag
	}

	field := f.Tag.Get("js")
	if idx := strings.IndexByte(field, ','); idx != -1 {
		tag := field[:idx]
		tfm.cache[name] = tag
		return tag
	}

	fieldName := uncapitalize(f.Name)
	tfm.cache[name] = fieldName

	return fieldName
}

func uncapitalize(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

func (tfm *tagFieldNameMapper) MethodName(_ reflect.Type, m reflect.Method) string {
	return uncapitalize(m.Name)
}

var _ goja.FieldNameMapper = &tagFieldNameMapper{}
