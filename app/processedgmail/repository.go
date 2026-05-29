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

func (r *Repository) Exists(messageID string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM processed_gmail_messages
			WHERE gmail_message_id = $1
		)`,
		messageID,
	).Scan(&exists)

	if err != nil {
		logdr.Info("Error checking processed gmail message: %v", err)
		return false, err
	}

	return exists, nil
}

func (r *Repository) RecordSuccess(
	messageID string,
	internalDate int64,
	transactionID int64,
) error {

	_, err := r.db.Exec(
		`INSERT INTO processed_gmail_messages
			(gmail_message_id, internal_date, transaction_id, status)
		VALUES ($1, $2, $3, 'success')
		ON CONFLICT (gmail_message_id) DO NOTHING`,
		messageID,
		internalDate,
		transactionID,
	)

	if err != nil {
		logdr.Info("Error recording gmail success: %v", err)
		return fmt.Errorf("could not record gmail success: %v", err)
	}

	return nil
}

func (r *Repository) RecordFailure(
	messageID string,
	internalDate int64,
	errMsg string,
) error {

	_, err := r.db.Exec(
		`INSERT INTO processed_gmail_messages
			(gmail_message_id, internal_date, status, error_message)
		VALUES ($1, $2, 'failed', $3)
		ON CONFLICT (gmail_message_id)
		DO UPDATE SET
			status = 'failed',
			error_message = EXCLUDED.error_message,
			processed_at = NOW()`,
		messageID,
		internalDate,
		errMsg,
	)

	if err != nil {
		logdr.Info("Error recording gmail failure: %v", err)
		return fmt.Errorf("could not record gmail failure: %v", err)
	}

	return nil
}
