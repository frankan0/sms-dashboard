package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sms-dashboard/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func APITokenMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")

		if token == "" {
			token = c.GetHeader("X-Token")
		}

		if token == "" {
			token = c.PostForm("token")
		}

		contentType := c.GetHeader("Content-Type")
		if token == "" && strings.HasPrefix(contentType, "application/json") {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				var bodyMap map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &bodyMap); err == nil {
					if t, ok := bodyMap["token"].(string); ok {
						token = t
					}
				}
			}
		}

		if token != cfg.APIToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
