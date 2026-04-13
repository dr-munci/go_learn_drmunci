package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TeacherOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "teacher" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Bu işlem sadece öğretmenler için",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
