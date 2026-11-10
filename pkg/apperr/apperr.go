package apperr

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

type AppError struct {
	ErrNo   int
	ErrDesc string
	Cause   error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.ErrDesc
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// New — use when description needs to be dynamic
func New(errNo int, errDesc string) *AppError {
	return &AppError{ErrNo: errNo, ErrDesc: errDesc}
}

// NewCode — use when description is fixed in codes.go
func NewCode(errNo int) *AppError {
	return &AppError{ErrNo: errNo, ErrDesc: Msg(errNo)}
}

// NewWithCause — use when you need to log original error internally
func NewWithCause(errNo int, errDesc string, cause error) *AppError {
	return &AppError{ErrNo: errNo, ErrDesc: errDesc, Cause: cause}
}

// Cast — safely convert error to *AppError, nil if not AppError
func Cast(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// Is — check if error matches a specific ErrNo
func Is(err error, errNo int) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.ErrNo == errNo
	}

	if errNo == ErrDBNoRows && errors.Is(err, sql.ErrNoRows) {
		return true
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return int(mysqlErr.Number) == errNo
	}

	return false
}

