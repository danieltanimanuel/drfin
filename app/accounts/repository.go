package accounts

import (
	"database/sql"
	"drfin/app/accountgroups"
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

func parseBigFloat(value interface{}) (*big.Float, error) {
	var str string

	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return nil, nil // or return an error if you prefer strict handling
	}

	f, _, err := big.ParseFloat(str, 10, 256, big.ToNearestEven)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *Repository) GetAllAccounts() ([]AccountDTO, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.name, a.balance, a.account_group_id, a.status, a.inserted_at, a.closed_at, a.credit_limit,
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
		var balance, creditLimit []byte
		if err := rows.Scan(&account.ID, &account.Name, &balance, &accountGroup.ID, &account.Status, &account.InsertedAt, &account.ClosedAt, &creditLimit, &accountGroup.Name); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		account.AccountGroup = accountGroup
		account.Balance, _, err = new(big.Float).Parse(string(balance), 10)
		if err != nil {
			return nil, fmt.Errorf("could not parse balance: %v", err)
		}

		if creditLimit != nil {
			account.CreditLimit, _, err = new(big.Float).Parse(string(creditLimit), 10)
			if err != nil {
				return nil, fmt.Errorf("could not parse credit limit: %v", err)
			}
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return accounts, nil
}

func (r *Repository) AddAccountBalance(tx *sql.Tx, accountID int, amount *big.Float) error {
	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING balance`

	return r.updateAccountBalance(tx, query, accountID, amount)
}

func (r *Repository) SubstractAccountBalance(tx *sql.Tx, accountID int, amount *big.Float) error {
	query := `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
		RETURNING balance`

	return r.updateAccountBalance(tx, query, accountID, amount)
}

func (r *Repository) updateAccountBalance(tx *sql.Tx, query string, accountID int, amount *big.Float) error {
	var newBalance float64
	err := tx.QueryRow(query, amount.Text('f', 2), accountID).Scan(&newBalance)
	if err != nil {
		log.Printf("failed to update account balance: %v", err)
		return fmt.Errorf("failed to update account balance: %v", err)
	}
	return nil
}

func (r *Repository) GetAccountByName(name string) (*Account, error) {
	var c Account
	err := r.db.QueryRow(`SELECT id, name FROM accounts WHERE name = $1`, name).
		Scan(&c.ID, &c.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("accounts with name %s not found", name)
		}
		log.Printf("Error fetching Account by name: %v", err)
		return nil, fmt.Errorf("could not retrieve Account: %v", err)
	}
	return &c, nil
}

func (r *Repository) GetAccountByAccountNumber(accountNumber string) (*AccountDTO, error) {
	var account AccountDTO
	err := r.db.QueryRow(`SELECT name, account_number FROM accounts WHERE account_number = $1`, accountNumber).
		Scan(&account.Name, &account.AccountNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error fetching Account by account_number: %v", err)
		return nil, fmt.Errorf("could not retrieve account_number: %v", err)
	}
	return &account, nil
}
func (r *Repository) CreateReceivableAccount(
    tx *sql.Tx,
    name string,
) (*Account, error) {

    const loanAccountGroupID = 7

    var account Account

    err := tx.QueryRow(`
        INSERT INTO accounts (
            name,
            balance,
            account_group_id,
            status,
            inserted_at
        )
        VALUES ($1, 0, $2, 'O', NOW())
        RETURNING id, name
    `,
        name,
        loanAccountGroupID,
    ).Scan(
        &account.ID,
        &account.Name,
    )

    if err != nil {
        return nil, fmt.Errorf(
            "could not create receivable account: %w",
            err,
        )
    }

    return &account, nil
}