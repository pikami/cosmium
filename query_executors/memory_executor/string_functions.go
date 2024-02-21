package memoryexecutor

import (
	"fmt"
	"strings"

	"github.com/pikami/cosmium/parsers"
)

func strings_StringEquals(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	ignoreCase := false
	if len(arguments) > 2 && arguments[2] != nil {
		ignoreCaseItem := arguments[2].(parsers.SelectItem)
		if value, ok := getFieldValue(ignoreCaseItem, queryParameters, row).(bool); ok {
			ignoreCase = value
		}
	}

	ex1Item := arguments[0].(parsers.SelectItem)
	ex2Item := arguments[1].(parsers.SelectItem)

	ex1 := getFieldValue(ex1Item, queryParameters, row)
	ex2 := getFieldValue(ex2Item, queryParameters, row)

	var ok bool
	var str1 string
	var str2 string

	if str1, ok = ex1.(string); !ok {
		fmt.Println("StringEquals got parameters of wrong type")
	}
	if str2, ok = ex2.(string); !ok {
		fmt.Println("StringEquals got parameters of wrong type")
	}

	if ignoreCase {
		return strings.EqualFold(str1, str2)
	}

	return str1 == str2
}

func strings_Concat(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	result := ""

	for _, arg := range arguments {
		if selectItem, ok := arg.(parsers.SelectItem); ok {
			value := getFieldValue(selectItem, queryParameters, row)
			result += convertToString(value)
		}
	}

	return result
}

func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	}
	return ""
}
