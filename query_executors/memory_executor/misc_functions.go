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

func (r rowContext) misc_Iif(arguments []interface{}) interface{} {
	if len(arguments) != 3 {
		return nil
	}

	condition := r.resolveSelectItem(arguments[0].(parsers.SelectItem))
	if condition != nil && condition == true {
		return r.resolveSelectItem(arguments[1].(parsers.SelectItem))
	}

	return r.resolveSelectItem(arguments[2].(parsers.SelectItem))
}

func (r rowContext) misc_UDF(arguments []interface{}) interface{} {
	return "TODO"
}
