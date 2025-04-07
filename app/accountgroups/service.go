package accountgroups

// Service handles the business logic and uses the repository to interact with the database.
type Service struct {
	repository *Repository
}

// NewService creates and returns a new service instance.
func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// GetAllAccountGroups fetches all AccountGroups using the repository.
func (s *Service) GetAllAccountGroups() ([]AccountGroup, error) {
	return s.repository.GetAllAccountGroups()
}

// GetTransactionAccountGroupByID fetches a transaction type by its ID.
func (s *Service) GetAccountGroupByID(id int) (*AccountGroup, error) {
	return s.repository.GetAccountGroupByID(id)
}

// CreateTransactionAccountGroup creates a new transaction type.
func (s *Service) CreateAccountGroup(c AccountGroup) (AccountGroup, error) {
	return s.repository.CreateAccountGroup(c)
}

// UpdateTransactionAccountGroup updates an existing transaction type.
func (s *Service) UpdateAccountGroup(id int, c AccountGroup) (AccountGroup, error) {
	return s.repository.UpdateAccountGroup(id, c)
}

// DeleteTransactionAccountGroup deletes a transaction type by its ID.
func (s *Service) DeleteAccountGroup(id int) error {
	return s.repository.DeleteAccountGroup(id)
}
