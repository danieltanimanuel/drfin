package accounts

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAllAccounts() ([]AccountDTO, error) {
	return s.repository.GetAllAccounts()
}

// GetTransactionAccountGroupByID fetches a transaction type by its ID.
// func (s *Service) GetAccountGroupByID(id int) (*AccountDTO, error) {
// 	return s.repository.GetAccountGroupByID(id)
// }

// // CreateTransactionAccountGroup creates a new transaction type.
// func (s *Service) CreateAccountGroup(c AccountDTO) (AccountDTO, error) {
// 	return s.repository.CreateAccountGroup(c)
// }

// // UpdateTransactionAccountGroup updates an existing transaction type.
// func (s *Service) UpdateAccountGroup(id int, c AccountDTO) (AccountDTO, error) {
// 	return s.repository.UpdateAccountGroup(id, c)
// }

// // DeleteTransactionAccountGroup deletes a transaction type by its ID.
// func (s *Service) DeleteAccountGroup(id int) error {
// 	return s.repository.DeleteAccountGroup(id)
// }
