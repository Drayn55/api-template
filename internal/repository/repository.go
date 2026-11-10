// internal/repository/repository.go
package repository

import (
	"context"

	"api-template/internal/audit"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db    *sqlx.DB
	audit *audit.Service
}

func NewRepo(db *sqlx.DB, auditSvc *audit.Service) Repo {
	return Repo{db: db, audit: auditSvc}
}

func (r *Repo) Audit() *audit.Service {
	return r.audit
}

// DB exposes db to embedded structs
func (r *Repo) DB() *sqlx.DB {
	return r.db
}

// BeginTx starts a new transaction
func (r *Repo) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}
