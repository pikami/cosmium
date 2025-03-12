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
	ts := runTestServer()
	documents_InitializeDb(t, ts)
	defer ts.Server.Close()

	t.Run("Read doc with client that appends slash to path", func(t *testing.T) {
		resourceIdTemplate := "dbs/%s/colls/%s/docs/%s"
		path := fmt.Sprintf(resourceIdTemplate, testDatabaseName, testCollectionName, "12345")
		testUrl := ts.URL + "/" + path + "/"
		date := time.Now().Format(time.RFC1123)
		signature := authentication.GenerateSignature("GET", "docs", path, date, config.DefaultAccountKey)
		httpClient := &http.Client{}
		req, _ := http.NewRequest("GET", testUrl, nil)
		req.Header.Add("x-ms-date", date)
		req.Header.Add("authorization", "sig="+url.QueryEscape(signature))
		res, err := httpClient.Do(req)

		assert.Nil(t, err)

		if res != nil {
			defer res.Body.Close()
			assert.Equal(t, http.StatusOK, res.StatusCode, "Expected HTTP status 200 OK")
		} else {
			t.FailNow()
		}
	})
}
