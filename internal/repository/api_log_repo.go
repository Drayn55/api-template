package repository

import (
	"context"

	"api-template/internal/model"

	"github.com/jmoiron/sqlx"
)

type APILogRepo struct {
	Repo
}

func NewAPILogRepo(db *sqlx.DB) *APILogRepo {
	if db == nil {
		return nil
	}
	return &APILogRepo{Repo: NewRepo(db, nil)}
}

func (r *APILogRepo) Create(data model.APILog) error {
	if r == nil || r.DB() == nil {
		return nil
	}

	_, err := r.DB().ExecContext(context.Background(), `
		INSERT INTO api_logs (
			hash_id, method, endpoint, request_body, response_body,
			status_code, ip_address, user_agent,
			employee_id, execution_time_ms
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		data.HashID,
		data.Method,
		data.Endpoint,
		data.RequestBody,
		data.ResponseBody,
		data.StatusCode,
		data.IPAddress,
		data.UserAgent,
		data.EmployeeID,
		data.ExecutionTimeMs,
	)
	if err != nil {
		return err
	}
	return nil
}
