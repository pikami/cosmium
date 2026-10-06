package memoryexecutor

import (
	"strings"

	"github.com/pikami/cosmium/parsers"
)

func (r rowContext) aggregate_Avg(arguments []interface{}) interface{} {
	selectExpression := arguments[0].(parsers.SelectItem)
	sum := 0.0
	count := 0

	for _, item := range r.grouppedRows {
		value := item.resolveSelectItem(selectExpression)
		if numericValue, ok := value.(float64); ok {
			sum += numericValue
			count++
		} else if numericValue, ok := value.(int); ok {
			sum += float64(numericValue)
			count++
		}
	}

	if count > 0 {
		return sum / float64(count)
	} else {
		return nil
	}
}

func (r rowContext) aggregate_Count(arguments []interface{}) interface{} {
	selectExpression := arguments[0].(parsers.SelectItem)
	count := 0

	for _, item := range r.grouppedRows {
		value := item.resolveSelectItem(selectExpression)
		if value != nil {
			count++
		}
	}

	return count
}

func (r rowContext) aggregate_Max(arguments []interface{}) interface{} {
	return r.aggregateMinMax(arguments, true)
}

func (r rowContext) aggregate_Min(arguments []interface{}) interface{} {
	return r.aggregateMinMax(arguments, false)
}

func (r rowContext) aggregateMinMax(arguments []interface{}, isMax bool) interface{} {
	selectExpression := arguments[0].(parsers.SelectItem)
	var extremeNumber float64
	hasNumber := false
	var extremeString string
	hasString := false

	for _, item := range r.grouppedRows {
		value := item.resolveSelectItem(selectExpression)

		if numericValue, ok := aggregateNumber(value); ok {
			if !hasNumber || numberIsExtreme(numericValue, extremeNumber, isMax) {
				extremeNumber = numericValue
				hasNumber = true
			}
			continue
		}

		if strValue, ok := value.(string); ok {
			if !hasString || stringIsExtreme(strValue, extremeString, isMax) {
				extremeString = strValue
				hasString = true
			}
		}
	}

	if hasNumber {
		return extremeNumber
	}
	if hasString {
		return extremeString
	}
	return nil
}

func aggregateNumber(value interface{}) (float64, bool) {
	switch numericValue := value.(type) {
	case float64:
		return numericValue, true
	case int:
		return float64(numericValue), true
	default:
		return 0, false
	}
}

func numberIsExtreme(candidate, current float64, isMax bool) bool {
	if isMax {
		return candidate > current
	}
	return candidate < current
}

func stringIsExtreme(candidate, current string, isMax bool) bool {
	comparison := strings.Compare(candidate, current)
	if isMax {
		return comparison > 0
	}
	return comparison < 0
}

func (r rowContext) aggregate_Sum(arguments []interface{}) interface{} {
	selectExpression := arguments[0].(parsers.SelectItem)
	sum := 0.0
	count := 0

	for _, item := range r.grouppedRows {
		value := item.resolveSelectItem(selectExpression)
		if numericValue, ok := value.(float64); ok {
			sum += numericValue
			count++
		} else if numericValue, ok := value.(int); ok {
			sum += float64(numericValue)
			count++
		}
	}

	if count > 0 {
		return sum
	} else {
		return nil
	}
}
