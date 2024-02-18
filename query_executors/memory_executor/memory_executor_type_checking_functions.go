package memoryexecutor

import (
	"github.com/pikami/cosmium/parsers"
)

func typeChecking_IsDefined(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	return ex != nil
}
