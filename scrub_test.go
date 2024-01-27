package scrub

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaggedField(t *testing.T) {
	primitiveTypeIsUnchanged := func(v any) bool {
		typ := reflect.TypeOf(v)
		actual := reflect.New(typ).Interface()
		expected := reflect.New(typ).Interface()
		TaggedFields(actual)
		return reflect.DeepEqual(actual, expected)
	}

	t.Run("with a primitive type, leaves the value unchanged", func(t *testing.T) {
		for _, v := range []any{1, 1.0, "hello"} {
			vPtr := &v
			assert.True(t, primitiveTypeIsUnchanged(v))
			assert.True(t, primitiveTypeIsUnchanged(vPtr))
			assert.True(t, primitiveTypeIsUnchanged(&vPtr))
		}
	})

	t.Run("nil is handled safely", func(t *testing.T) {
		var nilPtr *int
		TaggedFields(nilPtr)
		assert.Nil(t, nilPtr)

		TaggedFields(nil)
	})

	t.Run("with a mix of tagged and untagged fields, only scrubs the tagged ones", func(t *testing.T) {
		type person struct {
			Name     string `scrub:"true"`
			Age      int
			Interest string
			Height   float64 `scrub:"true"`
		}
		actual := person{Name: "Testy Tester", Age: 26, Interest: "Bitwise arithmetic", Height: 5.8}
		expected := person{Name: "", Age: 26, Interest: "Bitwise arithmetic", Height: 0}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a pointer to a struct, scrubs the tagged fields", func(t *testing.T) {
		type place struct {
			Name      string `scrub:"true"`
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name  string `scrub:"true"`
			Age   int
			Place *place
		}
		actual := person{Name: "Testy Tester", Age: 26, Place: &place{Name: "Testy Tester", Latitude: 1.0, Longitude: 2.0}}
		expected := person{
			Name: "",
			Age:  26,
			Place: &place{
				Name:      "",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a non-nil pointer to a non-struct, scrubs the tagged fields", func(t *testing.T) {
		type person struct {
			Name  string `scrub:"true"`
			Age   int
			Place *string
		}
		place := "earth"
		actual := person{Name: "Testy Tester", Age: 26, Place: &place}
		expected := person{
			Name:  "",
			Age:   26,
			Place: &place,
		}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a tagged struct, scrubs the tagged fields", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name  string `scrub:"true"`
			Age   int
			Place place `scrub:"true"`
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Place: place{
				Name:      "Testy Tester",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		expected := person{
			Name: "",
			Age:  26,
			Place: place{
				Name:      "",
				Latitude:  0.0,
				Longitude: 0.0,
			},
		}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing an untagged struct, scrubs the tagged fields", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name  string `scrub:"true"`
			Age   int
			Place place
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Place: place{
				Name:      "Testy Tester",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		expected := person{
			Name: "",
			Age:  26,
			Place: place{
				Name:      "Testy Tester",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a tagged struct pointer, scrubs the tagged fields", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name  string `scrub:"true"`
			Age   int
			Place *place `scrub:"true"`
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Place: &place{
				Name:      "Testy Tester",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		expected := person{
			Name:  "",
			Age:   26,
			Place: nil,
		}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a tagged struct pointer that's nil, leaves the pointer nil", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name  string `scrub:"true"`
			Age   int
			Place *place `scrub:"true"`
		}
		actual := person{
			Name:  "Testy Tester",
			Age:   26,
			Place: nil,
		}
		expected := person{
			Name:  "",
			Age:   26,
			Place: nil,
		}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})
}

func TestNamedFields(t *testing.T) {
	t.Run("with a non-struct type, leaves the value unchanged", func(t *testing.T) {
		nonStructValues := []any{
			1,
			reflect.New(reflect.TypeOf(1)).Interface().(*int),
			reflect.New(reflect.TypeOf(uint(1))).Interface().(*uint),
			reflect.New(reflect.TypeOf(uintptr(1))).Interface().(*uintptr),
			*reflect.New(reflect.TypeOf(uintptr(1))).Interface().(*uintptr),
			1.0,
			reflect.New(reflect.TypeOf(1.0)).Interface().(*float64),
			"hello",
			reflect.New(reflect.TypeOf("hello")).Interface().(*string),
		}

		for _, v := range nonStructValues {
			copy := v
			NamedFields(&copy, []string{}...)
			assert.Equal(t, v, copy)
		}
	})

	t.Run("nil is handled safely", func(t *testing.T) {
		var nilPtr *int
		NamedFields(nilPtr)
		assert.Nil(t, nilPtr)

		NamedFields(nil)
	})

	t.Run("with a mix of named and unnamed fields, only scrubs the named ones", func(t *testing.T) {
		type person struct {
			Name     string
			Age      int
			Interest string
			Height   float64
		}
		actual := person{Name: "Testy Tester", Age: 26, Interest: "Bitwise arithmetic", Height: 5.8}
		expected := person{Name: "", Age: 26, Interest: "Bitwise arithmetic", Height: 0}
		NamedFields(&actual, "Name", "Height")
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a pointer to a struct, scrubs the named fields", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name  string
			Age   int
			Place *place
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Place: &place{
				Name:      "Testy Tester",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		expected := person{
			Name: "",
			Age:  26,
			Place: &place{
				Name:      "",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		NamedFields(&actual, "Name")
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a non-nil pointer to a non-struct, scrubs the tagged fields", func(t *testing.T) {
		type person struct {
			Name  string
			Age   int
			Place *string
		}
		place := "earth"
		actual := person{Name: "Testy Tester", Age: 26, Place: &place}
		expected := person{
			Name:  "",
			Age:   26,
			Place: &place,
		}
		NamedFields(&actual, "Name")
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a named struct, scrubs the named fields", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name  string
			Age   int
			Place place
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Place: place{
				Name:      "Testy Tester",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		expected := person{
			Name: "",
			Age:  26,
			Place: place{
				Name:      "",
				Latitude:  0.0,
				Longitude: 0.0,
			},
		}
		NamedFields(&actual, "Name", "Place")
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing an unnamed struct, scrubs the named fields", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name  string
			Age   int
			Place place
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Place: place{
				Name:      "Testy Tester",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		expected := person{
			Name: "",
			Age:  26,
			Place: place{
				Name:      "",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		NamedFields(&actual, "Name")
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a named struct pointer that's nil, leaves the pointer nil", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name  string
			Age   int
			Place *place
		}
		actual := person{
			Name:  "Testy Tester",
			Age:   26,
			Place: nil,
		}
		expected := person{
			Name:  "",
			Age:   26,
			Place: nil,
		}
		NamedFields(&actual, "Name", "Place")
		assert.Equal(t, expected, actual)
	})
}
