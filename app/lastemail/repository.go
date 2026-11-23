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

func (r *Repository) GetLastEmail(ownerID int) (string, error) {
	var lastEmail string
	err := r.db.QueryRow(`select message_id from last_email where owner_id = $1`, ownerID).
		Scan(&lastEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("last email not found")
		}
		logdr.Error("Error fetching last email : %v", err)
		return "", fmt.Errorf("could not retrieve last email: %v", err)
	}
	return lastEmail, nil
}

func (r *Repository) UpdateLastEmail(messageID string, ownerID int) (string, error) {
	_, err := r.db.Exec(`UPDATE last_email 
		SET message_id = $1
		WHERE owner_id = $2`,
		messageID, ownerID)
	if err != nil {
		logdr.Error("Error updating transaction: %v", err)
		return messageID, fmt.Errorf("could not update transaction: %v", err)
	}
	return messageID, nil
}

func (r *Repository) GetLastInternalDate(ownerID int) (int64, error) {
	var internalDate sql.NullInt64
	err := r.db.QueryRow(`select internal_date from last_email where owner_id = $1`, ownerID).Scan(&internalDate)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("no last_email row for owner %d (insert one first): %w", ownerID, err)
	}
	if err != nil {
		return 0, err
	}
	if internalDate.Valid {
		return internalDate.Int64, nil
	}
	return 0, nil // fallback for old data
}

func (r *Repository) UpdateLastInternalDate(ownerID int, internalDate int64) error {
	logdr.Info("internaldate", "owner_id", ownerID, "internaldate", internalDate)
	res, err := r.db.Exec(`UPDATE last_email SET internal_date = $1 WHERE owner_id = $2`, internalDate, ownerID)
	if err != nil {
		logdr.Error("Error updating last_email: %v", err)
		return fmt.Errorf("could not update last_email: %v", err)
	}
	if n, err := res.RowsAffected(); err != nil || n != 1 {
		return fmt.Errorf("expected 1 last_email row for owner %d, updated %d (err: %v)", ownerID, n, err)
	}
	return nil
}

