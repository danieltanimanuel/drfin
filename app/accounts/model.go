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
	AccountNumber  string       `json:account_number`
	Balance        *big.Float   `json:"balance"`
	AccountGroupID int          `json:"account_group_id"`
	Status         string       `json:"status"`
	InsertedAt     time.Time    `json:"inserted_at"`
	ClosedAt       sql.NullTime `json:"closed_at"`
}
type AccountDTO struct {
	ID            int                        `json:"id"`
	Name          string                     `json:"name"`
	AccountNumber string                     `json:account_number`
	Balance       *big.Float                 `json:"balance"`
	AccountGroup  accountgroups.AccountGroup `json:"account_group"`
	Status        string                     `json:"status"`
	InsertedAt    time.Time                  `json:"inserted_at"`
	ClosedAt      sql.NullTime               `json:"closed_at"`
}

type AccountPage struct {
	Accounts     []AccountDTO
	AccountsJSON template.JS // safe to inject JSON into template
}
