package gmail

import (
	"drfin/app/transactions"
	"fmt"
	"math/big"
	"strings"
	"time"
)

const (
	EmailTitleCreditCardBCA = "Credit Card Transaction Notification"
	EmailTitleTransaksiBlu  = "Transaksimu Pakai blu Berhasil"
)

type Service struct {
	transactionService *transactions.Service
}

func NewService(transactionService *transactions.Service) *Service {
	return &Service{transactionService: transactionService}
}

func (s *Service) createTransactionFromGmail(extractedData *GmailExtractedData) error {
	switch extractedData.Title {
	case EmailTitleCreditCardBCA:
		tradeDate, err := convertStringToDate(extractedData.DetailTransaction[5].Value)
		if err != nil {
			return fmt.Errorf("error parsing date: %v", err)
		}
		amount, ok := convertStringToFloat(extractedData.DetailTransaction[6].Value)
		if !ok {
			return fmt.Errorf("%f invalid number format", amount)
		}
		createTransactionDTO := &transactions.CreateTransactionDTO{
			TradeDate:   tradeDate,
			Type:        "Expense",
			Category:    "Misc",
			Account:     "CC BCA",
			AccountDest: "",
			Amount:      amount,
			Description: extractedData.DetailTransaction[2].Value,
			Status:      "Waiting",
			Source:      "Email",
		}
		_, err = s.transactionService.CreateTransaction(*createTransactionDTO)
		if err != nil {
			return fmt.Errorf("error creating transaction : %v", err)
		}
	case EmailTitleTransaksiBlu:
		// tradeDate, err := convertStringToDate(extractedData.DetailTransaction[5].Value)
		// if err != nil {
		// 	return fmt.Errorf("error parsing date: %v", err)
		// }
		// amount, ok := convertStringToFloat(extractedData.DetailTransaction[6].Value)
		// if !ok {
		// 	return fmt.Errorf("%f invalid number format", amount)
		// }
	}
	return nil
}

func convertStringToDate(datetime string) (time.Time, error) {
	layout := "02-01-2006 15:04:05"
	tradeDate, err := time.Parse(layout, datetime[:len(datetime)-4])
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %v", err)
	}
	return tradeDate, err
}

func convertStringToFloat(number string) (*big.Float, bool) {
	// Step 1: Normalize the number string
	normalized := strings.ReplaceAll(number[2:], ".", "") // "Rp108.200,00 > 108200,00"
	normalized = strings.ReplaceAll(normalized, ",", ".") // "108200,00 > 108200.00"
	amount, ok := new(big.Float).SetString(normalized)
	if !ok {
		return &big.Float{}, ok
	}
	return amount, ok
}
