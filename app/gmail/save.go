package gmail

import (
	"context"
	"drfin/app/transactions"
	"fmt"
)

func (s *Service) saveTransaction(ctx context.Context, dto *transactions.CreateTransactionDTO) error {

	_, err := s.transactionService.CreateTransaction(ctx, *dto)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}
	return nil
}
