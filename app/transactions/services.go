package transactions

import (
	"drfin/app/accounts"
	"fmt"
	"log"
)

type Service struct {
	Repository  *Repository
	AccountRepo *accounts.Repository
}

func NewService(repository *Repository, accountRepo *accounts.Repository) *Service {
	return &Service{
		Repository:  repository,
		AccountRepo: accountRepo,
	}
}

func (s *Service) GetAllTransactions() ([]TransactionDTO, error) {
	transactions, err := s.Repository.GetAllTransactions()
	if err != nil {
		log.Printf("Error getting all transactions: %v", err)
		return nil, fmt.Errorf("could not fetch transactions: %v", err)
	}
	return transactions, nil
}

func (s *Service) GetTransactionByID(id int) (*Transaction, error) {
	transaction, err := s.Repository.GetTransactionByID(id)
	if err != nil {
		log.Printf("Error getting transaction by ID: %v", err)
		return nil, fmt.Errorf("could not fetch transaction with ID %d: %v", id, err)
	}
	return transaction, nil
}

func (s *Service) CreateTransaction(t Transaction) (Transaction, error) {
	// Start a new transaction
	tx, err := s.Repository.db.Begin()
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to start transaction: %v", err)
	}
	// Use a local variable for error tracking
	var txErr error
	// Ensure rollback if anything fails
	defer func() {
		if txErr != nil {
			log.Printf("Rolling back transaction due to error: %v", txErr)
			tx.Rollback()
		}
	}()

	// Create Transaction
	createdTransaction, txErr := s.Repository.CreateTransaction(tx, t)
	if txErr != nil {
		return Transaction{}, fmt.Errorf("could not create transaction: %v", txErr)
	}

	// Update account balance
	if createdTransaction.TypeID == 1 {
		// For Expense
		log.Printf("Substracting Balance for ID : %d", createdTransaction.AccountID)
		txErr = s.AccountRepo.SubstractAccountBalance(tx, createdTransaction.AccountID, createdTransaction.Amount)
		if txErr != nil {
			return Transaction{}, txErr
		}
	} else if createdTransaction.TypeID == 2 {
		// For Income
		log.Printf("Add Balance for ID : %d", createdTransaction.AccountID)
		txErr = s.AccountRepo.AddAccountBalance(tx, createdTransaction.AccountID, createdTransaction.Amount)
		if txErr != nil {
			return Transaction{}, txErr
		}
	} else {
		// For Transfer
		log.Printf("Transfer Balance from ID %d to : %d", createdTransaction.AccountID, createdTransaction.AccountDestID)
		txErr = s.AccountRepo.SubstractAccountBalance(tx, createdTransaction.AccountID, createdTransaction.Amount)
		if txErr != nil {
			return Transaction{}, txErr
		}
		txErr = s.AccountRepo.AddAccountBalance(tx, *createdTransaction.AccountDestID, createdTransaction.Amount)
		if txErr != nil {
			return Transaction{}, txErr
		}
	}

	// Commit the transaction
	txErr = tx.Commit()
	if txErr != nil {
		return Transaction{}, fmt.Errorf("failed to commit transaction: %v", txErr)
	}
	return createdTransaction, nil
}

func (s *Service) UpdateTransaction(id int, t Transaction) (Transaction, error) {
	updatedTransaction, err := s.Repository.UpdateTransaction(id, t)
	if err != nil {
		log.Printf("Error updating transaction with ID %d: %v", id, err)
		return Transaction{}, fmt.Errorf("could not update transaction with ID %d: %v", id, err)
	}
	return updatedTransaction, nil
}

func (s *Service) DeleteTransaction(id int) error {
	err := s.Repository.DeleteTransaction(id)
	if err != nil {
		log.Printf("Error deleting transaction with ID %d: %v", id, err)
		return fmt.Errorf("could not delete transaction with ID %d: %v", id, err)
	}
	return nil
}
