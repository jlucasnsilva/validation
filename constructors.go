package validation

func FieldValidationError(field, err string) Map {
	return Map{field: &Error{Type: err}}
}
