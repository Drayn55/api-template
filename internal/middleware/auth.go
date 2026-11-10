// internal/middleware/auth.go
package middleware

import (
	"context"

	"api-template/pkg/apperr"
	"api-template/pkg/jwtutil"

	"github.com/gin-gonic/gin"
)

func Auth(jwt *jwtutil.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Error(apperr.NewCode(apperr.ErrTokenMissing))
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(token)
		if err != nil {
			switch err.Error() {
			case "token expired":
				c.Error(apperr.NewCode(apperr.ErrTokenExpired))
			default:
				c.Error(apperr.NewCode(apperr.ErrTokenInvalid))
			}
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), "user_id", claims.EmployeeID)
		ctx = context.WithValue(ctx, "authorization", token)
		ctx = context.WithValue(ctx, "audit_user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_name", claims.Name)

		c.Request = c.Request.WithContext(ctx)

		c.Set("employeeId", claims.EmployeeID)
		c.Set("userName", claims.Name)
		c.Set("userId", claims.UserID)

		c.Next()
	}
}
