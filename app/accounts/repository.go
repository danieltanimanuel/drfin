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

func (r *Repository) GetAllAccounts(ownerID *int) ([]AccountDTO, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.name, a.balance, a.account_group_id, a.status, a.inserted_at, a.closed_at, a.credit_limit,
				 ag.name
		FROM accounts a
		LEFT JOIN account_groups ag ON a.account_group_id = ag.id
		WHERE ($1::int IS NULL OR a.owner_id = $1)
		ORDER BY id ASC`, ownerID)
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

func (r *Repository) GetAccountByName(name string, ownerID int) (*Account, error) {
	var c Account
	err := r.db.QueryRow(`SELECT id, name FROM accounts WHERE name = $1 AND owner_id = $2`, name, ownerID).
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

func (r *Repository) GetAccountByAccountNumber(accountNumber string, ownerID int) (*AccountDTO, error) {
	var account AccountDTO
	err := r.db.QueryRow(`SELECT name, account_number FROM accounts WHERE account_number = $1 AND owner_id = $2`, accountNumber, ownerID).
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

func (r *Repository) GetOwnerID(accountID int) (int, error) {
	var ownerID int
	err := r.db.QueryRow(`SELECT owner_id FROM accounts WHERE id = $1`, accountID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("%w: id %d", ErrAccountNotFound, accountID)
	}
	if err != nil {
		return 0, fmt.Errorf("could not read account owner: %w", err)
	}
	return ownerID, nil
}

func (r *Repository) CreateReceivableAccount(
    tx *sql.Tx,
    name string,
	ownerID int,
) (*Account, error) {

    const loanAccountGroupID = 7

    var account Account

    err := tx.QueryRow(`
        INSERT INTO accounts (
            name,
            balance,
            account_group_id,
            status,
            inserted_at,
			owner_id
        )
        VALUES ($1, 0, $2, 'O', NOW(), $3)
        RETURNING id, name
    `,
        name,
        loanAccountGroupID,
		ownerID,
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

func (r *Repository) AccountGroupExists(tx *sql.Tx, id int) (bool, error) {
	var exists bool
	err := tx.QueryRow(`SELECT EXISTS (SELECT 1 FROM account_groups WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("could not check account group: %w", err)
	}
	return exists, nil
}

func (r *Repository) AccountNameExists(tx *sql.Tx, name string, ownerID int) (bool, error) {
	var exists bool
	err := tx.QueryRow(`SELECT EXISTS (SELECT 1 FROM accounts WHERE LOWER(name) = LOWER($1) AND owner_id = $2)`, name, ownerID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("could not check account name: %w", err)
	}
	return exists, nil
}

func (r *Repository) CreateAccount(tx *sql.Tx, name string, groupID int, balance string, creditLimit *string, ownerID int) (*Account, error) {
	var account Account
	err := tx.QueryRow(`
		INSERT INTO accounts (name, balance, account_group_id, status, inserted_at, credit_limit, owner_id)
		VALUES ($1, $2::numeric, $3, 'O', NOW(), $4::numeric, $5)
		RETURNING id, name`,
		name, balance, groupID, creditLimit, ownerID,
	).Scan(&account.ID, &account.Name)
	if err != nil {
		return nil, fmt.Errorf("could not create account: %w", err)
	}
	return &account, nil
}

// LockAndDiffBalance locks the account row and returns its current balance and
// (newBalance - currentBalance), both computed by Postgres as numeric.
func (r *Repository) LockAndDiffBalance(tx *sql.Tx, accountID int, newBalance string) (oldBalance, delta string, err error) {
	err = tx.QueryRow(`
		SELECT balance::text, ($2::numeric - balance)::text
		FROM accounts
		WHERE id = $1
		FOR UPDATE`, accountID, newBalance).Scan(&oldBalance, &delta)
	if err == sql.ErrNoRows {
		return "", "", fmt.Errorf("%w: id %d", ErrAccountNotFound, accountID)
	}
	if err != nil {
		log.Printf("Error locking account balance: %v", err)
		return "", "", fmt.Errorf("could not read account balance: %w", err)
	}
	return oldBalance, delta, nil
}

func (r *Repository) SetAccountBalance(tx *sql.Tx, accountID int, balance string) error {
	_, err := tx.Exec(`UPDATE accounts SET balance = $1::numeric WHERE id = $2`, balance, accountID)
	if err != nil {
		return fmt.Errorf("could not set account balance: %w", err)
	}
	return nil
}

