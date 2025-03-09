package structhidrators

import (
	"reflect"

	"github.com/pikami/cosmium/internal/datastore"
)

func Hidrate(input interface{}) interface{} {
	if reflect.TypeOf(input) == reflect.TypeOf(datastore.Collection{}) {
		return hidrate(input, defaultCollection)
	}
	return input
}

func hidrate(input interface{}, defaults interface{}) interface{} {
	inputVal := reflect.ValueOf(input)
	defaultsVal := reflect.ValueOf(defaults)

	if inputVal.Kind() != reflect.Struct || defaultsVal.Kind() != reflect.Struct {
		panic("Both input and defaults must be structs")
	}

	output := reflect.New(inputVal.Type()).Elem()
	for i := 0; i < inputVal.NumField(); i++ {
		inputField := inputVal.Field(i)
		defaultField := defaultsVal.Field(i)

		if inputField.Kind() == reflect.Struct && defaultField.Kind() == reflect.Struct {
			filledNested := hidrate(inputField.Interface(), defaultField.Interface())
			output.Field(i).Set(reflect.ValueOf(filledNested))
		} else {
			if isEmptyValue(inputField) {
				output.Field(i).Set(defaultField)
			} else {
				output.Field(i).Set(inputField)
			}
		}
	}

	return output.Interface()
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
