package authentication

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

// https://learn.microsoft.com/en-us/rest/api/cosmos-db/access-control-on-cosmosdb-resources
func GenerateSignature(verb string, resourceType string, resourceId string, date string, masterKey string) string {
	isNameBased := resourceId != "" && ((len(resourceId) > 4 && resourceId[3] == '/') || strings.HasPrefix(strings.ToLower(resourceId), "interopusers"))
	if !isNameBased {
		resourceId = strings.ToLower(resourceId)
	}

	payload := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n",
		strings.ToLower(verb),
		strings.ToLower(resourceType),
		resourceId,
		strings.ToLower(date),
		"")

	masterKeyBytes, _ := base64.StdEncoding.DecodeString(masterKey)
	hash := hmac.New(sha256.New, masterKeyBytes)
	hash.Write([]byte(payload))
	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return signature
}
