package memoryexecutor

import (
	"github.com/pikami/cosmium/parsers"
)

func (c memoryExecutorContext) misc_In(arguments []interface{}, row RowType) bool {
	value := c.getFieldValue(arguments[0].(parsers.SelectItem), row)

	for i := 1; i < len(arguments); i++ {
		compareValue := c.getFieldValue(arguments[i].(parsers.SelectItem), row)
		if compareValues(value, compareValue) == 0 {
			return true
		}
	}

	return false
}
