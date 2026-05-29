package lastemail

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

func (r *Repository) GetLastEmail() (string, error) {
	var lastEmail string
	err := r.db.QueryRow(`select message_id from last_email`).
		Scan(&lastEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("last email not found")
		}
		logdr.Info("Error fetching last email : %v", err)
		return "", fmt.Errorf("could not retrieve last email: %v", err)
	}
	return lastEmail, nil
}

func (r *Repository) UpdateLastEmail(messageID string) (string, error) {
	_, err := r.db.Exec(`UPDATE last_email 
		SET message_id = $1`,
		messageID)
	if err != nil {
		logdr.Info("Error updating transaction: %v", err)
		return messageID, fmt.Errorf("could not update transaction: %v", err)
	}
	return messageID, nil
}

func (r *Repository) GetLastInternalDate() (int64, error) {
	var internalDate sql.NullInt64

	err := r.db.QueryRow(`select internal_date from last_email`).Scan(&internalDate)

	if err != nil {
		return 0, err
	}

	if internalDate.Valid {
		return internalDate.Int64, nil
	}

	// fallback for old data
	return 0, nil
}

func (r *Repository) UpdateLastInternalDate(internalDate int64) error {
	logdr.Info("internaldate", internalDate)
	_, err := r.db.Exec(
		`UPDATE last_email SET internal_date = $1`,
		internalDate,
	)
	if err != nil {
		logdr.Info("Error updating transaction: %v", err)
		return fmt.Errorf("could not update transaction: %v", err)
	}
	return nil
}
