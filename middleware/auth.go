package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("your-secret-key-change-in-production")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header gerekli"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz Authorization formatı"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz veya süresi dolmuş token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims okunamadı"})
			c.Abort()
			return
		}

		c.Set("userID", uint(claims["user_id"].(float64)))
		c.Set("userRole", claims["role"].(string))
		c.Next()
	}
}

func TeacherOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkilendirme bilgisi bulunamadı"})
			c.Abort()
			return
		}

		if role != "teacher" && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Bu işlem için öğretmen yetkisi gerekli"})
			c.Abort()
			return
		}

		c.Next()
	}
}
