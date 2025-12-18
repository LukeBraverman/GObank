package account

import "fmt"

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{
        repo: repo,
    }
}

func (s *Service) OpenAccount(
    accountNumber string,
    name string,
    initialBalance float64,
) (*TokenAccount, error) {

    if accountNumber == "" {
        return nil, fmt.Errorf("account number required")
    }

    if initialBalance < 0 {
        return nil, fmt.Errorf("initial balance cannot be negative")
    }

    acc := &TokenAccount{
        AccountNumber: accountNumber,
        Name:          name,
        Balance:       initialBalance,
    }

    if err := s.repo.CreateAccount(acc); err != nil {
        return nil, err
    }

    return acc, nil
}

func (s *Service) GetAccount(accountNumber string) (*TokenAccount, error) {
    acc, ok := s.repo.GetAccount(accountNumber)
    if !ok {
        return nil, fmt.Errorf("account not found")
    }
    return acc, nil
}
