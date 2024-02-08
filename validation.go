package validation

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	Validation interface {
		error
		insert([]pathElem, *Error) Validation
		get(keys ...any) Validation // get é um método para auxiliar nos testes
	}

	pathElem struct {
		Type  int
		Index int
		Field string
	}
)

const (
	elemTypeField = iota
	elemTypeSlice
)

var (
	indexRegExp = regexp.MustCompile(`\[([0-9]*)\]$`)
	keyRegExp   = regexp.MustCompile(`\[([a-zA-Z0-9]*)\]$`)
)

func New(verrs validator.ValidationErrors) Validation {
	if verrs == nil {
		return nil
	}

	m := make(Map)
	for _, e := range verrs {
		if p := splitNamespace(e.Namespace()); p != nil {
			m.insert(p, newError(e))
		}
	}
	return m
}

func splitNamespace(namespace string) []pathElem {
	parts := strings.Split(namespace, ".")

	switch len(parts) {
	case 1:
		return []pathElem{{Type: elemTypeField, Field: namespace}}
	case 2:
		return createPathElems(parts[1], true)
	default:
		result := make([]pathElem, 0, len(parts)-1)
		tail := parts[1:]
		lastIdx := len(tail) - 1
		for i, part := range tail {
			elems := createPathElems(part, i == lastIdx)
			result = append(result, elems...)
		}
		return result
	}
}

func createPathElems(part string, isLast bool) []pathElem {
	if !strings.HasSuffix(part, "]") {
		return []pathElem{{
			Type:  elemTypeField,
			Field: part,
		}}
	}

	name, index, key := pathSection(part)
	if key == "" {
		return []pathElem{
			{
				Field: name,
				Type:  elemTypeField,
			},
			{
				Index: index,
				Type:  elemTypeSlice,
			},
		}
	}

	return []pathElem{
		{
			Field: name,
			Index: index,
			Type:  elemTypeField,
		},
		{
			Field: key,
			Index: index,
			Type:  elemTypeField,
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
