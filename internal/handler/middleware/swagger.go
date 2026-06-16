package middleware

import (
	"github.com/gin-gonic/gin"
)

func SwaggerBasicAuth(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {

		user, pass, ok := c.Request.BasicAuth()
		if !ok {
			c.Header("WWW-Authenticate", `Basic realm="Swagger Restricted"`)
			c.AbortWithStatusJSON(401, gin.H{
				"message": "authentication required",
			})
			return
		}

		if user != username || pass != password {
			c.AbortWithStatusJSON(403, gin.H{
				"message": "invalid credentials",
			})
			return
		}

		c.Next()
	}
}
