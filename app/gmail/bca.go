package gmail

import (
	"context"
	"drfin/app/transactions"
	"fmt"
	"strings"
)

// accountFromSourceOfFund maps the BCA "Source of Fund" field to an owned account name.
// A missing (empty) value falls back to "Tahapan Gold" to preserve the previous behavior.
// A non-empty value that matches no known pattern returns false so the caller can fail
// instead of posting to the wrong account.
func accountFromSourceOfFund(sourceOfFund string) (string, bool) {
	sourceOfFund = strings.TrimSpace(sourceOfFund)
	switch {
	case sourceOfFund == "":
		return "Tahapan Gold", true
	case strings.Contains(sourceOfFund, "BCA KKI"):
		return "CC BCA", true
	case strings.Contains(sourceOfFund, "TAHAPAN - 0840"):
		return "Tahapan Gold", true
	}
	return "", false
}

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
		account, ok := accountFromSourceOfFund(mapData["Source of Fund"])
		if !ok {
			return fmt.Errorf("unrecognized source of fund for QRIS payment")
		}
		dto := &transactions.CreateTransactionDTO{
			TradeDate:   tradeDate,
			Type:        "Expense",
			Category:    "Unknown",
			Account:     account,
			Amount:      amount,
			Description: mapData["Payment to"],
			Status:      "Waiting",
			Source:      "Email",
		}
		return s.saveTransaction(ctx, dto)

	case "Credit Card & Paylater - BCA":
		amount, ok := ParseAmount(mapData["Total Payment"], EmailTitleTransaksiBCA)
		if !ok {
			return fmt.Errorf("%f invalid number format", amount)
		}
		dto := &transactions.CreateTransactionDTO{
			TradeDate:   tradeDate,
			Type:        "Transfer",
			Category:    "Unknown",
			Account:     "Tahapan Gold",
			AccountDest: "CC BCA",
			Amount:      amount,
			Description: mapData["Card No. / Customer No."] + " " + mapData["Name"],
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

	case "Pocket Transaction":
		dto, err := buildPocketTransferDTO(mapData, tradeDate)
		if err != nil {
			return err
		}
		return s.saveTransaction(ctx, dto)
	}
	return nil
}
