package gmail

import (
	"context"
	"drfin/app/shared/logdr"
	"drfin/app/transactions"
	"fmt"
)

func (s *Service) handleCreditCardBCA(ctx context.Context, data *GmailExtractedData) error {
	logdr.Info("log detail transactions", "detailTransactions", data)
	mapData := toMap(data.DetailTransaction)
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
	return s.saveTransaction(ctx, dto)
}
