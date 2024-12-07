package memoryexecutor

import (
	"math"

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
	selectExpression := arguments[0].(parsers.SelectItem)
	max := 0.0
	count := 0

	for _, item := range r.grouppedRows {
		value := item.resolveSelectItem(selectExpression)
		if numericValue, ok := value.(float64); ok {
			if numericValue > max {
				max = numericValue
			}
			count++
		} else if numericValue, ok := value.(int); ok {
			if float64(numericValue) > max {
				max = float64(numericValue)
			}
			count++
		}
	}

	if count > 0 {
		return max
	} else {
		return nil
	}
}

func (r rowContext) aggregate_Min(arguments []interface{}) interface{} {
	selectExpression := arguments[0].(parsers.SelectItem)
	min := math.MaxFloat64
	count := 0

	for _, item := range r.grouppedRows {
		value := item.resolveSelectItem(selectExpression)
		if numericValue, ok := value.(float64); ok {
			if numericValue < min {
				min = numericValue
			}
			count++
		} else if numericValue, ok := value.(int); ok {
			if float64(numericValue) < min {
				min = float64(numericValue)
			}
			count++
		}
	}

	if count > 0 {
		return min
	} else {
		return nil
	}
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
