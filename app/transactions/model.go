package transactions

import (
	"math/big"
	"time"
)

type Transaction struct {
	ID            int        `json:"id"`
	ParentID      int        `json:"parent_id"`
	TradeDate     time.Time  `json:"trade_date"`
	TypeID        int        `json:"type_id"`
	CategoryID    int        `json:"category_id"`
	AccountID     int        `json:"account_id"`
	AccountDestID *int       `json:"account_dest_id"`
	Amount        *big.Float `json:"amount"`
	Description   string     `json:"description"`
	Status        string     `json:"status"`
	Source        string     `json:"source"`
	InsertedAt    time.Time  `json:"inserted_at"`
	LastUpdatedAt time.Time  `json:"last_updated_at"`
	ConfirmedAt   *time.Time `json:"confirmed_at"`
}

type TransactionDTO struct {
	ID            int             `json:"id"`
	ParentID      int             `json:"parent_id"`
	TradeDate     time.Time       `json:"trade_date"`
	Type          TransactionType `json:"type"`
	Category      Category        `json:"category"`
	Account       Account         `json:"account_id"`
	AccountDest   *Account         `json:"account_dest_id"`
	Amount        *big.Float      `json:"amount"` //Why *?
	Description   string          `json:"description"`
	Status        string          `json:"status"`
	Source        string          `json:"source"`
	InsertedAt    time.Time       `json:"inserted_at"`
	LastUpdatedAt time.Time       `json:"last_updated_at"`
	ConfirmedAt   *time.Time      `json:"confirmed_at"`
}

type DailyTransactionsDTO struct {
	Date         time.Time
	IncomeTotal  *big.Float
	ExpenseTotal *big.Float
	Transactions []TransactionDTO
}

type CreateTransactionDTO struct {
	TradeDate   time.Time  `json:"trade_date"`
	Type        string     `json:"type"`
	Category    string     `json:"category"`
	Account     string     `json:"account_id"`
	AccountDest string     `json:"account_dest_id"`
	Amount      *big.Float `json:"amount"` //Why *?
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Source      string     `json:"source"`
}

// ✅ Enhancement:
// Using pointers lets you distinguish:
// nil → field not updated
// value → field updated
type UpdateTransactionDTO struct {
    TradeDate        *time.Time `json:"trade_date,omitempty"`
    Type             *string    `json:"type,omitempty"`
    Category         *string    `json:"category,omitempty"`
    Account          *string    `json:"account,omitempty"`
    AccountDest      *string    `json:"account_dest,omitempty"`
    NewReceivable    *string    `json:"new_receivable,omitempty"`
    Amount           *big.Float `json:"amount,omitempty"`
    Description      *string    `json:"description,omitempty"`
    Status           *string    `json:"status,omitempty"`
    Source           *string    `json:"source,omitempty"`
}

// Internal DB DTO (ID-based, repository-facing)
type UpdateTransactionDBDTO struct {
	TradeDate     *time.Time
	TypeID        *int
	CategoryID    *int
	AccountID     *int
	AccountDestID *int
	ClearAccountDest bool
	Amount        *big.Float
	Description   *string
	Status        *string
}

type TransactionType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryAggregate struct {
	CategoryName string
	TotalAmount  *big.Float
}

type AdjustBalanceRequest struct {
	NewBalance string `json:"new_balance"`
	Note       string `json:"note"`
}

type AdjustBalanceResult struct {
	AccountID     int    `json:"account_id"`
	OldBalance    string `json:"old_balance"`
	NewBalance    string `json:"new_balance"`
	Delta         string `json:"delta"`
	TransactionID int    `json:"transaction_id,omitempty"` // 0 when nothing changed
}
