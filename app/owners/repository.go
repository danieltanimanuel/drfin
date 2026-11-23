package owners

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const selectOwner = `SELECT id, COALESCE(email, ''), COALESCE(name, ''), COALESCE(bank_holder_name, '') FROM users`

func (r *Repository) List() ([]Owner, error) {
	rows, err := r.db.Query(selectOwner + ` ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("could not list owners: %w", err)
	}
	defer rows.Close()

	var list []Owner
	for rows.Next() {
		var o Owner
		if err := rows.Scan(&o.ID, &o.Email, &o.Name, &o.BankHolderName); err != nil {
			return nil, fmt.Errorf("could not scan owner: %w", err)
		}
		list = append(list, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating owners: %w", err)
	}
	return list, nil
}

func (r *Repository) GetByID(id int) (*Owner, error) {
	var o Owner
	err := r.db.QueryRow(selectOwner+` WHERE id = $1`, id).
		Scan(&o.ID, &o.Email, &o.Name, &o.BankHolderName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("owner %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("could not retrieve owner: %w", err)
	}
	return &o, nil
}