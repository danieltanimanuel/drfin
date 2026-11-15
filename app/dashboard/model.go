package dashboard

import (
	"drfin/app/accounts"
	"drfin/app/categories"
	"drfin/app/transactions"
	"drfin/app/transactiontypes"
)

type EditTransactionView struct {
    Transaction *transactions.TransactionDTO
    Categories  []categories.Category
    Accounts    []accounts.AccountDTO
    TransactionTypes []transactiontypes.TransactionType
}

type CreateTransactionView struct {
    Categories interface{}
    Accounts    interface{}
}