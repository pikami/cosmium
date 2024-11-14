package middleware

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/authentication"
	"github.com/pikami/cosmium/internal/logger"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestUrl := c.Request.URL.String()
		if config.Config.DisableAuth ||
			strings.HasPrefix(requestUrl, config.Config.ExplorerBaseUrlLocation) ||
			strings.HasPrefix(requestUrl, "/cosmium") {
			return
		}

		resourceType := urlToResourceType(requestUrl)
		resourceId := requestToResourceId(c)

		authHeader := c.Request.Header.Get("authorization")
		date := c.Request.Header.Get("x-ms-date")
		expectedSignature := authentication.GenerateSignature(
			c.Request.Method, resourceType, resourceId, date, config.Config.AccountKey)

		decoded, _ := url.QueryUnescape(authHeader)
		params, _ := url.ParseQuery(decoded)
		clientSignature := strings.Replace(params.Get("sig"), " ", "+", -1)
		if clientSignature != expectedSignature {
			logger.Errorf("Got wrong signature from client.\n- Expected: %s\n- Got: %s\n", expectedSignature, clientSignature)
			c.IndentedJSON(401, gin.H{
				"code":    "Unauthorized",
				"message": "Wrong signature.",
			})
			c.Abort()
		}
	}
}

func urlToResourceType(requestUrl string) string {
	var resourceType string
	parts := strings.Split(requestUrl, "/")
	switch len(parts) {
	case 2, 3:
		resourceType = parts[1]
	case 4, 5:
		resourceType = parts[3]
	case 6, 7:
		resourceType = parts[5]
	}

	return resourceType
}

func requestToResourceId(c *gin.Context) string {
	databaseId, _ := c.Params.Get("databaseId")
	collId, _ := c.Params.Get("collId")
	docId, _ := c.Params.Get("docId")
	resourceType := urlToResourceType(c.Request.URL.String())

	var resourceId string
	if databaseId != "" {
		resourceId += "dbs/" + databaseId
	}
	if collId != "" {
		resourceId += "/colls/" + collId
	}
	if docId != "" {
		resourceId += "/docs/" + docId
	}

	isFeed := c.Request.Header.Get("A-Im") == "Incremental Feed"
	if resourceType == "pkranges" && isFeed {
		resourceId = collId
	}

	return resourceId
}
