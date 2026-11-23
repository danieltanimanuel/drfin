package processedgmail

import (
	"database/sql"
	"drfin/app/shared/logdr"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Exists(ownerID int, messageID string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM processed_gmail_messages
			WHERE owner_id = $1 AND gmail_message_id = $2 AND status = 'success''
		)`,
		ownerID, messageID,
	).Scan(&exists)

	if err != nil {
		logdr.Error("Error checking processed gmail message: %v", err)
		return false, err
	}

	return exists, nil
}

func (r *Repository) RecordSuccess(ownerID int, messageID string, internalDate int64, transactionID int64) error {

	_, err := r.db.Exec(
		`INSERT INTO processed_gmail_messages
			(owner_id, gmail_message_id, internal_date, transaction_id, status)
		VALUES ($1, $2, $3, $4, 'success')
		ON CONFLICT (gmail_message_id) DO UPDATE SET
		status = 'success', error_message = NULL,
		transaction_id = EXCLUDED.transaction_id, processed_at = NOW()`,
		ownerID,
		messageID,
		internalDate,
		transactionID,
	)

	if err != nil {
		logdr.Error("Error recording gmail success: %v", err)
		return fmt.Errorf("could not record gmail success: %v", err)
	}

	return nil
}

func (r *Repository) RecordFailure(ownerID int, messageID string, internalDate int64, errMsg string) error {

	_, err := r.db.Exec(
		`INSERT INTO processed_gmail_messages
			(owner_id, gmail_message_id, internal_date, status, error_message)
		VALUES ($1, $2, $3, 'failed', $4)
		ON CONFLICT (owner_id, gmail_message_id)
		DO UPDATE SET
			status = 'failed',
			error_message = EXCLUDED.error_message,
			processed_at = NOW()`,
		ownerID,
		messageID,
		internalDate,
		errMsg,
	)

	if err != nil {
		logdr.Error("Error recording gmail failure: %v", err)
		return fmt.Errorf("could not record gmail failure: %v", err)
	}

	return nil
}
