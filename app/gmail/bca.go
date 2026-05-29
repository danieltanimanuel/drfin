package gmail

import (
	"context"
	"drfin/app/transactions"
	"fmt"
)

func (s *Service) handleTransaksiBCA(ctx context.Context, data *GmailExtractedData) error {
	mapData := toMap(data.DetailTransaction)
	tradeDate, err := ParseDate(mapData["Transaction Date"], EmailTitleTransaksiBCA)
	if err != nil {
		return err
	}

	transactionType := mapData["Transfer Type"]

	if value, exists := mapData["Transaction Type"]; exists {
		transactionType = value
	}

	switch transactionType {
	case "Transfer to BCA Account":
		amount, ok := ParseAmount(mapData["Transfer Amount"], EmailTitleTransaksiBCA)
		if !ok {
			return fmt.Errorf("%f invalid number format", amount)
		}
		dto := &transactions.CreateTransactionDTO{
			TradeDate:   tradeDate,
			Type:        "Expense",
			Category:    "Unknown",
			Account:     "Tahapan Gold",
			Amount:      amount,
			Description: "Transfer ke " + mapData["Beneficiary Name"],
			Status:      "Waiting",
			Source:      "Email",
		}
		return s.saveTransaction(ctx, dto)

	case "QRIS Payment":
		amount, ok := ParseAmount(mapData["Total Payment"], EmailTitleTransaksiBCA)
		if !ok {
			return fmt.Errorf("%f invalid number format", amount)
		}
		dto := &transactions.CreateTransactionDTO{
			TradeDate:   tradeDate,
			Type:        "Expense",
			Category:    "Unknown",
			Account:     "Tahapan Gold",
			Amount:      amount,
			Description: mapData["Payment to"],
			Status:      "Waiting",
			Source:      "Email",
		}
		return s.saveTransaction(ctx, dto)

	case "Flazz Top Up":
		amount, ok := ParseAmount(mapData["Top Up Amount"], EmailTitleTransaksiBCA)
		if !ok {
			return fmt.Errorf("%f invalid number format", amount)
		}
		dto := &transactions.CreateTransactionDTO{
			TradeDate:   tradeDate,
			Type:        "Expense",
			Category:    "Transport",
			Account:     "Tahapan Gold",
			Amount:      amount,
			Description: mapData["Transaction Type"] + " " + mapData["Flazz Card Number"],
			Status:      "Waiting",
			Source:      "Email",
		}
		return s.saveTransaction(ctx, dto)
	}
	return nil
}
