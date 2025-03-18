package types

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

// GetAllTypes retrieves all types from the database.
func (r *Repository) GetAllTypes() ([]Type, error) {
	rows, err := r.db.Query(`SELECT id, name FROM types ORDER BY id ASC`)
	if err != nil {
		log.Printf("Error fetching types: %v", err)
		return nil, fmt.Errorf("could not retrieve types: %v", err)
	}
	defer rows.Close()

	var types []Type
	for rows.Next() {
		var c Type
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		types = append(types, c)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return types, nil
}

// GetTypeByID retrieves a single Type by its ID.
func (r *Repository) GetTypeByID(id int) (*Type, error) {
	var c Type
	err := r.db.QueryRow(`SELECT id, name FROM types WHERE id = $1`, id).
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
func (r *Repository) CreateType(c Type) (Type, error) {
	var TypeID int
	err := r.db.QueryRow(`INSERT INTO types (name) VALUES ($1) RETURNING id`,
		c.Name).Scan(&TypeID)
	if err != nil {
		log.Printf("Error inserting Type: %v", err)
		return Type{}, fmt.Errorf("could not insert Type: %v", err)
	}
	c.ID = TypeID
	return c, nil
}

// UpdateType updates an existing Type.
func (r *Repository) UpdateType(id int, c Type) (Type, error) {
	_, err := r.db.Exec(`UPDATE types SET name = $1 WHERE id = $2`,
		c.Name, id)
	if err != nil {
		log.Printf("Error updating Type: %v", err)
		return Type{}, fmt.Errorf("could not update type: %v", err)
	}
	c.ID = id
	return c, nil
}

// DeleteTransactionType deletes a transaction type by its ID.
func (r *Repository) DeleteType(id int) error {
	_, err := r.db.Exec(`DELETE FROM types WHERE id = $1`, id)
	if err != nil {
		log.Printf("Error deleting Type: %v", err)
		return fmt.Errorf("could not delete Type: %v", err)
	}
	return nil
}
