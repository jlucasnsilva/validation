package validation

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	fieldPathElem struct {
		Type  int
		Index int
		Field string
	}

	TranslatorFunc func(Error) string

	Validation interface {
		error
		Translate(TranslatorFunc) any

		insert([]fieldPathElem, *Error) Validation
	}
)

var (
	indexRegExp = regexp.MustCompile(`\[([0-9]*)\]$`)
	keyRegExp   = regexp.MustCompile(`\[([a-zA-Z0-9]*)\]$`)
)

const (
	errField = iota
	errSlice
)

func Convert(err error) error {
	if err == nil {
		return nil
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	m := make(Map)
	for _, e := range errs {
		pathElems := splitNamespace(e.Namespace())
		verr := &Error{
			Type:  e.Tag(),
			Field: e.Field(),
			Param: e.Param(),
			Value: e.Value(),
			Path:  e.Namespace(),
		}
		m.insert(pathElems, verr)
	}

	return m
}

func splitNamespace(errNamespace string) []fieldPathElem {
	parts := strings.Split(errNamespace, ".")

	switch len(parts) {
	case 0:
		return nil
	case 1:
		return []fieldPathElem{{Type: errField, Field: errNamespace}}
	case 2:
		return createFieldErrorPathElem(parts[1], true)
	default:
		result := make([]fieldPathElem, 0, len(parts)-1)
		tail := parts[1:]
		lastIdx := len(tail) - 1
		for i, part := range tail {
			pathElems := createFieldErrorPathElem(part, i == lastIdx)
			result = append(result, pathElems...)
		}
		return result
	}
}

func createFieldErrorPathElem(part string, isLast bool) []fieldPathElem {
	if !strings.HasSuffix(part, "]") {
		return []fieldPathElem{{
			Type:  errField,
			Field: part,
		}}
	}

	name, index, key := pathSection(part)
	if key == "" {
		return []fieldPathElem{
			{
				Field: name,
				Type:  errField,
			},
			{
				Index: index,
				Type:  errSlice,
			},
		}
	}

	return []fieldPathElem{
		{
			Field: name,
			Index: index,
			Type:  errField,
		},
		{
			Field: key,
			Index: index,
			Type:  errField,
		},
	}
}

// Retorna as propriedades da seção de um caminho de validação (por
// exemplo: campoA[chaveA].campoB).
func pathSection(s string) (name string, index int, key string) {
	name = strings.Split(s, "[")[0]
	match := indexRegExp.FindStringSubmatch(s)
	if len(match) < 2 {
		// Essa condição é alcançada quando o conteúdo entre colchetes
		// for um identificador, e não um número (fazendo com que indexRegExp
		// não gere matches).
		match = keyRegExp.FindStringSubmatch(s)
		return name, 0, match[1]
	}
	index, _ = strconv.Atoi(match[1])
	return name, index, ""
}
