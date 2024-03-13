package validation

type (
	Error struct {
		Type  string
		Field string
		Path  string
		Param string
		Value any
	}
)

func (err *Error) insert(valuePath []fieldPathElem, verr *Error) Validation {
	return err
}

func (err *Error) MarshalJSON() ([]byte, error) {
	return []byte(`"` + err.Type + `"`), nil
}

func (err *Error) Error() string {
	return err.Type
}

func (err *Error) Translate(tr TranslatorFunc) any {
	return tr(*err)
}
