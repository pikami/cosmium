package authentication_test

import (
	"testing"

	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/authentication"
	"github.com/stretchr/testify/assert"
)

const (
	testDate = "Fri, 17 Dec 1926 03:15:00 GMT"
)

func Test_GenerateSignature(t *testing.T) {
	t.Run("Should generate GET signature", func(t *testing.T) {
		signature := authentication.GenerateSignature("GET", "colls", "dbs/Test Database/colls/Test Collection", testDate, config.DefaultAccountKey)
		assert.Equal(t, "cugjaA51bjCvxVi8LXg3XB+ZVKaFAZshILoJZF9nfEY=", signature)
	})

	t.Run("Should generate POST signature", func(t *testing.T) {
		signature := authentication.GenerateSignature("POST", "colls", "dbs/Test Database", testDate, config.DefaultAccountKey)
		assert.Equal(t, "E92FgDG9JiNX+NfsI+edOFtgkZRDkrrJxIfl12Vsu8A=", signature)
	})

	t.Run("Should generate DELETE signature", func(t *testing.T) {
		signature := authentication.GenerateSignature("DELETE", "dbs", "dbs/Test Database", testDate, config.DefaultAccountKey)
		assert.Equal(t, "LcuXXg0TcXxZG0kUCj9tZIWRy2yCzim3oiqGiHpRqGs=", signature)
	})

	t.Run("Should generate PKRANGES signature", func(t *testing.T) {
		signature := authentication.GenerateSignature("GET", "pkranges", "m4d+xG08uVM=", testDate, config.DefaultAccountKey)
		assert.Equal(t, "6S5ceZsl2EXWB3Jo5bJcK7zv8NxXnsxWPWD9TH3nNMo=", signature)
	})

	t.Run("Should generate PATCH signature", func(t *testing.T) {
		signature := authentication.GenerateSignature("PATCH", "docs", "dbs/test-db/colls/test-coll/docs/67890", testDate, config.DefaultAccountKey)
		assert.Equal(t, "VR1ddfxKBXnoaT+b3WkhyYVc9JmGNpTnaRmyDM44398=", signature)
	})
}
