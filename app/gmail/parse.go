package gmail

import (
	"fmt"
	"math/big"
	"strings"
	"time"
)

func prerequisiteParseDate(datetime string) (datetimeClean string) {
	if strings.Contains(datetime, "Agt") {
		return strings.ReplaceAll(datetime, "Agt", "Aug")
	} else if strings.Contains(datetime, "Des") {
		return strings.ReplaceAll(datetime, "Des", "Dec")
	} else if strings.Contains(datetime, "Mei") {
		return strings.ReplaceAll(datetime, "Mei", "May")
	}
	return datetime
}

func ParseDate(datetime string, emailTitle EmailTitle) (time.Time, error) {
	var tradeDate time.Time
	var err error
	datetime = prerequisiteParseDate(datetime)
	switch emailTitle {
	case EmailTitleTransaksiBCA:
		tradeDate, err = time.Parse(string(LayoutDateTimeDDMMMYYYYHHMMSS), datetime)
	case EmailTitleTransaksiBlu:
		tradeDate, err = time.Parse(string(LayoutDateTimeDDMMMYYYYHHMMSS), datetime[:len(datetime)-4]) //Remove WIB for Blu, CC BCA
	case EmailTitleCreditCardBCA:
		tradeDate, err = time.Parse(string(LayoutDateTimeCCBCA), datetime[:len(datetime)-4]) //Remove WIB for Blu, CC BCA
	default:
		return time.Time{}, fmt.Errorf("email title not mapping yet for conversion of trade date: %v", err)
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %v", err)
	}
	return tradeDate, err
}

func ParseAmount(number string, emailTitle EmailTitle) (*big.Float, bool) {
	normalized := ""
	switch emailTitle {
	case EmailTitleTransaksiBCA:
		normalized = strings.ReplaceAll(number[4:], ",", "") // "IDR 108,200.00 > 108200.00"
	case EmailTitleCreditCardBCA:
		if strings.Contains(number, "IDR") {
			normalized = strings.ReplaceAll(number[4:], ".", "") // "IDR 108,200.00 > 108200.00"
		} else {
			normalized = strings.ReplaceAll(number[2:], ".", "")  // "Rp108.200,00 > 108200,00"
			normalized = strings.ReplaceAll(normalized, ",", ".") // "108200,00 > 108200.00"
		}
	default:
		normalized = strings.ReplaceAll(number[2:], ".", "")  // "Rp108.200,00 > 108200,00"
		normalized = strings.ReplaceAll(normalized, ",", ".") // "108200,00 > 108200.00"
	}
	amount, ok := new(big.Float).SetString(normalized)
	if !ok {
		return &big.Float{}, ok
	}
	return amount, ok
}
