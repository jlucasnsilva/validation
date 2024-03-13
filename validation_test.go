package verno

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type (
	// TODO remover meta: json
	outterType struct {
		Name  string    `validate:"required,min=4" json:"name"`
		Inner innerType `validate:"required"       json:"inner"`
		Slice []string  `validate:"min=2"          json:"slice"`
	}

	// TODO remover meta: json
	innerType struct {
		Inner string     `validate:"required"            json:"inner"`
		Count int        `validate:"min=2"               json:"count"`
		Deep  []deepType `validate:"required,min=2,dive" json:"deep"`
	}

	// TODO remover meta: json
	deepType struct {
		Maximum int                   `validate:"max=-1"        json:"maximum"`
		Deeper  map[string]deeperType `validate:"required,dive" json:"deeper"`
	}

	// TODO remover meta: json
	deeperType struct {
		A string `validate:"required,min=4" json:"a"`
		B string `validate:"required,min=4" json:"b"`
	}
)

func TestPathSection(t *testing.T) {
	key := ""
	index := 10
	name := "hello"
	resName, resIndex, resKey := pathSection(fmt.Sprintf("%v[%v]", name, index))
	assert.Equal(t, name, resName)
	assert.Equal(t, index, resIndex)
	assert.Equal(t, key, resKey)

	key = ""
	index = 111
	name = "FIELD"
	resName, resIndex, resKey = pathSection(fmt.Sprintf("%v[%v]", name, index))
	assert.Equal(t, name, resName)
	assert.Equal(t, index, resIndex)
	assert.Equal(t, key, resKey)

	key = ""
	index = 6
	name = "bazooka"
	resName, resIndex, resKey = pathSection(fmt.Sprintf("%v[%v]", name, index))
	assert.Equal(t, name, resName)
	assert.Equal(t, index, resIndex)
	assert.Equal(t, key, resKey)

	key = "10aaaa"
	index = 0
	name = ""
	resName, resIndex, resKey = pathSection(fmt.Sprintf("%v[%v]", name, key))
	assert.Equal(t, name, resName)
	assert.Equal(t, index, resIndex)
	assert.Equal(t, key, resKey)
}

func TestSplitErrNamespace(t *testing.T) {
	ns := "Struct.field"
	expected := []fieldPathElem{{
		Type:  errField,
		Field: "field",
	}}
	assert.Equal(t, expected, splitNamespace(ns))

	ns = "Struct.field.subfield"
	expected = []fieldPathElem{
		{
			Type:  errField,
			Field: "field",
		},
		{
			Type:  errField,
			Field: "subfield",
		},
	}
	assert.Equal(t, expected, splitNamespace(ns))

	ns = "Struct.field[12]"
	expected = []fieldPathElem{
		{
			Type:  errField,
			Field: "field",
		},
		{
			Type:  errSlice,
			Index: 12,
		},
	}
	assert.Equal(t, expected, splitNamespace(ns))

	ns = "Struct.field[12].subfield"
	expected = []fieldPathElem{
		{
			Type:  errField,
			Field: "field",
		},
		{
			Type:  errSlice,
			Index: 12,
		},
		{
			Type:  errField,
			Field: "subfield",
		},
	}
	assert.Equal(t, expected, splitNamespace(ns))

	ns = "Struct.field.subfield[80].deepField"
	expected = []fieldPathElem{
		{
			Type:  errField,
			Field: "field",
		},
		{
			Type:  errField,
			Field: "subfield",
		},
		{
			Type:  errSlice,
			Index: 80,
		},
		{
			Type:  errField,
			Field: "deepField",
		},
	}
	assert.Equal(t, expected, splitNamespace(ns))
}

func TestInsert(t *testing.T) {
	m := make(Map)

	tag := "required"
	path := splitNamespace("useless.hello.world")
	expected := Map{"hello": Map{"world": &Error{Type: tag}}}
	assert.Equal(t, expected, m.insert(path, &Error{Type: tag}))

	tag = "min"
	path = splitNamespace("useless.foo.bar[0].baz")
	expected = Map{
		"hello": Map{"world": &Error{Type: "required"}},
		"foo": Map{
			"bar": Slice{
				Map{"baz": &Error{Type: tag}},
			},
		},
	}
	assert.Equal(t, expected, m.insert(path, &Error{Type: tag}))

	tag = "max"
	path = splitNamespace("useless.foo.bar[0].jazz")
	expected = Map{
		"hello": Map{"world": &Error{Type: "required"}},
		"foo": Map{
			"bar": Slice{
				Map{"baz": &Error{Type: "min"}, "jazz": &Error{Type: tag}},
			},
		},
	}
	assert.Equal(t, expected, m.insert(path, &Error{Type: tag}))

	tag = "min"
	path = splitNamespace("useless.foo.bar[1].baz")
	expected = Map{
		"hello": Map{"world": &Error{Type: "required"}},
		"foo": Map{
			"bar": Slice{
				Map{"baz": &Error{Type: "min"}, "jazz": &Error{Type: "max"}},
				Map{"baz": &Error{Type: tag}},
			},
		},
	}
	assert.Equal(t, expected, m.insert(path, &Error{Type: tag}))

	tag = "max"
	path = splitNamespace("useless.foo.bar[1].jazz")
	expected = Map{
		"hello": Map{"world": &Error{Type: "required"}},
		"foo": Map{
			"bar": Slice{
				Map{"baz": &Error{Type: "min"}, "jazz": &Error{Type: "max"}},
				Map{"baz": &Error{Type: "min"}, "jazz": &Error{Type: tag}},
			},
		},
	}
	assert.Equal(t, expected, m.insert(path, &Error{Type: tag}))

	tag = "time"
	path = splitNamespace("useless.foo.zhoda")
	expected = Map{
		"hello": Map{"world": &Error{Type: "required"}},
		"foo": Map{
			"bar": Slice{
				Map{"baz": &Error{Type: "min"}, "jazz": &Error{Type: "max"}},
				Map{"baz": &Error{Type: "min"}, "jazz": &Error{Type: "max"}},
			},
			"zhoda": &Error{Type: tag},
		},
	}
	assert.Equal(t, expected, m.insert(path, &Error{Type: tag}))
}

func TestValidate(t *testing.T) {
	v := validator.New()
	result := v.Struct(&outterType{
		Slice: []string{""},
		Inner: innerType{
			Deep: []deepType{
				{
					Deeper: map[string]deeperType{
						"hadouken": {},
						"flight":   {},
					},
				},
				{
					Deeper: map[string]deeperType{
						"hadouken": {},
						"flight":   {},
					},
				},
			},
		},
	})

	expected := Map{
		"Name": &Error{
			Type:  "required",
			Field: "Name",
			Path:  "outterType.Name",
			Param: "",
			Value: "",
		},
		"Slice": &Error{
			Type:  "min",
			Field: "Slice",
			Param: "2",
			Value: []string{""},
			Path:  "outterType.Slice",
		},
		"Inner": Map{
			"Inner": &Error{
				Type:  "required",
				Field: "Inner",
				Path:  "outterType.Inner.Inner",
				Value: "",
			},
			"Count": &Error{
				Type:  "min",
				Field: "Count",
				Param: "2",
				Value: 0,
				Path:  "outterType.Inner.Count",
			},
			"Deep": Slice{
				Map{
					"Deeper": Map{
						"hadouken": Map{
							"A": &Error{
								Type:  "required",
								Field: "A",
								Param: "",
								Value: "",
								Path:  "outterType.Inner.Deep[0].Deeper[hadouken].A",
							},
							"B": &Error{
								Type:  "required",
								Field: "B",
								Param: "",
								Value: "",
								Path:  "outterType.Inner.Deep[0].Deeper[hadouken].B",
							},
						},
						"flight": Map{
							"A": &Error{
								Type:  "required",
								Field: "A",
								Param: "",
								Value: "",
								Path:  "outterType.Inner.Deep[0].Deeper[flight].A",
							},
							"B": &Error{
								Type:  "required",
								Field: "B",
								Param: "",
								Value: "",
								Path:  "outterType.Inner.Deep[0].Deeper[flight].B",
							},
						},
					},
					"Maximum": &Error{
						Type:  "max",
						Field: "Maximum",
						Param: "-1",
						Value: 0,
						Path:  "outterType.Inner.Deep[0].Maximum",
					},
				},
				Map{
					"Deeper": Map{
						"hadouken": Map{
							"A": &Error{
								Type:  "required",
								Field: "A",
								Param: "",
								Value: "",
								Path:  "outterType.Inner.Deep[1].Deeper[hadouken].A",
							},
							"B": &Error{
								Type:  "required",
								Field: "B",
								Param: "",
								Value: "",
								Path:  "outterType.Inner.Deep[1].Deeper[hadouken].B",
							},
						},
						"flight": Map{
							"A": &Error{
								Type:  "required",
								Field: "A",
								Param: "",
								Value: "",
								Path:  "outterType.Inner.Deep[1].Deeper[flight].A",
							},
							"B": &Error{
								Type:  "required",
								Field: "B",
								Param: "",
								Value: "",
								Path:  "outterType.Inner.Deep[1].Deeper[flight].B",
							},
						},
					},
					"Maximum": &Error{
						Type:  "max",
						Field: "Maximum",
						Param: "-1",
						Value: 0,
						Path:  "outterType.Inner.Deep[1].Maximum",
					},
				},
			},
		},
	}

	assert.Equal(t, expected, Convert(result))
}

func TestTranslate(t *testing.T) {
	tr := func(e Error) string {
		return fmt.Sprintf("%v|%v[%v, %v]", e.Type, e.Field, e.Param, e.Value)
	}
	err := Map{
		"hello": Slice{
			&Error{
				Type:  "min",
				Field: "field",
				Path:  "a.path",
				Param: "1",
				Value: "2",
			},
			&Error{
				Type:  "min",
				Field: "field",
				Path:  "b.path",
				Param: "1",
				Value: "3",
			},
		},
	}
	expected := map[string]any{
		"hello": []any{
			"min|field[1, 2]",
			"min|field[1, 3]",
		},
	}
	result := err.Translate(tr)

	assert.Equal(t, expected, result)
}
