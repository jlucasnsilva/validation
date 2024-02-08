package validation

import "strings"

type (
	Slice []Validation
)

func (err Slice) insert(valuePath []pathElem, e *Error) Validation {
	if len(valuePath) < 1 {
		return err
	}

	head := valuePath[0]
	tail := valuePath[1:]

	if len(tail) < 1 {
		if head.Index >= len(err) {
			// como os valores aparecem sempre em ordem, se o índice
			// for maior ou igual ao comprimento, significa que ele
			// deve ser adicionado ao final do array (appende).
			return append(err, e)
		}
		err[head.Index] = e
		return err
	}

	if head.Index < len(err) {
		err[head.Index] = err[head.Index].insert(tail, e)
		return err
	}

	var next Validation
	if t := tail[0]; t.Type == elemTypeField {
		next = make(Map)
	} else if t.Type == elemTypeSlice {
		next = make(Slice, 0, 5)
	}
	return append(err, next.insert(tail, e))
}

func (err Slice) get(keys ...any) Validation {
	if len(keys) < 1 {
		return err
	}

	i, ok := keys[0].(int)
	if !ok {
		panic("ErrorSlice should be int")
	}

	if len(keys) < 2 {
		return err[i]
	}
	return err[i].get(keys[1:]...)
}

func (err Slice) Error() string {
	parts := make([]string, len(err))
	for i, v := range err {
		parts[i] = v.Error()
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
