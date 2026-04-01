package utils

import (
	"fmt"
	"homework04/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 生成 JWT
func GenerateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // 24 小时过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	InfoLogger.Println(token)
	re, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		InfoLogger.Println(err)
	}
	return re, err
}

// 验证 JWT
func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})
}
