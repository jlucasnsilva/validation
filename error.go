package validation

import "github.com/go-playground/validator/v10"

type (
	Error struct {
		Type  string
		Param string
		Value any
	}
)

func newError(e validator.FieldError) *Error {
	return &Error{
		Type:  e.Tag(),
		Param: e.Param(),
		Value: e.Value(),
	}
}

func (err *Error) insert(valuePath []pathElem, e *Error) Validation {
	return err
}

func (err *Error) get(keys ...any) Validation {
	return err
}

func (err *Error) MarshalJSON() ([]byte, error) {
	return []byte(`"` + err.Type + `"`), nil
}

func (err *Error) Error() string {
	return err.Type
}
