package account

import (
	"fmt"
	"sync"
)

type InMemoryRepository struct {
	mu       sync.RWMutex
	accounts map[string]*TokenAccount
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{
		accounts: make(map[string]*TokenAccount),
	}
}

func (r *InMemoryRepository) CreateAccount(account *TokenAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[account.AccountNumber]; exists {
		return fmt.Errorf("account %s already exists", account.AccountNumber)
	}

	r.accounts[account.AccountNumber] = account
	return nil
}

func (r *InMemoryRepository) GetAccount(accountNumber string) (*TokenAccount, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[accountNumber]
	return account, exists, nil
}

// func (r *Repository) UpdateBalance(accountNumber string, delta float64) error {
//     r.mu.Lock()
//     defer r.mu.Unlock()

//     account, exists := r.accounts[accountNumber]
//     if !exists {
//         return fmt.Errorf("account not found")
//     }

//     account.Balance += delta
//     return nil
// }
