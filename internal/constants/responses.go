package constants

import (
	"github.com/gin-gonic/gin"
)

var QueryPlanResponse = gin.H{
	"partitionedQueryExecutionInfoVersion": 2,
	"queryInfo": map[string]interface{}{
		"distinctType":                "None",
		"top":                         nil,
		"offset":                      nil,
		"limit":                       nil,
		"orderBy":                     []interface{}{},
		"orderByExpressions":          []interface{}{},
		"groupByExpressions":          []interface{}{},
		"groupByAliases":              []interface{}{},
		"aggregates":                  []interface{}{},
		"groupByAliasToAggregateType": map[string]interface{}{},
		"rewrittenQuery":              "",
		"hasSelectValue":              false,
		"dCountInfo":                  nil,
	},
	"queryRanges": []interface{}{
		map[string]interface{}{
			"min":            "",
			"max":            "FF",
			"isMinInclusive": true,
			"isMaxInclusive": false,
		},
	},
}

var UnknownErrorResponse = gin.H{"message": "Unknown error"}
var NotFoundResponse = gin.H{"message": "NotFound"}
var ConflictResponse = gin.H{"message": "Conflict"}
var BadRequestResponse = gin.H{"message": "BadRequest"}
