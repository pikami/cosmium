package memoryexecutor

import (
	"reflect"

	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/parsers"
)

func (c memoryExecutorContext) array_Concat(arguments []interface{}, row RowType) []interface{} {
	var result []interface{}
	for _, arg := range arguments {
		array := c.parseArray(arg, row)
		result = append(result, array...)
	}
	return result
}

func (c memoryExecutorContext) array_Length(arguments []interface{}, row RowType) int {
	array := c.parseArray(arguments[0], row)
	if array == nil {
		return 0
	}

	return len(array)
}

func (c memoryExecutorContext) array_Slice(arguments []interface{}, row RowType) []interface{} {
	var ok bool
	var start int
	var length int
	array := c.parseArray(arguments[0], row)
	startEx := c.getFieldValue(arguments[1].(parsers.SelectItem), row)

	if arguments[2] != nil {
		lengthEx := c.getFieldValue(arguments[2].(parsers.SelectItem), row)

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

func (c memoryExecutorContext) set_Intersect(arguments []interface{}, row RowType) []interface{} {
	set1 := c.parseArray(arguments[0], row)
	set2 := c.parseArray(arguments[1], row)

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

func (c memoryExecutorContext) set_Union(arguments []interface{}, row RowType) []interface{} {
	set1 := c.parseArray(arguments[0], row)
	set2 := c.parseArray(arguments[1], row)

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

func (c memoryExecutorContext) parseArray(argument interface{}, row RowType) []interface{} {
	exItem := argument.(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

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
