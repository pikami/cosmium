package memoryexecutor

import (
	"fmt"
	"strings"

	"github.com/pikami/cosmium/parsers"
)

func strings_StringEquals(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	str1 := parseString(arguments[0], queryParameters, row)
	str2 := parseString(arguments[1], queryParameters, row)
	ignoreCase := getBoolFlag(arguments, queryParameters, row)

	if ignoreCase {
		return strings.EqualFold(str1, str2)
	}

	return str1 == str2
}

func strings_Contains(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	str1 := parseString(arguments[0], queryParameters, row)
	str2 := parseString(arguments[1], queryParameters, row)
	ignoreCase := getBoolFlag(arguments, queryParameters, row)

	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	return strings.Contains(str1, str2)
}

func strings_EndsWith(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	str1 := parseString(arguments[0], queryParameters, row)
	str2 := parseString(arguments[1], queryParameters, row)
	ignoreCase := getBoolFlag(arguments, queryParameters, row)

	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	return strings.HasSuffix(str1, str2)
}

func strings_StartsWith(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	str1 := parseString(arguments[0], queryParameters, row)
	str2 := parseString(arguments[1], queryParameters, row)
	ignoreCase := getBoolFlag(arguments, queryParameters, row)

	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	return strings.HasPrefix(str1, str2)
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

func getBoolFlag(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	ignoreCase := false
	if len(arguments) > 2 && arguments[2] != nil {
		ignoreCaseItem := arguments[2].(parsers.SelectItem)
		if value, ok := getFieldValue(ignoreCaseItem, queryParameters, row).(bool); ok {
			ignoreCase = value
		}
	}

	return ignoreCase
}

func parseString(argument interface{}, queryParameters map[string]interface{}, row RowType) string {
	exItem := argument.(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)
	if str1, ok := ex.(string); ok {
		fmt.Println("StringEquals got parameters of wrong type")
		return str1
	}

	return ""
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
