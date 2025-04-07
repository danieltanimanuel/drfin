package accounts

import (
	"database/sql"
	"drfin/app/accountgroups"
	"fmt"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllAccounts() ([]AccountDTO, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.name, a.balance, a.account_group_id, a.status, a.inserted_at, a.closed_at,
				 ag.name
		FROM accounts a
		LEFT JOIN account_groups ag ON a.account_group_id = ag.id
		ORDER BY id ASC`)
	if err != nil {
		log.Printf("Error fetching accounts: %v", err)
		return nil, fmt.Errorf("could not retrieve accounts: %v", err)
	}
	defer rows.Close()

	var accounts []AccountDTO
	for rows.Next() {
		var account AccountDTO
		var accountGroup accountgroups.AccountGroup
		if err := rows.Scan(&account.ID, &account.Name, &account.Balance, &accountGroup.ID, &account.Status, &account.InsertedAt, &account.ClosedAt, &accountGroup.Name); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		account.AccountGroup = accountGroup
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return accounts, nil
}

func (r *Repository) AddAccountBalance(tx *sql.Tx, accountID int, amount float64) error {
	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING balance`

	return r.updateAccountBalance(tx, query, accountID, amount)
}

func (r *Repository) SubstractAccountBalance(tx *sql.Tx, accountID int, amount float64) error {
	query := `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
		RETURNING balance`

	return r.updateAccountBalance(tx, query, accountID, amount)
}

func (r *Repository) updateAccountBalance(tx *sql.Tx, query string, accountID int, amount float64) error {
	var newBalance float64
	err := tx.QueryRow(query, amount, accountID).Scan(&newBalance)
	if err != nil {
		log.Printf("failed to update account balance: %v", err)
		return fmt.Errorf("failed to update account balance: %v", err)
	}
	return nil
}

// func (r *Repository) GetAccountGroupByID(id int) (*Account, error) {
// 	var c Account
// 	err := r.db.QueryRow(`SELECT id, name FROM accounts WHERE id = $1`, id).
// 		Scan(&c.ID, &c.Name)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("accounts with ID %d not found", id)
// 		}
// 		log.Printf("Error fetching Account by ID: %v", err)
// 		return nil, fmt.Errorf("could not retrieve Account: %v", err)
// 	}
// 	return &c, nil
// }

// func (r *Repository) CreateAccountGroup(c AccountGroup) (AccountGroup, error) {
// 	var AccountGroupID int
// 	err := r.db.QueryRow(`INSERT INTO account_groups (name) VALUES ($1) RETURNING id`,
// 		c.Name).Scan(&AccountGroupID)
// 	if err != nil {
// 		log.Printf("Error inserting AccountGroup: %v", err)
// 		return AccountGroup{}, fmt.Errorf("could not insert AccountGroup: %v", err)
// 	}
// 	c.ID = AccountGroupID
// 	return c, nil
// }

// func (r *Repository) UpdateAccountGroup(id int, c AccountGroup) (AccountGroup, error) {
// 	_, err := r.db.Exec(`UPDATE account_groups SET name = $1 WHERE id = $2`,
// 		c.Name, id)
// 	if err != nil {
// 		log.Printf("Error updating AccountGroup: %v", err)
// 		return AccountGroup{}, fmt.Errorf("could not update AccountGroup: %v", err)
// 	}
// 	c.ID = id
// 	return c, nil
// }

// func (r *Repository) DeleteAccountGroup(id int) error {
// 	_, err := r.db.Exec(`DELETE FROM account_groups WHERE id = $1`, id)
// 	if err != nil {
// 		log.Printf("Error deleting AccountGroup: %v", err)
// 		return fmt.Errorf("could not delete AccountGroup: %v", err)
// 	}
// 	return nil
// }
