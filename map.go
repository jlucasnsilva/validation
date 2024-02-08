package validation

import (
	"fmt"
	"strings"
)

type (
	Map map[string]Validation
)

func (err Map) insert(valuePath []pathElem, e *Error) Validation {
	if len(valuePath) < 1 {
		return err
	}

	head := valuePath[0]
	tail := valuePath[1:]

	if len(tail) < 1 {
		err[head.Field] = e
		return err
	}

	next, ok := err[head.Field]
	if !ok {
		if t := tail[0]; t.Type == elemTypeField {
			next = make(Map)
		} else if t.Type == elemTypeSlice {
			next = make(Slice, 0, 5)
		}
	}
	err[head.Field] = next.insert(tail, e)
	return err
}

func (err Map) get(keys ...any) Validation {
	if len(keys) < 1 {
		return err
	}

	k, ok := keys[0].(string)
	if !ok {
		s := fmt.Sprintf(
			"##)) ErrorMap: key '%v' (of type %T) should be of type string",
			keys[0],
			keys[0],
		)
		panic(s)
	}

	if len(keys) < 2 {
		return err[k]
	}
	return err[k].get(keys[1:]...)
}

func (err Map) Error() string {
	parts := make([]string, 0, len(err))
	for k, v := range err {
		parts = append(parts, fmt.Sprintf(`"%v": %v`, k, v.Error()))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
