package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetServerInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"_self":     "",
		"id":        h.config.DatabaseAccount,
		"_rid":      fmt.Sprintf("%s.%s", h.config.DatabaseAccount, h.config.DatabaseDomain),
		"media":     "//media/",
		"addresses": "//addresses/",
		"_dbs":      "//dbs/",
		"writableLocations": []map[string]interface{}{
			{
				"name":                    "South Central US",
				"databaseAccountEndpoint": h.config.DatabaseEndpoint,
			},
		},
		"readableLocations": []map[string]interface{}{
			{
				"name":                    "South Central US",
				"databaseAccountEndpoint": h.config.DatabaseEndpoint,
			},
		},
		"enableMultipleWriteLocations":   false,
		"continuousBackupEnabled":        false,
		"enableNRegionSynchronousCommit": false,
		"userReplicationPolicy": map[string]interface{}{
			"asyncReplication":  false,
			"minReplicaSetSize": 1,
			"maxReplicasetSize": 4,
		},
		"userConsistencyPolicy":    map[string]interface{}{"defaultConsistencyLevel": "Session"},
		"systemReplicationPolicy":  map[string]interface{}{"minReplicaSetSize": 1, "maxReplicasetSize": 4},
		"readPolicy":               map[string]interface{}{"primaryReadCoefficient": 1, "secondaryReadCoefficient": 1},
		"queryEngineConfiguration": "{\"allowNewKeywords\":true,\"maxJoinsPerSqlQuery\":10,\"maxQueryRequestTimeoutFraction\":0.9,\"maxSqlQueryInputLength\":524288,\"maxUdfRefPerSqlQuery\":10,\"queryMaxInMemorySortDocumentCount\":-1000,\"spatialMaxGeometryPointCount\":256,\"sqlAllowNonFiniteNumbers\":false,\"sqlDisableOptimizationFlags\":0,\"enableSpatialIndexing\":true,\"maxInExpressionItemsCount\":2147483647,\"maxLogicalAndPerSqlQuery\":2147483647,\"maxLogicalOrPerSqlQuery\":2147483647,\"maxSpatialQueryCells\":2147483647,\"sqlAllowAggregateFunctions\":true,\"sqlAllowGroupByClause\":true,\"sqlAllowLike\":true,\"sqlAllowSubQuery\":true,\"sqlAllowScalarSubQuery\":true,\"sqlAllowTop\":true}",
	})
}
