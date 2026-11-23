package gmail

import (
	"context"
	"drfin/app/shared/logdr"
	"drfin/app/transactions"
	"fmt"
	"strings"
)

// isSkippedCurrency reports whether the amount's currency prefix
// (e.g. "MYR" in "MYR 99,47") is one we deliberately do not import yet.
// See "TODO: foreign-currency credit-card transactions" in docs/FEATURES.md.
func isSkippedCurrency(amountText string) bool {
	fields := strings.Fields(amountText)
	if len(fields) == 0 {
		return false
	}
	return strings.EqualFold(fields[0], "MYR")
}

func (s *Service) handleCreditCardBCA(ctx context.Context, data *GmailExtractedData) error {
	logdr.Info("log detail transactions", "detailTransactions", data)
	mapData := toMap(data.DetailTransaction)

	// Skip foreign-currency emails before any parsing so they are not treated as failures.
	if isSkippedCurrency(mapData["Sejumlah"]) {
		return ErrSkipEmail
	}

	tradeDate, err := ParseDate(mapData["Pada Tanggal"], EmailTitleCreditCardBCA)
	if err != nil {
		return err
	}
	amount, ok := ParseAmount(mapData["Sejumlah"], EmailTitleCreditCardBCA)
	if !ok {
		return fmt.Errorf("%f invalid number format", amount)
	}

	dto := &transactions.CreateTransactionDTO{
		TradeDate:   tradeDate,
		Type:        "Expense",
		Category:    "Unknown",
		Account:     "CC BCA",
		Amount:      amount,
		Description: mapData["Merchant / ATM"],
		Status:      "Waiting",
		Source:      "Email",
	}
	return s.saveTransaction(ctx, data.OwnerID, dto)
}
