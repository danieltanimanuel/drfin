package accountGroups

import (
	"database/sql"
	"fmt"
	"log"
)

// Repository handles the database queries for transaction types.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new repository instance.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetAllAccountGroups retrieves all accountgroups from the database.
func (r *Repository) GetAllAccountGroups() ([]AccountGroup, error) {
	rows, err := r.db.Query(`SELECT id, name FROM account_groups ORDER BY id ASC`)
	if err != nil {
		log.Printf("Error fetching account groups: %v", err)
		return nil, fmt.Errorf("could not retrieve account groups: %v", err)
	}
	defer rows.Close()

	var accountGroups []AccountGroup
	for rows.Next() {
		var c AccountGroup
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		accountGroups = append(accountGroups, c)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return accountGroups, nil
}

// GetAccountGroupByID retrieves a single Type by its ID.
func (r *Repository) GetAccountGroupByID(id int) (*AccountGroup, error) {
	var c AccountGroup
	err := r.db.QueryRow(`SELECT id, name FROM account_groups WHERE id = $1`, id).
		Scan(&c.ID, &c.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Type with ID %d not found", id)
		}
		log.Printf("Error fetching Type by ID: %v", err)
		return nil, fmt.Errorf("could not retrieve Type: %v", err)
	}
	return &c, nil
}

// CreateAccountGroup inserts a new AccountGroup into the database.
func (r *Repository) CreateAccountGroup(c AccountGroup) (AccountGroup, error) {
	var AccountGroupID int
	err := r.db.QueryRow(`INSERT INTO account_groups (name) VALUES ($1) RETURNING id`,
		c.Name).Scan(&AccountGroupID)
	if err != nil {
		log.Printf("Error inserting AccountGroup: %v", err)
		return AccountGroup{}, fmt.Errorf("could not insert AccountGroup: %v", err)
	}
	c.ID = AccountGroupID
	return c, nil
}

// UpdateAccountGroup updates an existing Type.
func (r *Repository) UpdateAccountGroup(id int, c AccountGroup) (AccountGroup, error) {
	_, err := r.db.Exec(`UPDATE account_groups SET name = $1 WHERE id = $2`,
		c.Name, id)
	if err != nil {
		log.Printf("Error updating Type: %v", err)
		return AccountGroup{}, fmt.Errorf("could not update AccountGroup: %v", err)
	}
	c.ID = id
	return c, nil
}

// DeleteTransactionType deletes a transaction type by its ID.
func (r *Repository) DeleteAccountGroup(id int) error {
	_, err := r.db.Exec(`DELETE FROM account_groups WHERE id = $1`, id)
	if err != nil {
		log.Printf("Error deleting Type: %v", err)
		return fmt.Errorf("could not delete Type: %v", err)
	}
	return nil
}
