package types

// Service handles the business logic and uses the repository to interact with the database.
type Service struct {
	repository *Repository
}

// NewService creates and returns a new service instance.
func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// GetAllTypes fetches all Types using the repository.
func (s *Service) GetAllTypes() ([]Type, error) {
	return s.repository.GetAllTypes()
}

// GetTransactionTypeByID fetches a transaction type by its ID.
func (s *Service) GetTypeByID(id int) (*Type, error) {
	return s.repository.GetTypeByID(id)
}

// CreateTransactionType creates a new transaction type.
func (s *Service) CreateType(c Type) (Type, error) {
	return s.repository.CreateType(c)
}

// UpdateTransactionType updates an existing transaction type.
func (s *Service) UpdateType(id int, c Type) (Type, error) {
	return s.repository.UpdateType(id, c)
}

// DeleteTransactionType deletes a transaction type by its ID.
func (s *Service) DeleteType(id int) error {
	return s.repository.DeleteType(id)
}
