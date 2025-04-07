package accounts

import (
	"drfin/app/accountgroups"
	"math/big"
	"time"
)

type Account struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Balance        big.Float `json:"balance"`
	AccountGroupID int       `json:"account_group_id"`
	Status         string    `json:"status"`
	InsertedAt     time.Time `json:"inserted_at"`
	ClosedAt       time.Time `json:"closed_at"`
}

type AccountDTO struct {
	ID           int                        `json:"id"`
	Name         string                     `json:"name"`
	Balance      big.Float                  `json:"balance"`
	AccountGroup accountgroups.AccountGroup `json:"account_group"`
	Status       string                     `json:"status"`
	InsertedAt   time.Time                  `json:"inserted_at"`
	ClosedAt     time.Time                  `json:"closed_at"`
}
