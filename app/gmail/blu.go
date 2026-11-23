package gmail

import (
	"context"
	"database/sql"
	"drfin/app/shared/logdr"
	"drfin/app/transactions"
	"fmt"
)

func toMap(details []GmailDetailTransaction) map[string]string {
	m := make(map[string]string, len(details))
	for _, d := range details {
		m[d.Key] = d.Value
	}
	return m
}

func (s *Service) handleTransaksiBlu(ctx context.Context, data *GmailExtractedData) error {
	logdr.Info("log detail transactions", "detailTransactions", data)
	mapData := toMap(data.DetailTransaction)
	account, err := s.accountService.GetAccountByAccountNumber(mapData["Beneficiary"],data.OwnerID)
	var transactionType string
	switch {
	case err == sql.ErrNoRows:
		transactionType = "Expense"
		err = nil
	case err != nil: //Technical Error
		return err
	default:
		transactionType = "Transfer"
	}
	tradeDate, err := ParseDate(mapData["Tgl & Jam Transaksi"], EmailTitleTransaksiBlu)
	if err != nil {
		return err
	}

	switch mapData["Tipe Transaksi"] {
	case "Top Up Flazz":
		amount, ok := ParseAmount(mapData["Total Bayar"], EmailTitleTransaksiBlu)
		if !ok {
			return fmt.Errorf("%f invalid number format", amount)
		}
		dto := &transactions.CreateTransactionDTO{
			TradeDate:   tradeDate,
			Type:        "Expense",
			Category:    "Transport",
			Account:     "Blu",
			Amount:      amount,
			Description: mapData["Tipe Transaksi"] + " " + mapData["Flazz"],
			Status:      "Waiting",
			Source:      "Email",
		}
		return s.saveTransaction(ctx, data.OwnerID, dto)
	case "BI-FAST":
		amount, ok := ParseAmount(mapData["Total"], EmailTitleTransaksiBlu)
		if !ok {
			return fmt.Errorf("%f invalid number format", amount)
		}
		accountName := ""
		if account != nil {
			accountName = account.Name
		}
		dto := &transactions.CreateTransactionDTO{
			TradeDate:   tradeDate,
			Type:        transactionType,
			Category:    "Unknown",
			Account:     "Blu",
			AccountDest: accountName,
			Amount:      amount,
			Description: mapData["Beneficiary Full"],
			Status:      "Waiting",
			Source:      "Email",
		}
		return s.saveTransaction(ctx, data.OwnerID, dto)
	case "QRIS":
		amount, ok := ParseAmount(mapData["Nominal Tagihan"], EmailTitleTransaksiBlu)
		if !ok {
			return fmt.Errorf("%f invalid number format", amount)
		}
		accountName := ""
		if account != nil {
			accountName = account.Name
		}
		dto := &transactions.CreateTransactionDTO{
			TradeDate:   tradeDate,
			Type:        transactionType,
			Category:    "Unknown",
			Account:     "Blu",
			AccountDest: accountName,
			Amount:      amount,
			Description: mapData["Beneficiary Full"],
			Status:      "Waiting",
			Source:      "Email",
		}
		return s.saveTransaction(ctx, data.OwnerID, dto)
	}
	return nil
}
