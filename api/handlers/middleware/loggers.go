package middleware

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/logger"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := io.ReadAll(c.Request.Body)
		rdr1 := io.NopCloser(bytes.NewBuffer(buf))
		rdr2 := io.NopCloser(bytes.NewBuffer(buf))

		bodyStr := readBody(rdr1)
		if bodyStr != "" {
			logger.Debug(bodyStr)
		}

		c.Request.Body = rdr2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
