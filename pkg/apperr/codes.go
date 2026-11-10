package apperr

// ── error number constants ────────────────────────────────
const (
	// internal
	ErrInternal = 500

	// input errors (handler layer) 400-404
	ErrInvalidJSON   = 400
	ErrMissingField  = 401
	ErrInvalidFormat = 402
	ErrInvalidParam  = 403
	ErrNotFound      = 404

	// auth errors (middleware) 700-703
	ErrTokenMissing = 700
	ErrTokenExpired = 701
	ErrTokenInvalid = 702
	ErrNoPermission = 703

	// generic business rules 600-609
	ErrBusinessRule     = 600
	ErrInsufficientBal  = 601
	ErrAccountInactive  = 602
	ErrLimitExceeded    = 603
	ErrAmountInvalid    = 604
	ErrAlreadyProcessed = 605
	ErrAccountNotFound  = 606
	ErrSameAccount      = 607

	// MariaDB / MySQL error numbers (repository layer)
	ErrDBNoRows          = 1032
	ErrDBDuplicateEntry  = 1062
	ErrDBForeignKeyFail  = 1452
	ErrDBCannotBeNull    = 1048
	ErrDBDataTooLong     = 1406
	ErrDBUnknownColumn   = 1054
	ErrDBTableNotExist   = 1146
	ErrDBDeadlock        = 1213
	ErrDBLockWaitTimeout = 1205
)

// ── error descriptions ────────────────────────────────────
var errMessages = map[int]string{
	// internal
	ErrInternal: "internal server error",

	// input errors
	ErrInvalidJSON:   "invalid request body",
	ErrMissingField:  "required field is missing",
	ErrInvalidFormat: "invalid format",
	ErrInvalidParam:  "invalid parameter",
	ErrNotFound:      "not found",

	// auth errors
	ErrTokenMissing: "authorization token is required",
	ErrTokenExpired: "token has expired",
	ErrTokenInvalid: "token is invalid",
	ErrNoPermission: "you don't have permission",

	// generic business rules
	ErrBusinessRule:     "business rule violation",
	ErrInsufficientBal:  "insufficient balance",
	ErrAccountInactive:  "account is inactive",
	ErrLimitExceeded:    "limit exceeded",
	ErrAmountInvalid:    "amount must be greater than zero",
	ErrAlreadyProcessed: "already processed",
	ErrAccountNotFound:  "account not found",
	ErrSameAccount:      "cannot use same account",

	// DB errors
	ErrDBNoRows:          "data not found",
	ErrDBDuplicateEntry:  "duplicate entry",
	ErrDBForeignKeyFail:  "foreign key constraint failed",
	ErrDBCannotBeNull:    "column cannot be null",
	ErrDBDataTooLong:     "data too long for column",
	ErrDBUnknownColumn:   "unknown column",
	ErrDBTableNotExist:   "table does not exist",
	ErrDBDeadlock:        "database deadlock, please retry",
	ErrDBLockWaitTimeout: "database lock timeout, please retry",
}

// Msg — get description by error number
func Msg(errNo int) string {
	if msg, ok := errMessages[errNo]; ok {
		return msg
	}
	return "unknown error"
}

