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
	Amount        float64    `json:"amount"`
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
	AccountDest   Account         `json:"account_dest_id"`
	Amount        *big.Float      `json:"amount"` //Why *?
	Description   string          `json:"description"`
	Status        string          `json:"status"`
	Source        string          `json:"source"`
	InsertedAt    time.Time       `json:"inserted_at"`
	LastUpdatedAt time.Time       `json:"last_updated_at"`
	ConfirmedAt   *time.Time      `json:"confirmed_at"`
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
	ID   *int    `json:"id"`
	Name *string `json:"name"`
}
