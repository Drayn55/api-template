package audit

import (
	"context"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func UserIDFromContext(ctx context.Context) *int64 {
	if v, ok := ctx.Value("audit_user_id").(int64); ok && v > 0 {
		return &v
	}
	if s, ok := ctx.Value("user_id").(string); ok && s != "" {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			return &id
		}
	}
	return nil
}

func (s *Service) LogInsert(
	ctx context.Context,
	tx *sqlx.Tx,
	table string,
	recordID int64,
	newData any,
	userID *int64,
) error {
	if s.repo == nil {
		return nil
	}
	return s.repo.Log(ctx, tx, Entry{
		TableName: table,
		RecordID:  recordID,
		Action:    ActionInsert,
		NewData:   newData,
		CreatedBy: userID,
	})
}

func (s *Service) LogUpdate(
	ctx context.Context,
	tx *sqlx.Tx,
	table string,
	recordID int64,
	oldData any,
	newData any,
	userID *int64,
) error {
	if s.repo == nil {
		return nil
	}
	return s.repo.Log(ctx, tx, Entry{
		TableName:     table,
		RecordID:      recordID,
		Action:        ActionUpdate,
		OldData:       oldData,
		NewData:       newData,
		ChangedFields: ChangedFields(oldData, newData),
		CreatedBy:     userID,
	})
}

func (s *Service) LogDelete(
	ctx context.Context,
	tx *sqlx.Tx,
	table string,
	recordID int64,
	oldData any,
	userID *int64,
) error {
	if s.repo == nil {
		return nil
	}
	return s.repo.Log(ctx, tx, Entry{
		TableName: table,
		RecordID:  recordID,
		Action:    ActionDelete,
		OldData:   oldData,
		CreatedBy: userID,
	})
}
