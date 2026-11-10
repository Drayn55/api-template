package apperr

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

// FromDBError — optional helper to normalize DB errors into AppError.
func FromDBError(err error) error {
	if err == nil {
		return nil
	}

	if appErr := Cast(err); appErr != nil {
		return appErr
	}

	// check for no rows first — standard Go error, not MySQL
	if errors.Is(err, sql.ErrNoRows) {
		return &AppError{
			ErrNo:   ErrDBNoRows,
			ErrDesc: "no rows found",
			Cause:   err,
		}
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return &AppError{
			ErrNo:   int(mysqlErr.Number),
			ErrDesc: mysqlErr.Message,
			Cause:   err,
		}
	}

	// other DB error
	return &AppError{
		ErrNo:   ErrInternal,
		ErrDesc: "database error",
		Cause:   err,
	}
}

// IsDBError — check specific MariaDB/MySQL error number
func IsDBError(err error, errNo int) bool {
	return Is(err, errNo)
}
