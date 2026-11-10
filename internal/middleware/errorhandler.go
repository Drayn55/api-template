package middleware

import (
	"api-template/internal/logger"
	"api-template/pkg/apperr"
	"api-template/pkg/response"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		appErr := apperr.Cast(err)
		if appErr != nil && appErr.Cause != nil {
			logger.L().Error("request error",
				"requestId", requestID(c),
				"errNo", appErr.ErrNo,
				"err", appErr.Cause.Error(),
			)
		}

		response.Error(c, err)
	}
}

