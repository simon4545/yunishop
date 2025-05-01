package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("your-secret-key") // 替换成你的密钥

// 生成 JWT Token
func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24小时后过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 解析 JWT Token
func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 从 Header 获取 Token
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Authorization header missing",
			})
		}

		// 检查格式是否为 "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid token format",
			})
		}

		tokenString := tokenParts[1]

		// 解析 Token
		token, err := ParseToken(tokenString)
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid or expired token",
			})
		}

		// 将 Token Claims 存入 Context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid token claims",
			})
		}

		c.Set("user_id", claims["user_id"]) // 存储用户 ID
		return next(c)
	}
}
