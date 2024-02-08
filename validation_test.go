package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	outterType struct {
		Name  string    `validate:"required,min=4" json:"name"`
		Inner innerType `validate:"required"       json:"inner"`
		Slice []string  `validate:"min=2"          json:"slice"`
	}

	innerType struct {
		Inner string     `validate:"required"            json:"inner"`
		Count int        `validate:"min=2"               json:"count"`
		Deep  []deepType `validate:"required,min=2,dive" json:"deep"`
	}

	deepType struct {
		Maximum int                   `validate:"max=-1"        json:"maximum"`
		Deeper  map[string]deeperType `validate:"required,dive" json:"deeper"`
	}

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
	expected := []pathElem{{
		Type:  elemTypeField,
		Field: "field",
	}}
	assert.Equal(t, expected, splitNamespace(ns))

	ns = "Struct.field.subfield"
	expected = []pathElem{
		{
			Type:  elemTypeField,
			Field: "field",
		},
		{
			Type:  elemTypeField,
			Field: "subfield",
		},
	}
	assert.Equal(t, expected, splitNamespace(ns))

	ns = "Struct.field[12]"
	expected = []pathElem{
		{
			Type:  elemTypeField,
			Field: "field",
		},
		{
			Type:  elemTypeSlice,
			Index: 12,
		},
	}
	assert.Equal(t, expected, splitNamespace(ns))

	ns = "Struct.field[12].subfield"
	expected = []pathElem{
		{
			Type:  elemTypeField,
			Field: "field",
		},
		{
			Type:  elemTypeSlice,
			Index: 12,
		},
		{
			Type:  elemTypeField,
			Field: "subfield",
		},
	}
	assert.Equal(t, expected, splitNamespace(ns))

	ns = "Struct.field.subfield[80].deepField"
	expected = []pathElem{
		{
			Type:  elemTypeField,
			Field: "field",
		},
		{
			Type:  elemTypeField,
			Field: "subfield",
		},
		{
			Type:  elemTypeSlice,
			Index: 80,
		},
		{
			Type:  elemTypeField,
			Field: "deepField",
		},
	}
	assert.Equal(t, expected, splitNamespace(ns))
}

func TestInsert(t *testing.T) {
	m := make(Map)

	errRequired := &Error{Type: "required"}
	path := splitNamespace("useless.hello.world")
	expected := Map{"hello": Map{"world": errRequired}}
	assert.Equal(t, expected, m.insert(path, errRequired))

	errMin := &Error{Type: "min"}
	path = splitNamespace("useless.foo.bar[0].baz")
	expected = Map{
		"hello": Map{"world": errRequired},
		"foo": Map{
			"bar": Slice{
				Map{"baz": errMin},
			},
		},
	}
	assert.Equal(t, expected, m.insert(path, errMin))

	errMax := &Error{Type: "max"}
	path = splitNamespace("useless.foo.bar[0].jazz")
	expected = Map{
		"hello": Map{"world": errRequired},
		"foo": Map{
			"bar": Slice{
				Map{"baz": errMin, "jazz": errMax},
			},
		},
	}
	assert.Equal(t, expected, m.insert(path, errMax))

	path = splitNamespace("useless.foo.bar[1].baz")
	expected = Map{
		"hello": Map{"world": errRequired},
		"foo": Map{
			"bar": Slice{
				Map{"baz": errMin, "jazz": errMax},
				Map{"baz": errMin},
			},
		},
	}
	assert.Equal(t, expected, m.insert(path, errMin))

	path = splitNamespace("useless.foo.bar[1].jazz")
	expected = Map{
		"hello": Map{"world": errRequired},
		"foo": Map{
			"bar": Slice{
				Map{"baz": errMin, "jazz": errMax},
				Map{"baz": errMin, "jazz": errMax},
			},
		},
	}
	assert.Equal(t, expected, m.insert(path, errMax))

	errTime := &Error{Type: "time"}
	path = splitNamespace("useless.foo.zhoda")
	expected = Map{
		"hello": Map{"world": errRequired},
		"foo": Map{
			"bar": Slice{
				Map{"baz": errMin, "jazz": errMax},
				Map{"baz": errMin, "jazz": errMax},
			},
			"zhoda": errTime,
		},
	}
	assert.Equal(t, expected, m.insert(path, errTime))
}
