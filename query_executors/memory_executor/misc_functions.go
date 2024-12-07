package memoryexecutor

import (
	"github.com/pikami/cosmium/parsers"
)

func (r rowContext) misc_In(arguments []interface{}) bool {
	value := r.resolveSelectItem(arguments[0].(parsers.SelectItem))

	for i := 1; i < len(arguments); i++ {
		compareValue := r.resolveSelectItem(arguments[i].(parsers.SelectItem))
		if compareValues(value, compareValue) == 0 {
			return true
		}
	}

	return false
}
