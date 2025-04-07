package categories

// Service handles the business logic and uses the repository to interact with the database.
type Service struct {
	repository *Repository
}

// NewService creates and returns a new service instance.
func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// GetAllCategories fetches all categories using the repository.
func (s *Service) GetAllCategories() ([]Category, error) {
	return s.repository.GetAllCategories()
}

// GetTransactionTypeByID fetches a transaction type by its ID.
func (s *Service) GetCategoryByID(id int) (*Category, error) {
	return s.repository.GetCategoryByID(id)
}

// CreateTransactionType creates a new transaction type.
func (s *Service) CreateCategory(c Category) (Category, error) {
	return s.repository.CreateCategory(c)
}

// UpdateTransactionType updates an existing transaction type.
func (s *Service) UpdateCategory(id int, c Category) (Category, error) {
	return s.repository.UpdateCategory(id, c)
}

// DeleteTransactionType deletes a transaction type by its ID.
func (s *Service) DeleteCategory(id int) error {
	return s.repository.DeleteCategory(id)
}
