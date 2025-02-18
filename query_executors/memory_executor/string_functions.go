package memoryexecutor

import (
	"fmt"
	"strings"

	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/parsers"
)

func (r rowContext) strings_StringEquals(arguments []interface{}) bool {
	str1, str1ok := r.parseString(arguments[0])
	str2, str2ok := r.parseString(arguments[1])
	ignoreCase := r.getBoolFlag(arguments)

	if !str1ok || !str2ok {
		return false
	}

	if ignoreCase {
		return strings.EqualFold(str1, str2)
	}

	return str1 == str2
}

func (r rowContext) strings_Contains(arguments []interface{}) bool {
	str1, str1ok := r.parseString(arguments[0])
	str2, str2ok := r.parseString(arguments[1])
	ignoreCase := r.getBoolFlag(arguments)

	if !str1ok || !str2ok {
		return false
	}

	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	return strings.Contains(str1, str2)
}

func (r rowContext) strings_EndsWith(arguments []interface{}) bool {
	str1, str1ok := r.parseString(arguments[0])
	str2, str2ok := r.parseString(arguments[1])
	ignoreCase := r.getBoolFlag(arguments)

	if !str1ok || !str2ok {
		return false
	}

	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	return strings.HasSuffix(str1, str2)
}

func (r rowContext) strings_StartsWith(arguments []interface{}) bool {
	str1, str1ok := r.parseString(arguments[0])
	str2, str2ok := r.parseString(arguments[1])
	ignoreCase := r.getBoolFlag(arguments)

	if !str1ok || !str2ok {
		return false
	}

	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	return strings.HasPrefix(str1, str2)
}

func (r rowContext) strings_Concat(arguments []interface{}) string {
	result := ""

	for _, arg := range arguments {
		if selectItem, ok := arg.(parsers.SelectItem); ok {
			value := r.resolveSelectItem(selectItem)
			result += convertToString(value)
		}
	}

	return result
}

func (r rowContext) strings_IndexOf(arguments []interface{}) int {
	str1, str1ok := r.parseString(arguments[0])
	str2, str2ok := r.parseString(arguments[1])

	if !str1ok || !str2ok {
		return -1
	}

	start := 0
	if len(arguments) > 2 && arguments[2] != nil {
		if startPos, ok := r.resolveSelectItem(arguments[2].(parsers.SelectItem)).(int); ok {
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

func (r rowContext) strings_ToString(arguments []interface{}) string {
	value := r.resolveSelectItem(arguments[0].(parsers.SelectItem))
	return convertToString(value)
}

func (r rowContext) strings_Upper(arguments []interface{}) string {
	value := r.resolveSelectItem(arguments[0].(parsers.SelectItem))
	return strings.ToUpper(convertToString(value))
}

func (r rowContext) strings_Lower(arguments []interface{}) string {
	value := r.resolveSelectItem(arguments[0].(parsers.SelectItem))
	return strings.ToLower(convertToString(value))
}

func (r rowContext) strings_Left(arguments []interface{}) string {
	var ok bool
	var length int
	str, strOk := r.parseString(arguments[0])
	lengthEx := r.resolveSelectItem(arguments[1].(parsers.SelectItem))

	if !strOk {
		return ""
	}

	if length, ok = lengthEx.(int); !ok {
		logger.ErrorLn("strings_Left - got parameters of wrong type")
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

func (r rowContext) strings_Length(arguments []interface{}) int {
	str, strOk := r.parseString(arguments[0])
	if !strOk {
		return 0
	}

	return len(str)
}

func (r rowContext) strings_LTrim(arguments []interface{}) string {
	str, strOk := r.parseString(arguments[0])
	if !strOk {
		return ""
	}

	return strings.TrimLeft(str, " ")
}

func (r rowContext) strings_Replace(arguments []interface{}) string {
	str, strOk := r.parseString(arguments[0])
	oldStr, oldStrOk := r.parseString(arguments[1])
	newStr, newStrOk := r.parseString(arguments[2])

	if !strOk || !oldStrOk || !newStrOk {
		return ""
	}

	return strings.Replace(str, oldStr, newStr, -1)
}

func (r rowContext) strings_Replicate(arguments []interface{}) string {
	var ok bool
	var times int
	str, strOk := r.parseString(arguments[0])
	timesEx := r.resolveSelectItem(arguments[1].(parsers.SelectItem))

	if !strOk {
		return ""
	}

	if times, ok = timesEx.(int); !ok {
		logger.ErrorLn("strings_Replicate - got parameters of wrong type")
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

func (r rowContext) strings_Reverse(arguments []interface{}) string {
	str, strOk := r.parseString(arguments[0])
	runes := []rune(str)

	if !strOk {
		return ""
	}

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func (r rowContext) strings_Right(arguments []interface{}) string {
	var ok bool
	var length int
	str, strOk := r.parseString(arguments[0])
	lengthEx := r.resolveSelectItem(arguments[1].(parsers.SelectItem))

	if !strOk {
		return ""
	}

	if length, ok = lengthEx.(int); !ok {
		logger.ErrorLn("strings_Right - got parameters of wrong type")
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

func (r rowContext) strings_RTrim(arguments []interface{}) string {
	str, strOk := r.parseString(arguments[0])
	if !strOk {
		return ""
	}

	return strings.TrimRight(str, " ")
}

func (r rowContext) strings_Substring(arguments []interface{}) string {
	var ok bool
	var startPos int
	var length int
	str, strOk := r.parseString(arguments[0])
	startPosEx := r.resolveSelectItem(arguments[1].(parsers.SelectItem))
	lengthEx := r.resolveSelectItem(arguments[2].(parsers.SelectItem))

	if !strOk {
		return ""
	}

	if startPos, ok = startPosEx.(int); !ok {
		logger.ErrorLn("strings_Substring - got start parameters of wrong type")
		return ""
	}
	if length, ok = lengthEx.(int); !ok {
		logger.ErrorLn("strings_Substring - got length parameters of wrong type")
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

func (r rowContext) strings_Trim(arguments []interface{}) string {
	str, strOk := r.parseString(arguments[0])
	if !strOk {
		return ""
	}

	return strings.TrimSpace(str)
}

func (r rowContext) getBoolFlag(arguments []interface{}) bool {
	ignoreCase := false
	if len(arguments) > 2 && arguments[2] != nil {
		ignoreCaseItem := arguments[2].(parsers.SelectItem)
		if value, ok := r.resolveSelectItem(ignoreCaseItem).(bool); ok {
			ignoreCase = value
		}
	}

	return ignoreCase
}

func (r rowContext) parseString(argument interface{}) (value string, ok bool) {
	exItem := argument.(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)
	if str1, ok := ex.(string); ok {
		return str1, true
	}

	logger.ErrorLn("StringEquals got parameters of wrong type")
	return "", false
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
