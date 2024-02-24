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

func strings_IndexOf(arguments []interface{}, queryParameters map[string]interface{}, row RowType) int {
	str1 := parseString(arguments[0], queryParameters, row)
	str2 := parseString(arguments[1], queryParameters, row)

	start := 0
	if len(arguments) > 2 && arguments[2] != nil {
		if startPos, ok := getFieldValue(arguments[2].(parsers.SelectItem), queryParameters, row).(int); ok {
			start = startPos
		}
	}

	if len(str1) <= start {
		return -1
	}

	str1 = str1[start:]
	result := strings.Index(str1, str2)

	if result == -1 {
		return result
	} else {
		return result + start
	}
}

func strings_ToString(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	value := getFieldValue(arguments[0].(parsers.SelectItem), queryParameters, row)
	return convertToString(value)
}

func strings_Upper(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	value := getFieldValue(arguments[0].(parsers.SelectItem), queryParameters, row)
	return strings.ToUpper(convertToString(value))
}

func strings_Lower(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	value := getFieldValue(arguments[0].(parsers.SelectItem), queryParameters, row)
	return strings.ToLower(convertToString(value))
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
		return str1
	}

	fmt.Println("StringEquals got parameters of wrong type")
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
	case bool:
		return fmt.Sprintf("%t", v)
	}
	return ""
}
