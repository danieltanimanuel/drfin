package accounts

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrAccountNotFound = errors.New("account not found")
)

// numeric(24,2): up to 22 integer digits and 2 decimals.
var amountPattern = regexp.MustCompile(`^-?\d{1,22}(\.\d{1,2})?$`)

// NormalizeAmount validates a plain decimal string and returns it with exactly
// 2 decimals. No float conversion is involved.
func NormalizeAmount(s string) (string, error) {
	s = strings.TrimSpace(s)
	if !amountPattern.MatchString(s) {
		return "", fmt.Errorf("%w: amount %q must be a plain number with at most 2 decimals", ErrInvalidInput, s)
	}
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		return "", fmt.Errorf("%w: amount %q could not be parsed", ErrInvalidInput, s)
	}
	return r.FloatString(2), nil
}