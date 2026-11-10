// pkg/response/response.go
package response

import (
	"net/http"

	"api-template/pkg/apperr"

	"github.com/gin-gonic/gin"
)

// Response — single shape for ALL API responses.
type Response struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    any    `json:"meta,omitempty"`
	Summary any    `json:"summary,omitempty"`
}

// Success — operation succeeded; data is the payload (may be nil).
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMeta — success with pagination or other metadata.
func SuccessWithMeta(c *gin.Context, data, meta any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    0,
		Message: "success",
		Data:    data,
		Meta:    meta,
	})
}

// SuccessWithSummary — success with meta and summary (e.g. totals).
func SuccessWithSummary(c *gin.Context, data, meta, summary any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    0,
		Message: "success",
		Data:    data,
		Meta:    meta,
		Summary: summary,
	})
}

// Error — used by errorhandler middleware; formats AppError into Response.
func Error(c *gin.Context, err error) {
	appErr := apperr.Cast(err)
	if appErr == nil {
		if apperr.Is(err, apperr.ErrDBNoRows) {
			appErr = apperr.Cast(apperr.FromDBError(err))
		} else if wrapped := apperr.FromDBError(err); wrapped != nil {
			if casted := apperr.Cast(wrapped); casted != nil && casted.ErrNo != apperr.ErrInternal {
				appErr = casted
			}
		}
	}
	if appErr != nil {
		c.JSON(httpStatusFromErrNo(appErr.ErrNo), Response{
			Success: false,
			Code:    appErr.ErrNo,
			Message: appErr.ErrDesc,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Code:    apperr.ErrInternal,
		Message: "internal server error",
		Data:    nil,
	})
}

// httpStatusFromErrNo maps internal error codes to HTTP status codes.
func httpStatusFromErrNo(errNo int) int {
	switch {
	// input validation errors → 400 Bad Request
	case errNo == apperr.ErrInvalidJSON,
		errNo == apperr.ErrMissingField,
		errNo == apperr.ErrInvalidFormat,
		errNo == apperr.ErrInvalidParam,
		errNo == apperr.ErrBusinessRule,
		errNo == apperr.ErrInsufficientBal,
		errNo == apperr.ErrLimitExceeded,
		errNo == apperr.ErrAmountInvalid,
		errNo == apperr.ErrAlreadyProcessed,
		errNo == apperr.ErrSameAccount,
		errNo == apperr.ErrAccountInactive:
		return http.StatusBadRequest

	// auth errors → 401 Unauthorized
	case errNo == apperr.ErrTokenMissing,
		errNo == apperr.ErrTokenExpired,
		errNo == apperr.ErrTokenInvalid:
		return http.StatusUnauthorized

	// permission errors → 403 Forbidden
	case errNo == apperr.ErrNoPermission:
		return http.StatusForbidden

	// not found errors → 404 Not Found
	case errNo == apperr.ErrNotFound,
		errNo == apperr.ErrDBNoRows,
		errNo == apperr.ErrAccountNotFound:
		return http.StatusNotFound

	// conflict errors → 409 Conflict
	case errNo == apperr.ErrDBDuplicateEntry:
		return http.StatusConflict

	// all other errors → 500 Internal Server Error
	default:
		return http.StatusInternalServerError
	}
}
