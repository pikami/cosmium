package handlers

import (
	"fmt"
	"net/http"
	"strings"

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

type Address struct {
	IsPrimary                     bool   `json:"isPrimary"`
	PhyscialUri                   string `json:"physcialUri"`
	IsAuxiliary                   bool   `json:"isAuxiliary"`
	PartitionTargetReplicaSetSize int    `json:"partitionTargetReplicaSetSize"`
	Protocol                      string `json:"protocol"`
	PartitionKeyRangeId           string `json:"partitionKeyRangeId"`
	PartitionIndex                string `json:"partitionIndex"`
}

func (h *Handlers) GetAddresses(c *gin.Context) {
	addresses := []Address{}

	if h.config.EnableRntbd {
		addresses = append(addresses, Address{
			IsPrimary:                     true,
			PhyscialUri:                   h.config.RntbdEndpoint,
			IsAuxiliary:                   false,
			PartitionTargetReplicaSetSize: 1,
			Protocol:                      "rntbd",
			PartitionKeyRangeId:           "0",
			PartitionIndex:                "0@0",
		})
	}

	if !strings.Contains(c.Request.RequestURI, "protocol%20eq%20rntbd") {
		addresses = append(addresses, Address{
			IsPrimary:                     true,
			PhyscialUri:                   h.config.DatabaseEndpoint,
			IsAuxiliary:                   false,
			PartitionTargetReplicaSetSize: 1,
			Protocol:                      "https",
			PartitionKeyRangeId:           "0",
			PartitionIndex:                "0@0",
		})
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"Addresss": addresses,
		"_count":   len(addresses),
	})
}
