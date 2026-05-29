package dashboard

import (
	"drfin/app/categories"
	"drfin/app/transactions"
)

type EditTransactionView struct {
	Transaction *transactions.TransactionDTO
	Categories  []categories.Category
}