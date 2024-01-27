package scrub

import (
	"reflect"
	"slices"
)

// TaggedFields takes a struct and recursively sets all fields annotated with a `scrub:"true"`
// struct tag to their zero value. This is useful when you control the struct definition.
//
// This function is a no-op for non-struct types.
func TaggedFields(src any) {
	scrub(src, func(field reflect.StructField) bool {
		return field.Tag.Get("scrub") == "true"
	})
}

// NamedFields takes a struct and sets all fields with the given names to their zero value. This is useful
// when you want to scrub a struct type from a package that you don't control.
//
// This function is a no-op for non-struct types.
func NamedFields(src any, names ...string) {
	scrub(src, func(field reflect.StructField) bool {
		return slices.Contains(names, field.Name)
	})
}

func scrub(src any, shouldScrubFn func(field reflect.StructField) bool) {
	if src == nil {
		return
	}
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		switch field.Kind() {
		case reflect.Struct:
			structField := v.Type().Field(i)
			if shouldScrubFn(structField) {
				zero := reflect.Zero(field.Type())
				field.Set(zero)
			} else {
				scrub(field.Addr().Interface(), shouldScrubFn)
			}
		case reflect.Ptr:
			if field.IsNil() {
				continue
			}
			if field.Elem().Kind() == reflect.Struct {
				structField := v.Type().Field(i)
				if shouldScrubFn(structField) {
					zero := reflect.Zero(field.Type())
					field.Set(zero)
				} else {
					scrub(field.Interface(), shouldScrubFn)
				}
			} else {
				continue
			}
		default:
			structField := v.Type().Field(i)
			if shouldScrubFn(structField) {
				zero := reflect.Zero(field.Type())
				field.Set(zero)
			}
		}
	}
}
