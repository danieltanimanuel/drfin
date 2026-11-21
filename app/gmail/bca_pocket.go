package gmail

import (
	"drfin/app/transactions"
	"fmt"
	"strings"
	"time"
)

const (
	mainAccountName  = "Tahapan Gold"
	pocketAccountPfx = "MyBCAPocket "
)

func pocketAccountName(pocket string) (string, error) {
	pocket = strings.TrimSpace(pocket)
	if pocket == "" {
		return "", fmt.Errorf("pocket transaction: missing pocket name")
	}
	return pocketAccountPfx + pocket, nil
}

// buildPocketTransferDTO maps a BCA "Pocket Transaction" email to a Transfer DTO.
// Layouts (by which keys are present):
//   - Beneficiary Pocket + Source of Fund   : Tahapan Gold -> pocket
//   - Beneficiary Pocket + Source Pocket    : pocket -> pocket
//   - Source Pocket + Beneficiary Account   : pocket -> Tahapan Gold (withdrawal)
//
// Anything else returns an error so the email is recorded as failed, not silently skipped.
func buildPocketTransferDTO(m map[string]string, tradeDate time.Time) (*transactions.CreateTransactionDTO, error) {
	// ParseAmount slices the string, so guard against missing/short values.
	if len(strings.TrimSpace(m["Amount"])) < 5 {
		return nil, fmt.Errorf("pocket transaction: missing amount")
	}
	amount, ok := ParseAmount(strings.TrimSpace(m["Amount"]), EmailTitleTransaksiBCA)
	if !ok {
		return nil, fmt.Errorf("pocket transaction: invalid amount format")
	}

	var src, dest string
	var err error
	_, hasBenPocket := m["Beneficiary Pocket"]
	_, hasSrcPocket := m["Source Pocket"]
	_, hasSrcFund := m["Source of Fund"]
	_, hasBenAccount := m["Beneficiary Account"]

	switch {
	case hasBenPocket && hasSrcFund: // main -> pocket
		src = mainAccountName
		dest, err = pocketAccountName(m["Beneficiary Pocket"])
	case hasBenPocket && hasSrcPocket: // pocket -> pocket
		if src, err = pocketAccountName(m["Source Pocket"]); err == nil {
			dest, err = pocketAccountName(m["Beneficiary Pocket"])
		}
	case hasSrcPocket && hasBenAccount: // pocket withdrawal -> always Tahapan Gold
		dest = mainAccountName
		src, err = pocketAccountName(m["Source Pocket"])
	default:
		return nil, fmt.Errorf("pocket transaction: unrecognized layout")
	}
	if err != nil {
		return nil, err
	}
	if src == dest {
		return nil, fmt.Errorf("pocket transaction: source and destination are the same")
	}

	return &transactions.CreateTransactionDTO{
		TradeDate:   tradeDate,
		Type:        "Transfer",
		Category:    "Unknown",
		Account:     src,
		AccountDest: dest,
		Amount:      amount,
		Description: "Pocket: " + strings.TrimPrefix(src, pocketAccountPfx) + " -> " + strings.TrimPrefix(dest, pocketAccountPfx),
		Status:      "Waiting",
		Source:      "Email",
	}, nil
}