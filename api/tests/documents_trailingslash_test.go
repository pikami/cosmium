package tests_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/authentication"
	"github.com/stretchr/testify/assert"
)

// Request document with trailing slash like python cosmosdb client does.
func Test_Documents_Read_Trailing_Slash(t *testing.T) {
	ts, _ := documents_InitializeDb(t)
	defer ts.Close()

	t.Run("Read doc with client that appends slash to path", func(t *testing.T) {
		resourceIdTemplate := "dbs/%s/colls/%s/docs/%s"
		path := fmt.Sprintf(resourceIdTemplate, testDatabaseName, testCollectionName, "12345")
		testUrl := ts.URL + "/" + path + "/"
		date := time.Now().Format(time.RFC1123)
		signature := authentication.GenerateSignature("GET", "docs", path, date, config.Config.AccountKey)
		httpClient := &http.Client{}
		req, _ := http.NewRequest("GET", testUrl, nil)
		req.Header.Add("x-ms-date", date)
		req.Header.Add("authorization", "sig="+url.QueryEscape(signature))
		_, err := httpClient.Do(req)

		assert.Nil(t, err)

	})

}
