package transactiontypes

// Service handles the business logic and uses the repository to interact with the database.
type Service struct {
	repository *Repository
}

// NewService creates and returns a new service instance.
func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// GetAllTransactionTypes fetches all TransactionTypes using the repository.
func (s *Service) GetAllTransactionTypes() ([]TransactionType, error) {
	return s.repository.GetAllTransactionTypes()
}

// GetTransactionTransactionTypeByID fetches a transaction type by its ID.
func (s *Service) GetTransactionTypeByID(id int) (*TransactionType, error) {
	return s.repository.GetTransactionTypeByID(id)
}

// CreateTransactionTransactionType creates a new transaction type.
func (s *Service) CreateTransactionType(c TransactionType) (TransactionType, error) {
	return s.repository.CreateTransactionType(c)
}

// UpdateTransactionTransactionType updates an existing transaction type.
func (s *Service) UpdateTransactionType(id int, c TransactionType) (TransactionType, error) {
	return s.repository.UpdateTransactionType(id, c)
}

// DeleteTransactionTransactionType deletes a transaction type by its ID.
func (s *Service) DeleteTransactionType(id int) error {
	return s.repository.DeleteTransactionType(id)
}
