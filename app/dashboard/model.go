package dashboard

import (
	"drfin/app/categories"
	"drfin/app/transactions"
)

type EditTransactionView struct {
	Transaction *transactions.TransactionDTO
	Categories  []categories.Category
	Accounts    interface{}
}

type CreateTransactionView struct {
    Categories interface{}
    Accounts    interface{}
}