package accounts

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAllAccounts(ownerID *int) ([]AccountDTO, error) {
	return s.repository.GetAllAccounts(ownerID)
}

func (s *Service) GetAccountByAccountNumber(accountNumber string, ownerID int) (*AccountDTO, error) {
	return s.repository.GetAccountByAccountNumber(accountNumber, ownerID)
}

func (s *Service) GetAccountOwnerID(accountID int) (int, error) {
	return s.repository.GetOwnerID(accountID)
}

func (s *Service) CreateReceivable(
	ctx context.Context,
	name string,
	ownerID int,
) (*Account, error) {

	tx, err := s.repository.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	account, err := s.repository.CreateReceivableAccount(tx, name, ownerID)
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

func validateCreateAccount(req CreateAccountRequest) (name, balance string, creditLimit *string, err error) {
	name = strings.TrimSpace(req.Name)
	if name == "" || utf8.RuneCountInString(name) > 255 {
		return "", "", nil, fmt.Errorf("%w: name is required (max 255 characters)", ErrInvalidInput)
	}
	if req.AccountGroupID <= 0 {
		return "", "", nil, fmt.Errorf("%w: account group is required", ErrInvalidInput)
	}

	balance = "0.00"
	if strings.TrimSpace(req.InitialBalance) != "" {
		if balance, err = NormalizeAmount(req.InitialBalance); err != nil {
			return "", "", nil, err
		}
	}

	if strings.TrimSpace(req.CreditLimit) != "" {
		cl, err := NormalizeAmount(req.CreditLimit)
		if err != nil {
			return "", "", nil, err
		}
		if strings.HasPrefix(cl, "-") {
			return "", "", nil, fmt.Errorf("%w: credit limit cannot be negative", ErrInvalidInput)
		}
		creditLimit = &cl
	}
	return name, balance, creditLimit, nil
}

func (s *Service) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	name, balance, creditLimit, err := validateCreateAccount(req)
	if req.OwnerID <= 0 {
		return nil, fmt.Errorf("%w: owner is required", ErrInvalidInput)
	}
	if err != nil {
		return nil, err
	}

	tx, err := s.repository.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() // no-op after a successful Commit

	ok, err := s.repository.AccountGroupExists(tx, req.AccountGroupID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%w: account group %d does not exist", ErrInvalidInput, req.AccountGroupID)
	}

	taken, err := s.repository.AccountNameExists(tx, name, req.OwnerID)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, fmt.Errorf("%w: an account named %q already exists", ErrInvalidInput, name)
	}

	account, err := s.repository.CreateAccount(tx, name, req.AccountGroupID, balance, creditLimit, req.OwnerID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit account: %w", err)
	}
	return account, nil
}