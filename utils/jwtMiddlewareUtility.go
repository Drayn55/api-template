package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JWTSecret(), nil

		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user", claims)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

func GetEmployeeIDFromToken(c *gin.Context) *string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret(), nil

	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		for _, key := range []string{"employee_id", "emp_id", "id", "sub", "user_id", "uid"} {
			if val, exists := claims[key]; exists {
				strVal := fmt.Sprintf("%v", val)
				return &strVal
			}
		}
	}

	return nil
}
