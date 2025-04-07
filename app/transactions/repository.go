package transactions

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllTransactions() ([]TransactionDTO, error) {
	rows, err := r.db.Query(queryGetAllTransaction)
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		return nil, fmt.Errorf("could not retrieve transactions: %v", err)
	}
	defer rows.Close()

	var transactions []TransactionDTO
	for rows.Next() {
		var transaction TransactionDTO
		var amount []byte
		if err := rows.Scan(
			&transaction.ID, &transaction.ParentID, &transaction.TradeDate, &transaction.Type.ID, &transaction.Category.ID, &transaction.Account.ID,
			&transaction.AccountDest.ID, &amount, &transaction.Description, &transaction.Status, &transaction.Source,
			&transaction.InsertedAt, &transaction.LastUpdatedAt, &transaction.ConfirmedAt,
			&transaction.Category.Name, &transaction.Type.Name, &transaction.Account.Name, &transaction.AccountDest.Name); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		transaction.Amount, _, err = new(big.Float).Parse(string(amount), 10)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return transactions, nil
}

func (r *Repository) GetTransactionByID(id int) (*Transaction, error) {
	var t Transaction
	err := r.db.QueryRow(`SELECT id, parent_id, trade_date, type_id, category_id, account_id, account_dest_id, amount, description, status, source, inserted_at, last_updated_at, confirmed_at FROM transactions WHERE id = $1`, id).
		Scan(&t.ID, &t.ParentID, &t.TradeDate, &t.TypeID, &t.CategoryID, &t.AccountID,
			&t.AccountDestID, &t.Amount, &t.Description, &t.Status, &t.Source, &t.InsertedAt,
			&t.LastUpdatedAt, &t.ConfirmedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction with ID %d not found", id)
		}
		log.Printf("Error fetching transaction by ID: %v", err)
		return nil, fmt.Errorf("could not retrieve transaction: %v", err)
	}
	return &t, nil
}

func (r *Repository) CreateTransaction(tx *sql.Tx, t Transaction) (Transaction, error) {
	log.Printf("%v", t)

	var transactionID int
	// Separate query for auto and manual insertion
	query := queryInsertTransactionAuto
	if t.Source == "Manual" {
		query = queryInsertTransactionManual
	}
	err := tx.QueryRow(query,
		t.ParentID, t.TradeDate, t.TypeID, t.CategoryID, t.AccountID, t.AccountDestID, t.Amount, t.Description, t.Status, t.Source).
		Scan(&transactionID)
	if err != nil {
		log.Printf("Error inserting transaction: %v", err)
		return Transaction{}, fmt.Errorf("could not insert transaction: %v", err)
	}
	t.ID = transactionID
	return t, nil
}

func (r *Repository) UpdateTransaction(id int, t Transaction) (Transaction, error) {
	_, err := r.db.Exec(`UPDATE transactions 
		SET parent_id = $1, trade_date = $2, type_id = $3, category_id = $4, account_id = $5, account_dest_id = $6, amount = $7, description = $8, status = $9, source = $10, last_updated_at = $11, confirmed_at = $12 
		WHERE id = $13`,
		t.ParentID, t.TradeDate, t.TypeID, t.CategoryID, t.AccountID, t.AccountDestID, t.Amount, t.Description, t.Status, t.Source, t.LastUpdatedAt, t.ConfirmedAt, id)
	if err != nil {
		log.Printf("Error updating transaction: %v", err)
		return Transaction{}, fmt.Errorf("could not update transaction: %v", err)
	}
	t.ID = id
	return t, nil
}

func (r *Repository) DeleteTransaction(id int) error {
	_, err := r.db.Exec(`DELETE FROM transactions WHERE id = $1`, id)
	if err != nil {
		log.Printf("Error deleting transaction: %v", err)
		return fmt.Errorf("could not delete transaction: %v", err)
	}
	return nil
}
