package account

type service struct {
    repo *repository
}

// Constructor (Go-style)
func NewService(repo *repository) *service {
    return &service{
        repo: repo,
    }
}

// Business logic method
func (s *service) GetTestAccount() *BankAccount {
    return s.repo.GetTestAccount()
}
