package audit

import (
	"context"
	"encoding/json"

	"api-template/pkg/apperr"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Log(ctx context.Context, tx *sqlx.Tx, entry Entry) error {
	var oldData, newData, changedFields []byte
	var err error

	if entry.OldData != nil {
		oldData, err = json.Marshal(entry.OldData)
		if err != nil {
			return err
		}
	}
	if entry.NewData != nil {
		newData, err = json.Marshal(entry.NewData)
		if err != nil {
			return err
		}
	}
	if len(entry.ChangedFields) > 0 {
		changedFields, err = json.Marshal(entry.ChangedFields)
		if err != nil {
			return err
		}
	}

	q := `INSERT INTO audits (
		table_name, record_id, action, old_data, new_data, changed_fields, created_by
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

	var oldArg, newArg, changedArg any
	if len(oldData) > 0 {
		oldArg = oldData
	}
	if len(newData) > 0 {
		newArg = newData
	}
	if len(changedFields) > 0 {
		changedArg = changedFields
	}

	if tx != nil {
		_, err = tx.ExecContext(ctx, q,
			entry.TableName, entry.RecordID, entry.Action,
			oldArg, newArg, changedArg, entry.CreatedBy,
		)
	} else {
		_, err = r.db.ExecContext(ctx, q,
			entry.TableName, entry.RecordID, entry.Action,
			oldArg, newArg, changedArg, entry.CreatedBy,
		)
	}
	if err != nil {
		return apperr.FromDBError(err)
	}
	return nil
}

