package gmail

import (
	"context"
	"drfin/app/transactions"
	"fmt"
)

func (s *Service) saveTransaction(ctx context.Context, ownerID int, dto *transactions.CreateTransactionDTO) error {
	dto.OwnerID = ownerID
	_, err := s.transactionService.CreateTransaction(ctx, *dto)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}
	return nil
}
