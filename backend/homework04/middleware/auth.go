package middleware

import (
	"homework04/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWT 中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code:    401,
				Message: "Authorization header missing",
			})
			c.Abort()
			return
		}

		token, err := utils.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code:    401,
				Message: "invalid or expired token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code:    401,
				Message: "invalid token claims",
			})
			c.Abort()
			return
		}

		// 保存用户信息到上下文
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))

		c.Next()
	}
}
