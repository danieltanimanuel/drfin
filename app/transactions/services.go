package transactions

import (
	"context"
	"database/sql"
	"drfin/app/accounts"
	"drfin/app/categories"
	"drfin/app/shared/logdr"
	"drfin/app/transactiontypes"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strings"
	"time"
)

type Service struct {
	Repository          *Repository
	AccountRepo         *accounts.Repository
	TransactionTypeRepo *transactiontypes.Repository
	CategoryRepo        *categories.Repository
}

const (
	adjustmentTypeName     = "Adjustment"
	adjustmentCategoryName = "Unknown"
	maxAdjustmentNoteLen   = 500
)

func NewService(repository *Repository, accountRepo *accounts.Repository, transactionTypeRepo *transactiontypes.Repository, categoryRepo *categories.Repository) *Service {
	return &Service{
		Repository:          repository,
		AccountRepo:         accountRepo,
		TransactionTypeRepo: transactionTypeRepo,
		CategoryRepo:        categoryRepo,
	}
}

func (s *Service) GetAllTransactions(ctx context.Context) ([]DailyTransactionsDTO, error) {

	transactions, err := s.Repository.GetAllTransactions(ctx)
	if err != nil {
		log.Printf("Error getting all transactions: %v", err)
		return nil, fmt.Errorf("could not fetch transactions: %w", err)
	}

	dailyMap := make(map[string]*DailyTransactionsDTO)

	for _, tx := range transactions {
		dayKey := tx.TradeDate.Format("2006-01-02")
		if _, ok := dailyMap[dayKey]; !ok {
			dailyMap[dayKey] = &DailyTransactionsDTO{
				Date: tx.TradeDate,
			}
		}
		if dailyMap[dayKey].IncomeTotal == nil {
			dailyMap[dayKey].IncomeTotal = big.NewFloat(0)
		}
		if dailyMap[dayKey].ExpenseTotal == nil {
			dailyMap[dayKey].ExpenseTotal = big.NewFloat(0)
		}
		switch tx.Type.Name {
		case "Income":
			dailyMap[dayKey].IncomeTotal =
				new(big.Float).Add(dailyMap[dayKey].IncomeTotal, tx.Amount)
		case "Expense":
			dailyMap[dayKey].ExpenseTotal =
				new(big.Float).Add(dailyMap[dayKey].ExpenseTotal, tx.Amount)
		}
		dailyMap[dayKey].Transactions =
			append(dailyMap[dayKey].Transactions, tx)
	}

	// convert map → slice
	var grouped []DailyTransactionsDTO
	for _, v := range dailyMap {
		grouped = append(grouped, *v)
	}

	// sort by date desc
	sort.Slice(grouped, func(i, j int) bool {
		return grouped[i].Date.After(grouped[j].Date)
	})

	return grouped, nil
}

func (s *Service) GetTransactionByID(ctx context.Context, id int) (*Transaction, error) {

	transaction, err := s.Repository.GetTransactionByID(ctx, id)
	if err != nil {
		log.Printf("Error getting transaction by ID: %v", err)
		return nil, fmt.Errorf("could not fetch transaction with ID %d: %v", id, err)
	}

	return transaction, nil
}

func (s *Service) GetTransactionByIDDTO(ctx context.Context, id int) (*TransactionDTO, error) {

	transaction, err := s.Repository.GetTransactionByIDDTO(ctx, id)
	if err != nil {
		log.Printf("Error getting transaction by ID: %v", err)
		return nil, fmt.Errorf("could not fetch transaction with ID %d: %v", id, err)
	}

	return transaction, nil
}

func (s *Service) GetExpenseCategoryAggregates(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]CategoryAggregate, error) {

	return s.Repository.GetCategoryAggregatesByType(
		ctx,
		startDate,
		endDate,
		"Expense",
	)
}

func (s *Service) CreateTransaction(ctx context.Context, t CreateTransactionDTO) (Transaction, error) {

	// Begin SQL transaction WITH context
	tx, err := s.Repository.db.BeginTx(ctx, nil)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to start transaction: %w", err)
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

	// Inquiry necessary ID (type, account, category)
	transactiontype, txErr := s.TransactionTypeRepo.GetTransactionTypeByName(t.Type)
	if txErr != nil {
		return Transaction{}, fmt.Errorf("could not inquiry transaction types due to error : %v", txErr)
	}
	category, txErr := s.CategoryRepo.GetCategoryByName(t.Category)
	if txErr != nil {
		return Transaction{}, fmt.Errorf("could not inquiry category due to error : %v", txErr)
	}
	account, txErr := s.AccountRepo.GetAccountByName(t.Account)
	if txErr != nil {
		return Transaction{}, fmt.Errorf("could not inquiry account due to error : %v", txErr)
	}

	// If Expense, accountDestID nil, else get the accountDestID from database
	var accountDestID *int
	accountDestID = nil
	if transactiontype.ID != 1 { // not Expense
		accountDest, txErr := s.AccountRepo.GetAccountByName(t.AccountDest)
		if txErr != nil {
			return Transaction{}, fmt.Errorf("could not inquiry account due to error : %v", txErr)
		}
		accountDestID = &accountDest.ID
	}
	// Status to "Waiting" if coming from email, else "Confirmed"
	status := "Waiting"
	if t.Source != "Email" {
		status = "Confirmed"
	}

	input := Transaction{
		TradeDate:     t.TradeDate,
		TypeID:        transactiontype.ID,
		CategoryID:    category.ID,
		AccountID:     account.ID,
		AccountDestID: accountDestID,
		Amount:        t.Amount,
		Description:   t.Description,
		Status:        status,
		Source:        t.Source,
	}

	// ---- Persist transaction ----

	created, txErr := s.Repository.CreateTransaction(tx, input)
	if txErr != nil {
		return Transaction{}, fmt.Errorf("could not create transaction: %w", txErr)
	}

	// ---- Apply balance (single source of truth) ----

	txErr = s.applyBalance(tx, created)
	if txErr != nil {
		return Transaction{}, txErr
	}

	// ---- Commit ----

	txErr = tx.Commit()
	if txErr != nil {
		return Transaction{}, fmt.Errorf("failed to commit transaction: %w", txErr)
	}

	return created, nil
}

func affectsBalance(dto UpdateTransactionDBDTO) bool {
	return dto.Amount != nil ||
		dto.AccountID != nil ||
		dto.AccountDestID != nil ||
		dto.TypeID != nil
}

func (s *Service) reverseBalance(
	tx *sql.Tx,
	t Transaction,
) error {

	switch t.TypeID {
	case 1: // Expense
		return s.AccountRepo.AddAccountBalance(
			tx,
			t.AccountID,
			t.Amount,
		)

	case 2: // Income
		return s.AccountRepo.SubstractAccountBalance(
			tx,
			t.AccountID,
			t.Amount,
		)

	case 3: // Transfer
		if err := s.AccountRepo.AddAccountBalance(
			tx,
			t.AccountID,
			t.Amount,
		); err != nil {
			return err
		}
		return s.AccountRepo.SubstractAccountBalance(
			tx,
			*t.AccountDestID,
			t.Amount,
		)
		
	default:
		return fmt.Errorf("unknown transaction type: %d", t.TypeID)
	}
	
}

func (s *Service) applyBalance(
	tx *sql.Tx,
	t Transaction,
) error {

	switch t.TypeID {
	case 1: // Expense
		return s.AccountRepo.SubstractAccountBalance(
			tx,
			t.AccountID,
			t.Amount,
		)

	case 2: // Income
		return s.AccountRepo.AddAccountBalance(
			tx,
			t.AccountID,
			t.Amount,
		)

	default: // Transfer
		if err := s.AccountRepo.SubstractAccountBalance(
			tx,
			t.AccountID,
			t.Amount,
		); err != nil {
			return err
		}
		return s.AccountRepo.AddAccountBalance(
			tx,
			*t.AccountDestID,
			t.Amount,
		)
	}
}

func (s *Service) UpdateTransaction(ctx context.Context, id int, dto UpdateTransactionDTO) (Transaction, error) {
	tx, err := s.Repository.db.BeginTx(ctx, nil)
	if err != nil {
		return Transaction{}, err
	}

	var txErr error
	defer func() {
		if txErr != nil {
			tx.Rollback()
		}
	}()

	// 1️⃣ Load existing transaction
	oldTx, txErr := s.Repository.GetTransactionByID(ctx, id)
	if txErr != nil {
		return Transaction{}, txErr
	}

	// 2️⃣ Resolve DTO → DB DTO (labels → IDs)
	dbDTO := UpdateTransactionDBDTO{}

	if dto.Amount != nil {
		dbDTO.Amount = dto.Amount
	}

	if dto.Type != nil {
		tt, err := s.TransactionTypeRepo.GetTransactionTypeByName(*dto.Type)
		if err != nil {
			return Transaction{}, err
		}
		dbDTO.TypeID = &tt.ID

		// Any non-transfer transaction cannot have destination.
		if tt.ID != 3 {
			dbDTO.ClearAccountDest = true
		}
	}

	if dto.Account != nil {
		acc, err := s.AccountRepo.GetAccountByName(*dto.Account)
		if err != nil {
			return Transaction{}, err
		}
		dbDTO.AccountID = &acc.ID
	}

	if dto.NewReceivable != nil {

		if dto.Type == nil {
			return Transaction{}, fmt.Errorf(
				"transaction type is required when creating a receivable",
			)
		}

		if *dto.Type != "Transfer" {
			return Transaction{}, fmt.Errorf(
				"receivable destination can only be used for Transfer",
			)
		}

		accountName := strings.TrimSpace(*dto.NewReceivable)

		if accountName == "" {
			return Transaction{}, fmt.Errorf(
				"receivable account name cannot be empty",
			)
		}

		account, err := s.AccountRepo.CreateReceivableAccount(
			tx,
			accountName,
		)
		if err != nil {
			return Transaction{}, err
		}

		dbDTO.AccountDestID = &account.ID
	}

	if dto.AccountDest != nil && *dto.AccountDest != "" {
		acc, err := s.AccountRepo.GetAccountByName(*dto.AccountDest)
		if err != nil {
			return Transaction{}, err
		}

		dbDTO.AccountDestID = &acc.ID
	}

	if dto.Category != nil {
		cat, err := s.CategoryRepo.GetCategoryByName(*dto.Category)
		if err != nil {
			txErr = err
			return Transaction{}, txErr
		}
		dbDTO.CategoryID = &cat.ID
	}

	dbDTO.Description = dto.Description
	dbDTO.Status = dto.Status
	dbDTO.TradeDate = dto.TradeDate

	logdr.Info("DBDTO Category", "ID", dbDTO.CategoryID)

	// 3️⃣ Update transaction row
	updatedTx, txErr :=
		s.Repository.UpdateTransactionTx(ctx, tx, id, dbDTO)
	if txErr != nil {
		return Transaction{}, txErr
	}

	// 4️⃣ Determine if balance is affected
	balanceAffected :=
		oldTx.Amount.Cmp(updatedTx.Amount) != 0 ||
			oldTx.TypeID != updatedTx.TypeID ||
			oldTx.AccountID != updatedTx.AccountID ||
			(oldTx.AccountDestID != nil && updatedTx.AccountDestID != nil &&
				*oldTx.AccountDestID != *updatedTx.AccountDestID) ||
			(oldTx.AccountDestID == nil && updatedTx.AccountDestID != nil) ||
			(oldTx.AccountDestID != nil && updatedTx.AccountDestID == nil)

	if balanceAffected {

		// Disallow same source & destination for transfer
		if updatedTx.TypeID == 3 &&
			updatedTx.AccountDestID != nil &&
			updatedTx.AccountID == *updatedTx.AccountDestID {

			txErr = fmt.Errorf("source and destination account must differ")
			return Transaction{}, txErr
		}

		if updatedTx.TypeID == 3 && updatedTx.AccountDestID == nil {
			txErr = fmt.Errorf("transfer requires a destination account")
			return Transaction{}, txErr
		}	

		if updatedTx.TypeID != 3 && updatedTx.AccountDestID != nil {
			txErr = fmt.Errorf("only transfer transactions can have a destination account")
			return Transaction{}, txErr
		}

		// Reverse old balance
		txErr = s.reverseBalance(tx, *oldTx)
		if txErr != nil {
			return Transaction{}, txErr
		}

		// Apply new balance
		txErr = s.applyBalance(tx, updatedTx)
		if txErr != nil {
			return Transaction{}, txErr
		}
	}

	// 5️⃣ Commit
	txErr = tx.Commit()
	if txErr != nil {
		return Transaction{}, txErr
	}

	return updatedTx, nil
}

func (s *Service) DeleteTransaction(ctx context.Context, id int) error {
	err := s.Repository.DeleteTransaction(ctx, id)
	if err != nil {
		log.Printf("Error deleting transaction with ID %d: %v", id, err)
		return fmt.Errorf("could not delete transaction with ID %d: %v", id, err)
	}
	return nil
}

func (s *Service) AdjustAccountBalance(ctx context.Context, accountID int, req AdjustBalanceRequest) (AdjustBalanceResult, error) {
	if accountID <= 0 {
		return AdjustBalanceResult{}, fmt.Errorf("%w: invalid account id", accounts.ErrInvalidInput)
	}
	newBalance, err := accounts.NormalizeAmount(req.NewBalance)
	if err != nil {
		return AdjustBalanceResult{}, err
	}
	note := strings.TrimSpace(req.Note)
	if len([]rune(note)) > maxAdjustmentNoteLen {
		return AdjustBalanceResult{}, fmt.Errorf("%w: note is too long", accounts.ErrInvalidInput)
	}

	// Lookups before the SQL transaction so it stays short.
	adjType, err := s.TransactionTypeRepo.GetTransactionTypeByName(adjustmentTypeName)
	if err != nil {
		return AdjustBalanceResult{}, err
	}
	category, err := s.CategoryRepo.GetCategoryByName(adjustmentCategoryName)
	if err != nil {
		return AdjustBalanceResult{}, err
	}

	tx, err := s.Repository.db.BeginTx(ctx, nil)
	if err != nil {
		return AdjustBalanceResult{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() // no-op after a successful Commit

	oldBalance, delta, err := s.AccountRepo.LockAndDiffBalance(tx, accountID, newBalance)
	if err != nil {
		return AdjustBalanceResult{}, err
	}

	result := AdjustBalanceResult{
		AccountID:  accountID,
		OldBalance: oldBalance,
		NewBalance: newBalance,
		Delta:      delta,
	}

	d, ok := new(big.Rat).SetString(delta)
	if !ok {
		return AdjustBalanceResult{}, fmt.Errorf("could not parse balance difference %q", delta)
	}
	if d.Sign() == 0 {
		return result, nil // nothing to adjust
	}

	description := fmt.Sprintf("Balance adjustment: %s -> %s", oldBalance, newBalance)
	if note != "" {
		description += " (" + note + ")"
	}

	id, err := s.Repository.InsertAdjustmentTx(tx, accountID, adjType.ID, category.ID, delta, description, time.Now())
	if err != nil {
		return AdjustBalanceResult{}, err
	}
	if err := s.AccountRepo.SetAccountBalance(tx, accountID, newBalance); err != nil {
		return AdjustBalanceResult{}, err
	}
	if err := tx.Commit(); err != nil {
		return AdjustBalanceResult{}, fmt.Errorf("failed to commit adjustment: %w", err)
	}

	result.TransactionID = id
	return result, nil
}