package accounts

import (
	"context"
	"fmt"
	"log"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAllAccounts() ([]AccountDTO, error) {
	return s.repository.GetAllAccounts()
}

func (s *Service) GetAccountByAccountNumber(accountNumber string) (*AccountDTO, error) {
	return s.repository.GetAccountByAccountNumber(accountNumber)
}

func (s *Service) CreateReceivable(
	ctx context.Context,
	name string,
) (*Account, error) {

	tx, err := s.repository.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	account, err := s.repository.CreateReceivableAccount(tx, name)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("rollback failed: %v", rollbackErr)
		}

		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit receivable account: %w", err)
	}

	return account, nil
}
