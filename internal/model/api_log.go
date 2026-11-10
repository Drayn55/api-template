package model

import "time"

type APILog struct {
	ID              int64      `db:"id"`
	HashID          string     `db:"hash_id"`
	Method          *string    `db:"method"`
	Endpoint        *string    `db:"endpoint"`
	RequestBody     *string    `db:"request_body"`
	ResponseBody    *string    `db:"response_body"`
	StatusCode      *int       `db:"status_code"`
	IPAddress       *string    `db:"ip_address"`
	UserAgent       *string    `db:"user_agent"`
	EmployeeID      *string    `db:"employee_id"`
	ExecutionTimeMs *float64   `db:"execution_time_ms"`
	CreatedAt       *time.Time `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}
