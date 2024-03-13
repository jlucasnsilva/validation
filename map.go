package validation

import (
	"fmt"
	"strings"
)

type (
	Map map[string]Validation
)

func New(field, err string) Validation {
	return Map{field: &Error{Type: err}}
}

func (err Map) insert(valuePath []fieldPathElem, verr *Error) Validation {
	if len(valuePath) < 1 {
		return err
	}

	head := valuePath[0]
	tail := valuePath[1:]

	if len(tail) < 1 {
		err[head.Field] = verr
		return err
	}

	next, ok := err[head.Field]
	if !ok {
		if t := tail[0]; t.Type == errField {
			next = make(Map)
		} else if t.Type == errSlice {
			next = make(Slice, 0, 5)
		}
	}
	err[head.Field] = next.insert(tail, verr)
	return err
}

func (err Map) Error() string {
	parts := make([]string, 0, len(err))
	for k, v := range err {
		parts = append(parts, fmt.Sprintf(`"%v": %v`, k, v.Error()))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

func (err Map) Translate(tr TranslatorFunc) any {
	m := make(map[string]any)
	for k, v := range err {
		m[k] = v.Translate(tr)
	}
	return m
}
