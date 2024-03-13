package verno

import "strings"

type (
	Slice []Validation
)

func NewSlice(length, capacity int) Slice {
	return make(Slice, length, capacity)
}

func (err Slice) insert(valuePath []fieldPathElem, verr *Error) Validation {
	if len(valuePath) < 1 {
		return err
	}

	head := valuePath[0]
	tail := valuePath[1:]

	if len(tail) < 1 {
		val := verr
		if head.Index >= len(err) {
			// como os valores aparecem sempre em ordem, se o índice
			// for maior ou igual ao comprimento, significa que ele
			// deve ser adicionado ao final do array (appende).
			return append(err, val)
		}
		err[head.Index] = val
		return err
	}

	if head.Index < len(err) {
		err[head.Index] = err[head.Index].insert(tail, verr)
		return err
	}

	var next Validation
	if t := tail[0]; t.Type == errField {
		next = make(Map)
	} else if t.Type == errSlice {
		next = make(Slice, 0, 5)
	}
	return append(err, next.insert(tail, verr))
}

func (err Slice) Error() string {
	parts := make([]string, len(err))
	for i, v := range err {
		parts[i] = v.Error()
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func (err Slice) Translate(tr TranslatorFunc) any {
	s := make([]any, 0, len(err))
	for _, e := range err {
		s = append(s, e.Translate(tr))
	}
	return s
}
