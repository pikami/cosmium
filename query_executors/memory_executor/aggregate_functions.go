package memoryexecutor

import (
	"math"

	"github.com/pikami/cosmium/parsers"
)

func (c memoryExecutorContext) aggregate_Avg(arguments []interface{}, row RowType) interface{} {
	selectExpression := arguments[0].(parsers.SelectItem)
	sum := 0.0
	count := 0

	if array, isArray := row.([]RowType); isArray {
		for _, item := range array {
			value := c.getFieldValue(selectExpression, item)
			if numericValue, ok := value.(float64); ok {
				sum += numericValue
				count++
			} else if numericValue, ok := value.(int); ok {
				sum += float64(numericValue)
				count++
			}
		}
	}

	if count > 0 {
		return sum / float64(count)
	} else {
		return nil
	}
}

func (c memoryExecutorContext) aggregate_Count(arguments []interface{}, row RowType) interface{} {
	selectExpression := arguments[0].(parsers.SelectItem)
	count := 0

	if array, isArray := row.([]RowType); isArray {
		for _, item := range array {
			value := c.getFieldValue(selectExpression, item)
			if value != nil {
				count++
			}
		}
	}

	return count
}

func (c memoryExecutorContext) aggregate_Max(arguments []interface{}, row RowType) interface{} {
	selectExpression := arguments[0].(parsers.SelectItem)
	max := 0.0
	count := 0

	if array, isArray := row.([]RowType); isArray {
		for _, item := range array {
			value := c.getFieldValue(selectExpression, item)
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
	}

	if count > 0 {
		return max
	} else {
		return nil
	}
}

func (c memoryExecutorContext) aggregate_Min(arguments []interface{}, row RowType) interface{} {
	selectExpression := arguments[0].(parsers.SelectItem)
	min := math.MaxFloat64
	count := 0

	if array, isArray := row.([]RowType); isArray {
		for _, item := range array {
			value := c.getFieldValue(selectExpression, item)
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
	}

	if count > 0 {
		return min
	} else {
		return nil
	}
}

func (c memoryExecutorContext) aggregate_Sum(arguments []interface{}, row RowType) interface{} {
	selectExpression := arguments[0].(parsers.SelectItem)
	sum := 0.0
	count := 0

	if array, isArray := row.([]RowType); isArray {
		for _, item := range array {
			value := c.getFieldValue(selectExpression, item)
			if numericValue, ok := value.(float64); ok {
				sum += numericValue
				count++
			} else if numericValue, ok := value.(int); ok {
				sum += float64(numericValue)
				count++
			}
		}
	}

	if count > 0 {
		return sum
	} else {
		return nil
	}
}
