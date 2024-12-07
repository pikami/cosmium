package memoryexecutor

import (
	"reflect"

	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/parsers"
)

func (r rowContext) array_Concat(arguments []interface{}) []interface{} {
	var result []interface{}
	for _, arg := range arguments {
		array := r.parseArray(arg)
		result = append(result, array...)
	}
	return result
}

func (r rowContext) array_Length(arguments []interface{}) int {
	array := r.parseArray(arguments[0])
	if array == nil {
		return 0
	}

	return len(array)
}

func (r rowContext) array_Slice(arguments []interface{}) []interface{} {
	var ok bool
	var start int
	var length int
	array := r.parseArray(arguments[0])
	startEx := r.resolveSelectItem(arguments[1].(parsers.SelectItem))

	if arguments[2] != nil {
		lengthEx := r.resolveSelectItem(arguments[2].(parsers.SelectItem))

		if length, ok = lengthEx.(int); !ok {
			logger.Error("array_Slice - got length parameters of wrong type")
			return []interface{}{}
		}
	}

	if start, ok = startEx.(int); !ok {
		logger.Error("array_Slice - got start parameters of wrong type")
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

func (r rowContext) set_Intersect(arguments []interface{}) []interface{} {
	set1 := r.parseArray(arguments[0])
	set2 := r.parseArray(arguments[1])

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

func (r rowContext) set_Union(arguments []interface{}) []interface{} {
	set1 := r.parseArray(arguments[0])
	set2 := r.parseArray(arguments[1])

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

func (r rowContext) parseArray(argument interface{}) []interface{} {
	exItem := argument.(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	arrValue := reflect.ValueOf(ex)
	if arrValue.Kind() != reflect.Slice {
		logger.Error("parseArray got parameters of wrong type")
		return nil
	}

	result := make([]interface{}, arrValue.Len())

	for i := 0; i < arrValue.Len(); i++ {
		result[i] = arrValue.Index(i).Interface()
	}

	return result
}
