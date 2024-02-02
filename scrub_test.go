package scrub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaggedField(t *testing.T) {
	primitiveTypeIsUnchanged := func(v any) bool {
		copy := v
		actual := v
		expected := copy
		TaggedFields(actual)
		return expected == actual
	}

	t.Run("with a primitive type, leaves the value unchanged", func(t *testing.T) {
		primitives := []any{
			int(1), int8(1), int16(1), int32(1), int64(1),
			uintptr(1), uint(1), uint8(1), uint16(1), uint32(1), uint64(1),
			float32(1.0), float64(1.0),
			"hello",
		}
		for _, v := range primitives {
			vPtr := &v
			assert.True(t, primitiveTypeIsUnchanged(v))
			assert.True(t, primitiveTypeIsUnchanged(vPtr))
			assert.True(t, primitiveTypeIsUnchanged(&vPtr))
		}

		intPtr := new(int)
		intDoublePtr := &intPtr
		*intPtr = 1
		TaggedFields(intPtr)
		assert.Equal(t, 1, *intPtr)
		assert.Equal(t, *intDoublePtr, intPtr)
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

	t.Run("with a struct containing an unexported pointer to a struct, leaves the unexported field alone", func(t *testing.T) {
		type place struct {
			Name      string `scrub:"true"`
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name  string `scrub:"true"`
			Age   int
			place *place
		}
		actual := person{Name: "Testy Tester", Age: 26, place: &place{Name: "Testy Tester", Latitude: 1.0, Longitude: 2.0}}
		expected := person{
			Name: "",
			Age:  26,
			place: &place{
				Name:      "Testy Tester",
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

	t.Run("with a struct containing a tagged slice, scrubs the slice", func(t *testing.T) {
		type person struct {
			Name   string `scrub:"true"`
			Age    int
			Places []string `scrub:"true"`
		}
		actual := person{
			Name:   "Testy Tester",
			Age:    26,
			Places: []string{"earth", "mars"},
		}
		expected := person{
			Name:   "",
			Age:    26,
			Places: nil,
		}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a slice of structs, walks the slice and scrubs any tagged fields on those structs", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64 `scrub:"true"`
			Longitude float64 `scrub:"true"`
		}
		type person struct {
			Name   string
			Age    int
			Places []place
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []place{
				{Name: "Place 1", Latitude: 1.0, Longitude: 2.0},
				{Name: "Place 2", Latitude: 3.14, Longitude: 1.5926},
			},
		}
		expected := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []place{
				{Name: "Place 1", Latitude: 0.0, Longitude: 0.0},
				{Name: "Place 2", Latitude: 0.0, Longitude: 0.0},
			},
		}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a slice of pointers to structs, walks the slice and scrubs any tagged fields on those structs", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64 `scrub:"true"`
			Longitude float64 `scrub:"true"`
		}
		type person struct {
			Name   string
			Age    int
			Places []*place
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []*place{
				{Name: "Place 1", Latitude: 1.0, Longitude: 2.0},
				{Name: "Place 2", Latitude: 3.14, Longitude: 1.5926},
			},
		}
		expected := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []*place{
				{Name: "Place 1", Latitude: 0.0, Longitude: 0.0},
				{Name: "Place 2", Latitude: 0.0, Longitude: 0.0},
			},
		}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a nil slice of pointers to structs, gracefully handles the nil slice", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64 `scrub:"true"`
			Longitude float64 `scrub:"true"`
		}
		type person struct {
			Name   string
			Age    int
			Places []*place
		}
		actual := person{
			Name:   "Testy Tester",
			Age:    26,
			Places: nil,
		}
		expected := person{
			Name:   "Testy Tester",
			Age:    26,
			Places: nil,
		}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a slice of pointers to structs, gracefully handles nil", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64 `scrub:"true"`
			Longitude float64 `scrub:"true"`
		}
		type person struct {
			Name   string
			Age    int
			Places []*place
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []*place{
				{Name: "Place 1", Latitude: 1.0, Longitude: 2.0},
				nil,
			},
		}
		expected := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []*place{
				{Name: "Place 1", Latitude: 0, Longitude: 0},
				nil,
			},
		}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing an unexported field, leaves the field unchanged", func(t *testing.T) {
		type person struct {
			Name  string `scrub:"true"`
			age   int
			Place string `scrub:"true"`
		}
		actual := person{Name: "Testy Tester", age: 26, Place: "earth"}
		expected := person{Name: "", age: 26, Place: ""}
		TaggedFields(&actual)
		assert.Equal(t, expected, actual)
	})
}

func TestNamedFields(t *testing.T) {
	primitiveTypeIsUnchanged := func(v any) bool {
		copy := v
		actual := v
		expected := copy
		NamedFields(actual, []string{}...)
		return expected == actual
	}

	t.Run("with a primitive type, leaves the value unchanged", func(t *testing.T) {
		primitives := []any{
			int(1), int8(1), int16(1), int32(1), int64(1),
			uintptr(1), uint(1), uint8(1), uint16(1), uint32(1), uint64(1),
			float32(1.0), float64(1.0),
			"hello",
		}
		for _, v := range primitives {
			vPtr := &v
			assert.True(t, primitiveTypeIsUnchanged(v))
			assert.True(t, primitiveTypeIsUnchanged(vPtr))
			assert.True(t, primitiveTypeIsUnchanged(&vPtr))
		}

		intPtr := new(int)
		intDoublePtr := &intPtr
		*intPtr = 1
		NamedFields(intPtr, []string{}...)
		assert.Equal(t, 1, *intPtr)
		assert.Equal(t, *intDoublePtr, intPtr)
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

	t.Run("with a struct containing an unexported pointer to a struct, leaves the unexported field alone", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name  string
			Age   int
			place *place
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			place: &place{
				Name:      "Testy Tester",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		expected := person{
			Name: "",
			Age:  26,
			place: &place{
				Name:      "Testy Tester",
				Latitude:  1.0,
				Longitude: 2.0,
			},
		}
		NamedFields(&actual, "Name")
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a non-nil pointer to a non-struct, scrubs the named fields", func(t *testing.T) {
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

	t.Run("with a struct containing a named slice, scrubs the slice", func(t *testing.T) {
		type person struct {
			Name   string `scrub:"true"`
			Age    int
			Places []string `scrub:"true"`
		}
		actual := person{
			Name:   "Testy Tester",
			Age:    26,
			Places: []string{"earth", "mars"},
		}
		expected := person{
			Name:   "",
			Age:    26,
			Places: nil,
		}
		NamedFields(&actual, "Name", "Places")
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a slice of structs, walks the slice and scrubs any named fields on those structs", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name   string
			Age    int
			Places []place
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []place{
				{Name: "Place 1", Latitude: 1.0, Longitude: 2.0},
				{Name: "Place 2", Latitude: 3.14, Longitude: 1.5926},
			},
		}
		expected := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []place{
				{Name: "Place 1", Latitude: 0.0, Longitude: 0.0},
				{Name: "Place 2", Latitude: 0.0, Longitude: 0.0},
			},
		}
		NamedFields(&actual, "Latitude", "Longitude")
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a slice of pointers to structs, walks the slice and scrubs any named fields on those structs", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name   string
			Age    int
			Places []*place
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []*place{
				{Name: "Place 1", Latitude: 1.0, Longitude: 2.0},
				{Name: "Place 2", Latitude: 3.14, Longitude: 1.5926},
			},
		}
		expected := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []*place{
				{Name: "Place 1", Latitude: 0.0, Longitude: 0.0},
				{Name: "Place 2", Latitude: 0.0, Longitude: 0.0},
			},
		}
		NamedFields(&actual, "Latitude", "Longitude")
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a nil slice of structs, gracefully handles the nil slice", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name   string
			Age    int
			Places []*place
		}
		actual := person{
			Name:   "Testy Tester",
			Age:    26,
			Places: nil,
		}
		expected := person{
			Name:   "Testy Tester",
			Age:    26,
			Places: nil,
		}
		NamedFields(&actual, "Latitude", "Longitude")
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing a slice of structs, gracefully handles nil", func(t *testing.T) {
		type place struct {
			Name      string
			Latitude  float64
			Longitude float64
		}
		type person struct {
			Name   string
			Age    int
			Places []*place
		}
		actual := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []*place{
				{Name: "Place 1", Latitude: 1.0, Longitude: 2.0},
				nil,
			},
		}
		expected := person{
			Name: "Testy Tester",
			Age:  26,
			Places: []*place{
				{Name: "Place 1", Latitude: 0, Longitude: 0},
				nil,
			},
		}
		NamedFields(&actual, "Latitude", "Longitude")
		assert.Equal(t, expected, actual)
	})

	t.Run("with a struct containing an unexported field, leaves the field unchanged", func(t *testing.T) {
		type person struct {
			Name  string `scrub:"true"`
			age   int
			Place string `scrub:"true"`
		}
		actual := person{Name: "Testy Tester", age: 26, Place: "earth"}
		expected := person{Name: "", age: 26, Place: ""}
		NamedFields(&actual, "Name", "Place")
		assert.Equal(t, expected, actual)
	})
}
