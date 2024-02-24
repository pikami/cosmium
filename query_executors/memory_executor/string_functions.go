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

func strings_Left(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	var ok bool
	var length int
	str := parseString(arguments[0], queryParameters, row)
	lengthEx := getFieldValue(arguments[1].(parsers.SelectItem), queryParameters, row)

	if length, ok = lengthEx.(int); !ok {
		fmt.Println("strings_Left - got parameters of wrong type")
		return ""
	}

	if length <= 0 {
		return ""
	}

	if len(str) <= length {
		return str
	}

	return str[:length]
}

func strings_Length(arguments []interface{}, queryParameters map[string]interface{}, row RowType) int {
	str := parseString(arguments[0], queryParameters, row)
	return len(str)
}

func strings_LTrim(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	str := parseString(arguments[0], queryParameters, row)
	return strings.TrimLeft(str, " ")
}

func strings_Replace(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	str := parseString(arguments[0], queryParameters, row)
	oldStr := parseString(arguments[1], queryParameters, row)
	newStr := parseString(arguments[2], queryParameters, row)
	return strings.Replace(str, oldStr, newStr, -1)
}

func strings_Replicate(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	var ok bool
	var times int
	str := parseString(arguments[0], queryParameters, row)
	timesEx := getFieldValue(arguments[1].(parsers.SelectItem), queryParameters, row)

	if times, ok = timesEx.(int); !ok {
		fmt.Println("strings_Replicate - got parameters of wrong type")
		return ""
	}

	if times <= 0 {
		return ""
	}

	if len(str) <= times {
		return str
	}

	return strings.Repeat(str, times)
}

func strings_Reverse(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	str := parseString(arguments[0], queryParameters, row)
	runes := []rune(str)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func strings_Right(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	var ok bool
	var length int
	str := parseString(arguments[0], queryParameters, row)
	lengthEx := getFieldValue(arguments[1].(parsers.SelectItem), queryParameters, row)

	if length, ok = lengthEx.(int); !ok {
		fmt.Println("strings_Right - got parameters of wrong type")
		return ""
	}

	if length <= 0 {
		return ""
	}

	if len(str) <= length {
		return str
	}

	return str[len(str)-length:]
}

func strings_RTrim(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	str := parseString(arguments[0], queryParameters, row)
	return strings.TrimRight(str, " ")
}

func strings_Substring(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	var ok bool
	var startPos int
	var length int
	str := parseString(arguments[0], queryParameters, row)
	startPosEx := getFieldValue(arguments[1].(parsers.SelectItem), queryParameters, row)
	lengthEx := getFieldValue(arguments[2].(parsers.SelectItem), queryParameters, row)

	if startPos, ok = startPosEx.(int); !ok {
		fmt.Println("strings_Substring - got start parameters of wrong type")
		return ""
	}
	if length, ok = lengthEx.(int); !ok {
		fmt.Println("strings_Substring - got length parameters of wrong type")
		return ""
	}

	if startPos >= len(str) {
		return ""
	}

	endPos := startPos + length
	if endPos > len(str) {
		endPos = len(str)
	}

	return str[startPos:endPos]
}

func strings_Trim(arguments []interface{}, queryParameters map[string]interface{}, row RowType) string {
	str := parseString(arguments[0], queryParameters, row)
	return strings.TrimSpace(str)
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
