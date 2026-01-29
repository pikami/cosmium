package continuationtoken

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Generate(t *testing.T) {
	token := Generate("test-resource-id", 1, 100)

	assert.Equal(t, "test-resource-id", token.Token.ResourceId)
	assert.Equal(t, 1, token.Token.PageIndex)
	assert.Equal(t, 100, token.Token.TotalResults)
}

func Test_FromString(t *testing.T) {
	token := FromString("{\"token\":\"-RID:~test-resource-id#RT:1#TRC:100#ISV:2#IEO:65567#QCF:8#LR:1\",\"range\":{\"min\":\"\",\"max\":\"FF\"}}")

	assert.Equal(t, "test-resource-id", token.Token.ResourceId)
	assert.Equal(t, 1, token.Token.PageIndex)
	assert.Equal(t, 100, token.Token.TotalResults)
}

func Test_ToString(t *testing.T) {
	token := Generate("test-resource-id", 1, 100)
	assert.Equal(t, "{\"token\":\"-RID:~test-resource-id#RT:1#TRC:100#ISV:2#IEO:65567#QCF:8#LR:1\",\"range\":{\"min\":\"\",\"max\":\"FF\"}}", token.ToString())
}

func Test_GenerateDefault(t *testing.T) {
	token := GenerateDefault("test-resource-id")
	assert.Equal(t, "test-resource-id", token.Token.ResourceId)
	assert.Equal(t, 0, token.Token.PageIndex)
	assert.Equal(t, 0, token.Token.TotalResults)
}
