package categories

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

// GetAllCategories retrieves all categories from the database.
func (r *Repository) GetAllCategories() ([]Category, error) {
	rows, err := r.db.Query(`SELECT id, name FROM categories ORDER BY id ASC`)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		return nil, fmt.Errorf("could not retrieve categories: %v", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return categories, nil
}

// GetCategoryByID retrieves a single category by its ID.
func (r *Repository) GetCategoryByID(id int) (*Category, error) {
	var c Category
	err := r.db.QueryRow(`SELECT id, name FROM categories WHERE id = $1`, id).
		Scan(&c.ID, &c.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category with ID %d not found", id)
		}
		log.Printf("Error fetching category by ID: %v", err)
		return nil, fmt.Errorf("could not retrieve category: %v", err)
	}
	return &c, nil
}

// CreateCategory inserts a new category into the database.
func (r *Repository) CreateCategory(c Category) (Category, error) {
	var categoryID int
	err := r.db.QueryRow(`INSERT INTO categories (name) VALUES ($1) RETURNING id`,
		c.Name).Scan(&categoryID)
	if err != nil {
		log.Printf("Error inserting category: %v", err)
		return Category{}, fmt.Errorf("could not insert category: %v", err)
	}
	c.ID = categoryID
	return c, nil
}

// UpdateCategory updates an existing category.
func (r *Repository) UpdateCategory(id int, c Category) (Category, error) {
	_, err := r.db.Exec(`UPDATE categories SET name = $1 WHERE id = $2`,
		c.Name, id)
	if err != nil {
		log.Printf("Error updating category: %v", err)
		return Category{}, fmt.Errorf("could not update category: %v", err)
	}
	c.ID = id
	return c, nil
}

// DeleteTransactionType deletes a transaction type by its ID.
func (r *Repository) DeleteCategory(id int) error {
	_, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
		return fmt.Errorf("could not delete category: %v", err)
	}
	return nil
}
