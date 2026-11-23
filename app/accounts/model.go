package accounts

import (
	"database/sql"
	"drfin/app/accountgroups"
	"html/template"
	"math/big"
	"time"
)

type Account struct {
	ID             int          `json:"id"`
	Name           string       `json:"name"`
	AccountNumber  string       `json:"account_number"`
	Balance        *big.Float   `json:"balance"`
	AccountGroupID int          `json:"account_group_id"`
	Status         string       `json:"status"`
	InsertedAt     time.Time    `json:"inserted_at"`
	ClosedAt       sql.NullTime `json:"closed_at"`
	CreditLimit    *big.Float   `json:"credit_limit"`
}
type AccountDTO struct {
	ID            int                        `json:"id"`
	Name          string                     `json:"name"`
	AccountNumber string                     `json:"account_number"`
	Balance       *big.Float                 `json:"balance"`
	AccountGroup  accountgroups.AccountGroup `json:"account_group"`
	Status        string                     `json:"status"`
	InsertedAt    time.Time                  `json:"inserted_at"`
	ClosedAt      sql.NullTime               `json:"closed_at"`
	CreditLimit   *big.Float                 `json:"credit_limit"`
}

type AccountPage struct {
	Accounts     []AccountDTO
	AccountsJSON template.JS // safe to inject JSON into template
}

type CreateAccountRequest struct {
	Name           string `json:"name"`
	AccountGroupID int    `json:"account_group_id"`
	InitialBalance string `json:"initial_balance"` // optional, defaults to 0
	CreditLimit    string `json:"credit_limit"`    // optional
	OwnerID        int    `json:"owner_id"`
}