package transactiontypes

import (
	"database/sql"
	"fmt"
	"log"
)

// Repository handles the database queries for transaction transactionTypes.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new repository instance.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetAllTypes retrieves all transactionTypes from the database.
func (r *Repository) GetAllTransactionTypes() ([]TransactionType, error) {
	rows, err := r.db.Query(`SELECT id, name FROM transaction_types ORDER BY id ASC`)
	if err != nil {
		log.Printf("Error fetching transactionTypes: %v", err)
		return nil, fmt.Errorf("could not retrieve transactionTypes: %v", err)
	}
	defer rows.Close()

	var transactionType []TransactionType
	for rows.Next() {
		var c TransactionType
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		transactionType = append(transactionType, c)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return transactionType, nil
}

// GetTypeByID retrieves a single Type by its ID.
func (r *Repository) GetTransactionTypeByID(id int) (*TransactionType, error) {
	var c TransactionType
	err := r.db.QueryRow(`SELECT id, name FROM transaction_types WHERE id = $1`, id).
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

// CreateType inserts a new Type into the database.
func (r *Repository) CreateTransactionType(c TransactionType) (TransactionType, error) {
	var TransactionTypeID int
	err := r.db.QueryRow(`INSERT INTO transaction_types (name) VALUES ($1) RETURNING id`,
		c.Name).Scan(&TransactionTypeID)
	if err != nil {
		log.Printf("Error inserting Type: %v", err)
		return TransactionType{}, fmt.Errorf("could not insert Type: %v", err)
	}
	c.ID = TransactionTypeID
	return c, nil
}

// UpdateType updates an existing Type.
func (r *Repository) UpdateTransactionType(id int, c TransactionType) (TransactionType, error) {
	_, err := r.db.Exec(`UPDATE transaction_types SET name = $1 WHERE id = $2`,
		c.Name, id)
	if err != nil {
		log.Printf("Error updating Type: %v", err)
		return TransactionType{}, fmt.Errorf("could not update type: %v", err)
	}
	c.ID = id
	return c, nil
}

// DeleteTransactionType deletes a transaction type by its ID.
func (r *Repository) DeleteTransactionType(id int) error {
	_, err := r.db.Exec(`DELETE FROM transaction_types WHERE id = $1`, id)
	if err != nil {
		log.Printf("Error deleting Type: %v", err)
		return fmt.Errorf("could not delete Type: %v", err)
	}
	return nil
}
