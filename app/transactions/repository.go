package transactions

import (
	"context" // ADD
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

//
// READ
//

func (r *Repository) GetAllTransactions(
	ctx context.Context, // ADD
) ([]TransactionDTO, error) {

	// Enhancement:
	// QueryContext respects request timeout & cancellation
	rows, err := r.db.QueryContext(ctx, queryGetAllTransaction)
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		return nil, fmt.Errorf("could not retrieve transactions: %w", err)
	}
	defer rows.Close()

	var transactions []TransactionDTO
	for rows.Next() {
		var transaction TransactionDTO
		var amount []byte
		var accountDestID sql.NullInt64
		var accountDestName sql.NullString

		if err := rows.Scan(
			&transaction.ID,
			&transaction.ParentID,
			&transaction.TradeDate,
			&transaction.Type.ID,
			&transaction.Category.ID,
			&transaction.Account.ID,
			&accountDestID,
			&amount,
			&transaction.Description,
			&transaction.Status,
			&transaction.Source,
			&transaction.InsertedAt,
			&transaction.LastUpdatedAt,
			&transaction.ConfirmedAt,
			&transaction.Category.Name,
			&transaction.Type.Name,
			&transaction.Account.Name,
			&accountDestName,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("could not scan row: %w", err)
		}

		if accountDestID.Valid {
			transaction.AccountDest = &Account{
				ID: int(accountDestID.Int64),
			}

			if accountDestName.Valid {
				transaction.AccountDest.Name = accountDestName.String
			}
		}

		transaction.Amount, _, err =
			new(big.Float).Parse(string(amount), 10)

		if err != nil {
			log.Printf("Error parsing amount: %v", err)
			return nil, fmt.Errorf("could not parse amount: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transactions, nil
}

func (r *Repository) GetTransactionByID(
	ctx context.Context, // ADD
	id int,
) (*Transaction, error) {

	var t Transaction
	var amount []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id, parent_id, trade_date, type_id, category_id,
			account_id, account_dest_id, amount, description,
			status, source, inserted_at, last_updated_at, confirmed_at
		FROM transactions
		WHERE id = $1
	`, id).Scan(
		&t.ID,
		&t.ParentID,
		&t.TradeDate,
		&t.TypeID,
		&t.CategoryID,
		&t.AccountID,
		&t.AccountDestID,
		&amount,
		&t.Description,
		&t.Status,
		&t.Source,
		&t.InsertedAt,
		&t.LastUpdatedAt,
		&t.ConfirmedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction with ID %d not found", id)
		}
		log.Printf("Error fetching transaction by ID: %v", err)
		return nil, fmt.Errorf("could not retrieve transaction: %w", err)
	}

	t.Amount, _, err =
		new(big.Float).Parse(string(amount), 10)
	if err != nil {
		log.Printf("Error parsing amount: %v", err)
		return nil, fmt.Errorf("could not parse amount: %w", err)
	}

	return &t, nil
}

func (r *Repository) GetTransactionByIDDTO(
	ctx context.Context, // ADD
	id int,
) (*TransactionDTO, error) {

	var transaction TransactionDTO
	var amount []byte
	var accountDestID sql.NullInt64
	var accountDestName sql.NullString

	err := r.db.QueryRowContext(ctx, queryGetTransactionByID, id).Scan(
		&transaction.ID,
		&transaction.ParentID,
		&transaction.TradeDate,
		&transaction.Type.ID,
		&transaction.Category.ID,
		&transaction.Account.ID,
		&accountDestID,
		&amount,
		&transaction.Description,
		&transaction.Status,
		&transaction.Source,
		&transaction.InsertedAt,
		&transaction.LastUpdatedAt,
		&transaction.ConfirmedAt,
		&transaction.Category.Name,
		&transaction.Type.Name,
		&transaction.Account.Name,
		&accountDestName,
	)

	if accountDestID.Valid {
			transaction.AccountDest = &Account{
				ID: int(accountDestID.Int64),
			}

			if accountDestName.Valid {
				transaction.AccountDest.Name = accountDestName.String
			}
		}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction with ID %d not found", id)
		}
		log.Printf("Error fetching transaction by ID: %v", err)
		return nil, fmt.Errorf("could not retrieve transaction: %w", err)
	}

	transaction.Amount, _, err =
		new(big.Float).Parse(string(amount), 10)
	if err != nil {
		log.Printf("Error parsing amount: %v", err)
		return nil, fmt.Errorf("could not parse amount: %w", err)
	}

	return &transaction, nil
}

//
// CREATE (inside SQL transaction)
//
func (r *Repository) CreateTransaction(
	tx *sql.Tx,
	t Transaction,
) (Transaction, error) {

	log.Printf("%v", t)

	var transactionID int

	query := queryInsertTransactionAuto
	if t.Source == "Manual" {
		query = queryInsertTransactionManual
	}

	// Convert big.Float → string (DB-safe)
	amount := t.Amount.Text('f', 2)

	err := tx.QueryRow(
		query,
		t.ParentID,
		t.TradeDate,
		t.TypeID,
		t.CategoryID,
		t.AccountID,
		t.AccountDestID,
		amount,
		t.Description,
		t.Status,
		t.Source,
	).Scan(&transactionID)

	if err != nil {
		log.Printf("Error inserting transaction: %v", err)
		return Transaction{}, fmt.Errorf("could not insert transaction: %w", err)
	}

	t.ID = transactionID
	return t, nil
}

//
// UPDATE
//

func (r *Repository) UpdateTransaction(
	ctx context.Context, // ADD
	id int,
	dto UpdateTransactionDBDTO,
) (Transaction, error) {

	query := `
		UPDATE transactions
		SET
			trade_date      = COALESCE($1, trade_date),
			type_id         = COALESCE($2, type_id),
			category_id     = COALESCE($3, category_id),
			account_id      = COALESCE($4, account_id),
			account_dest_id = COALESCE($5, account_dest_id),
			amount          = COALESCE($6, amount),
			description     = COALESCE($7, description),
			status          = COALESCE($8, status),
			last_updated_at = NOW()
		WHERE id = $9
		RETURNING *;
	`

	var t Transaction

	err := r.db.QueryRowContext(
		ctx,
		query,
		dto.TradeDate,
		dto.Amount,
		dto.Description,
		dto.CategoryID,
		dto.AccountID,
		dto.AccountDestID,
		dto.Status,
		id,
	).Scan(
		&t.ID,
		&t.ParentID,
		&t.TradeDate,
		&t.TypeID,
		&t.CategoryID,
		&t.AccountID,
		&t.AccountDestID,
		&t.Amount,
		&t.Description,
		&t.Status,
		&t.Source,
		&t.InsertedAt,
		&t.LastUpdatedAt,
		&t.ConfirmedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return Transaction{}, fmt.Errorf("transaction not found")
		}
		log.Printf("update transaction failed by ID: %v", err)
		return Transaction{}, fmt.Errorf("update transaction failed: %w", err)
	}
	return t, nil
}

func (r *Repository) UpdateTransactionTx(
	ctx context.Context,
	tx *sql.Tx,
	id int,
	dto UpdateTransactionDBDTO,
) (Transaction, error) {

	query := `
		UPDATE transactions
		SET
			trade_date      = COALESCE($1, trade_date),
			type_id         = COALESCE($2, type_id),
			category_id     = COALESCE($3, category_id),
			account_id      = COALESCE($4, account_id),
			account_dest_id =
							CASE
								WHEN $6 = true THEN NULL
								ELSE COALESCE($5, account_dest_id)
							END,
			amount          = COALESCE($7, amount),
			description     = COALESCE($8, description),
			status          = COALESCE($9, status),
			last_updated_at = NOW()
		WHERE id = $10
		RETURNING *;
	`

	var t Transaction
	var amountStr string

	err := tx.QueryRowContext(
		ctx,
		query,
		dto.TradeDate,
		dto.TypeID,
		dto.CategoryID,
		dto.AccountID,
		dto.AccountDestID,
		dto.ClearAccountDest,
		dto.Amount,
		dto.Description,
		dto.Status,
		id,
	).Scan(
		&t.ID,
		&t.ParentID,
		&t.TradeDate,
		&t.TypeID,
		&t.CategoryID,
		&t.AccountID,
		&t.AccountDestID,
		&amountStr,
		&t.Description,
		&t.Status,
		&t.Source,
		&t.InsertedAt,
		&t.LastUpdatedAt,
		&t.ConfirmedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return Transaction{}, fmt.Errorf("transaction not found")
		}
		log.Printf("update transaction failed by ID: %v", err)
		return Transaction{}, fmt.Errorf("update transaction failed: %w", err)
	}

	amount := new(big.Float)
	_, ok := amount.SetString(amountStr)
	if !ok {
		return Transaction{}, fmt.Errorf("invalid numeric value for amount: %s", amountStr)
	}

	t.Amount = amount

	return t, nil
}

//
// DELETE
//

func (r *Repository) DeleteTransaction(
	ctx context.Context, // ADD
	id int,
) error {

	_, err := r.db.ExecContext(ctx,
		`DELETE FROM transactions WHERE id = $1`,
		id,
	)
	if err != nil {
		log.Printf("Error deleting transaction: %v", err)
		return fmt.Errorf("could not delete transaction: %w", err)
	}

	return nil
}

func (r *Repository) GetCategoryAggregatesByType(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
	typeName string,
) ([]CategoryAggregate, error) {

	query := `
		SELECT
			c.name AS category_name,
			COALESCE(SUM(t.amount), 0) AS total_amount
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		JOIN transaction_types tt ON t.type_id = tt.id
		WHERE t.trade_date BETWEEN $1 AND $2
		  AND tt.name = $3
		GROUP BY c.name
		ORDER BY total_amount DESC;
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate, typeName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CategoryAggregate

	for rows.Next() {
		var categoryName string
		var totalStr string

		if err := rows.Scan(&categoryName, &totalStr); err != nil {
			return nil, err
		}

		totalFloat, _, err := big.ParseFloat(totalStr, 10, 128, big.ToNearestEven)
		if err != nil {
			return nil, err
		}

		results = append(results, CategoryAggregate{
			CategoryName: categoryName,
			TotalAmount:  totalFloat,
		})
	}

	return results, nil
}