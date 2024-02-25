package memoryexecutor

import (
	"fmt"
	"strings"

	"github.com/pikami/cosmium/parsers"
)

func (c memoryExecutorContext) strings_StringEquals(arguments []interface{}, row RowType) bool {
	str1 := c.parseString(arguments[0], row)
	str2 := c.parseString(arguments[1], row)
	ignoreCase := c.getBoolFlag(arguments, row)

	if ignoreCase {
		return strings.EqualFold(str1, str2)
	}

	return str1 == str2
}

func (c memoryExecutorContext) strings_Contains(arguments []interface{}, row RowType) bool {
	str1 := c.parseString(arguments[0], row)
	str2 := c.parseString(arguments[1], row)
	ignoreCase := c.getBoolFlag(arguments, row)

	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	return strings.Contains(str1, str2)
}

func (c memoryExecutorContext) strings_EndsWith(arguments []interface{}, row RowType) bool {
	str1 := c.parseString(arguments[0], row)
	str2 := c.parseString(arguments[1], row)
	ignoreCase := c.getBoolFlag(arguments, row)

	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	return strings.HasSuffix(str1, str2)
}

func (c memoryExecutorContext) strings_StartsWith(arguments []interface{}, row RowType) bool {
	str1 := c.parseString(arguments[0], row)
	str2 := c.parseString(arguments[1], row)
	ignoreCase := c.getBoolFlag(arguments, row)

	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	return strings.HasPrefix(str1, str2)
}

func (c memoryExecutorContext) strings_Concat(arguments []interface{}, row RowType) string {
	result := ""

	for _, arg := range arguments {
		if selectItem, ok := arg.(parsers.SelectItem); ok {
			value := c.getFieldValue(selectItem, row)
			result += convertToString(value)
		}
	}

	return result
}

func (c memoryExecutorContext) strings_IndexOf(arguments []interface{}, row RowType) int {
	str1 := c.parseString(arguments[0], row)
	str2 := c.parseString(arguments[1], row)

	start := 0
	if len(arguments) > 2 && arguments[2] != nil {
		if startPos, ok := c.getFieldValue(arguments[2].(parsers.SelectItem), row).(int); ok {
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

func (c memoryExecutorContext) strings_ToString(arguments []interface{}, row RowType) string {
	value := c.getFieldValue(arguments[0].(parsers.SelectItem), row)
	return convertToString(value)
}

func (c memoryExecutorContext) strings_Upper(arguments []interface{}, row RowType) string {
	value := c.getFieldValue(arguments[0].(parsers.SelectItem), row)
	return strings.ToUpper(convertToString(value))
}

func (c memoryExecutorContext) strings_Lower(arguments []interface{}, row RowType) string {
	value := c.getFieldValue(arguments[0].(parsers.SelectItem), row)
	return strings.ToLower(convertToString(value))
}

func (c memoryExecutorContext) strings_Left(arguments []interface{}, row RowType) string {
	var ok bool
	var length int
	str := c.parseString(arguments[0], row)
	lengthEx := c.getFieldValue(arguments[1].(parsers.SelectItem), row)

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

func (c memoryExecutorContext) strings_Length(arguments []interface{}, row RowType) int {
	str := c.parseString(arguments[0], row)
	return len(str)
}

func (c memoryExecutorContext) strings_LTrim(arguments []interface{}, row RowType) string {
	str := c.parseString(arguments[0], row)
	return strings.TrimLeft(str, " ")
}

func (c memoryExecutorContext) strings_Replace(arguments []interface{}, row RowType) string {
	str := c.parseString(arguments[0], row)
	oldStr := c.parseString(arguments[1], row)
	newStr := c.parseString(arguments[2], row)
	return strings.Replace(str, oldStr, newStr, -1)
}

func (c memoryExecutorContext) strings_Replicate(arguments []interface{}, row RowType) string {
	var ok bool
	var times int
	str := c.parseString(arguments[0], row)
	timesEx := c.getFieldValue(arguments[1].(parsers.SelectItem), row)

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

func (c memoryExecutorContext) strings_Reverse(arguments []interface{}, row RowType) string {
	str := c.parseString(arguments[0], row)
	runes := []rune(str)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func (c memoryExecutorContext) strings_Right(arguments []interface{}, row RowType) string {
	var ok bool
	var length int
	str := c.parseString(arguments[0], row)
	lengthEx := c.getFieldValue(arguments[1].(parsers.SelectItem), row)

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

func (c memoryExecutorContext) strings_RTrim(arguments []interface{}, row RowType) string {
	str := c.parseString(arguments[0], row)
	return strings.TrimRight(str, " ")
}

func (c memoryExecutorContext) strings_Substring(arguments []interface{}, row RowType) string {
	var ok bool
	var startPos int
	var length int
	str := c.parseString(arguments[0], row)
	startPosEx := c.getFieldValue(arguments[1].(parsers.SelectItem), row)
	lengthEx := c.getFieldValue(arguments[2].(parsers.SelectItem), row)

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

func (c memoryExecutorContext) strings_Trim(arguments []interface{}, row RowType) string {
	str := c.parseString(arguments[0], row)
	return strings.TrimSpace(str)
}

func (c memoryExecutorContext) getBoolFlag(arguments []interface{}, row RowType) bool {
	ignoreCase := false
	if len(arguments) > 2 && arguments[2] != nil {
		ignoreCaseItem := arguments[2].(parsers.SelectItem)
		if value, ok := c.getFieldValue(ignoreCaseItem, row).(bool); ok {
			ignoreCase = value
		}
	}

	return ignoreCase
}

func (c memoryExecutorContext) parseString(argument interface{}, row RowType) string {
	exItem := argument.(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)
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
