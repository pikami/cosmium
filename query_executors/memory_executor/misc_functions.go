package memoryexecutor

import (
	"github.com/pikami/cosmium/parsers"
)

func misc_In(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	value := getFieldValue(arguments[0].(parsers.SelectItem), queryParameters, row)

	for i := 1; i < len(arguments); i++ {
		compareValue := getFieldValue(arguments[i].(parsers.SelectItem), queryParameters, row)
		if compareValues(value, compareValue) == 0 {
			return true
		}
	}

	return false
}
