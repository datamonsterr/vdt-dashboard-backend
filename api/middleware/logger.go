package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger returns a Gin middleware for logging requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			log := logrus.WithFields(logrus.Fields{
				"status":     param.StatusCode,
				"method":     param.Method,
				"path":       param.Path,
				"ip":         param.ClientIP,
				"latency":    param.Latency,
				"user_agent": param.Request.UserAgent(),
			})

			if param.StatusCode >= 400 {
				log.Error("HTTP Request")
			} else {
				log.Info("HTTP Request")
			}

			return ""
		},
		Output:    nil, // We're using logrus, so we don't need gin's output
		SkipPaths: []string{"/health"},
	})
}

// Recovery returns a Gin middleware for panic recovery
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logrus.WithFields(logrus.Fields{
			"panic": recovered,
			"path":  c.Request.URL.Path,
		}).Error("Panic recovered")

		c.JSON(500, gin.H{
			"success": false,
			"message": "Internal server error",
			"error": map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"details": "An unexpected error occurred",
			},
		})
	})
}
