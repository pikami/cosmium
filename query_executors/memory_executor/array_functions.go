package memoryexecutor

import (
	"fmt"
	"reflect"

	"github.com/pikami/cosmium/parsers"
)

func array_Concat(arguments []interface{}, queryParameters map[string]interface{}, row RowType) []interface{} {
	var result []interface{}
	for _, arg := range arguments {
		array := parseArray(arg, queryParameters, row)
		result = append(result, array...)
	}
	return result
}

func array_Length(arguments []interface{}, queryParameters map[string]interface{}, row RowType) int {
	array := parseArray(arguments[0], queryParameters, row)
	if array == nil {
		return 0
	}

	return len(array)
}

func array_Slice(arguments []interface{}, queryParameters map[string]interface{}, row RowType) []interface{} {
	var ok bool
	var start int
	var length int
	array := parseArray(arguments[0], queryParameters, row)
	startEx := getFieldValue(arguments[1].(parsers.SelectItem), queryParameters, row)

	if arguments[2] != nil {
		lengthEx := getFieldValue(arguments[2].(parsers.SelectItem), queryParameters, row)

		if length, ok = lengthEx.(int); !ok {
			fmt.Println("array_Slice - got length parameters of wrong type")
			return []interface{}{}
		}
	}

	if start, ok = startEx.(int); !ok {
		fmt.Println("array_Slice - got start parameters of wrong type")
		return []interface{}{}
	}

	if start < 0 {
		start = len(array) + start
	}

	if start < 0 {
		start = 0
	}

	if array == nil || start >= len(array) {
		return []interface{}{}
	}

	end := start + length
	if end > len(array) {
		end = len(array)
	}
	return array[start:end]
}

func set_Intersect(arguments []interface{}, queryParameters map[string]interface{}, row RowType) []interface{} {
	set1 := parseArray(arguments[0], queryParameters, row)
	set2 := parseArray(arguments[1], queryParameters, row)

	intersection := make(map[interface{}]struct{})
	if set1 == nil || set2 == nil {
		return []interface{}{}
	}

	for _, item := range set1 {
		intersection[item] = struct{}{}
	}

	var result []interface{}
	for _, item := range set2 {
		if _, exists := intersection[item]; exists {
			result = append(result, item)
		}
	}

	return result
}

func set_Union(arguments []interface{}, queryParameters map[string]interface{}, row RowType) []interface{} {
	set1 := parseArray(arguments[0], queryParameters, row)
	set2 := parseArray(arguments[1], queryParameters, row)

	var result []interface{}
	union := make(map[interface{}]struct{})
	for _, item := range set1 {
		if _, ok := union[item]; !ok {
			union[item] = struct{}{}
			result = append(result, item)
		}
	}

	for _, item := range set2 {
		if _, ok := union[item]; !ok {
			union[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func parseArray(argument interface{}, queryParameters map[string]interface{}, row RowType) []interface{} {
	exItem := argument.(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	arrValue := reflect.ValueOf(ex)
	if arrValue.Kind() != reflect.Slice {
		fmt.Println("parseArray got parameters of wrong type")
		return nil
	}

	result := make([]interface{}, arrValue.Len())

	for i := 0; i < arrValue.Len(); i++ {
		result[i] = arrValue.Index(i).Interface()
	}

	return result
}
